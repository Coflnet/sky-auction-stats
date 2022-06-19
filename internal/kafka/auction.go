package kafka

import (
	"context"
	"github.com/Coflnet/auction-stats/internal/model"
	"github.com/Coflnet/auction-stats/internal/redis"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"sync"
	"time"
)

func ReadAuctions() error {

	messages := make([]kafka.Message, 0)
	for i := 0; i < 1000; i++ {
		m, err := auctionReader.FetchMessage(context.Background())
		if err != nil {
			return err
		}
		messages = append(messages, m)
	}

	processedMessages := processMessages(messages)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	count := 0
	for processedMessage := range processedMessages {
		err := auctionReader.CommitMessages(ctx, processedMessage)
		if err != nil {
			log.Error().Err(err).Msgf("error committing message")
			return err
		}

		count++
	}

	log.Info().Msgf("processed batch of %d messages", count)

	return nil
}

// processMessage takes a new-auction message from kafka and processes everything that has to be done
// example key of a message: 7807253170724460696506/09/2022 12:23:12
func processMessages(messages []kafka.Message) <-chan kafka.Message {

	ch := make(chan kafka.Message)

	go func() {
		wg := sync.WaitGroup{}
		for _, m := range messages {
			wg.Add(1)

			go func(msg kafka.Message) {
				defer wg.Done()

				uuid := string(msg.Key)[:20]
				t := string(msg.Key)[len(string(msg.Key))-len("06/14/2022 04:55:56"):]

				timestamp, err := time.Parse("01/02/2006 15:04:05", t)
				if err != nil {
					log.Error().Err(err).Msgf("can not convert timestamp: %s", t)
					return
				}

				auction := model.Auction{
					Start: timestamp,
					UUID:  uuid,
				}

				err = redis.CountAuction(&auction)
				if err != nil {
					log.Error().Err(err).Msgf("error counting auction")

					return
				}

				ch <- msg
			}(m)
		}
		wg.Wait()
		close(ch)
	}()

	return ch
}
