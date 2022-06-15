package kafka

import (
	"context"
	"github.com/Coflnet/auction-stats/internal/model"
	"github.com/Coflnet/auction-stats/internal/prometheus"
	"github.com/Coflnet/auction-stats/internal/redis"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"time"
)

func ReadAuctions() error {
	m, err := auctionReader.FetchMessage(context.Background())
	if err != nil {
		return err
	}

	err = processMessage(&m)
	if err != nil {
		log.Error().Err(err).Msgf("error processing message")
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	err = auctionReader.CommitMessages(ctx, m)

	return nil
}

// processMessage takes a new-auction message from kafka and processes everything that has to be done
// example key of a message: 7807253170724460696506/09/2022 12:23:12
func processMessage(msg *kafka.Message) error {

	uuid := string(msg.Key)[:20]
	t := string(msg.Key)[len(string(msg.Key))-len("06/14/2022 04:55:56"):]

	timestamp, err := time.Parse("01/02/2006 15:04:05", t)
	if err != nil {
		log.Error().Err(err).Msgf("can not convert timestamp: %s", t)
		return err
	}

	auction := model.Auction{
		Start: timestamp,
		UUID:  uuid,
	}

	err = redis.CountAuction(&auction)
	if err != nil {
		log.Error().Err(err).Msgf("error counting auction")
		return err
	}

	prometheus.ObserveAuctionAt()

	return nil
}
