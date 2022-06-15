package kafka

import (
	"fmt"
	"github.com/rs/zerolog/log"
	kafkago "github.com/segmentio/kafka-go"
	"os"
)

var (
	conn          *kafkago.Conn
	auctionTopic  string
	flipTopic     string
	auctionReader *kafkago.Reader
)

func Init() error {

	auctionTopic = os.Getenv("TOPIC_NEW_AUCTION")
	flipTopic = os.Getenv("TOPIC_NEW_FLIP")

	if flipTopic == "" || auctionTopic == "" {
		return fmt.Errorf("TOPIC env vars not set")
	}

	auctionReader = kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:   []string{os.Getenv("KAFKA_HOST")},
		GroupID:   os.Getenv("KAFKA_CONSUMER_GROUP"),
		Topic:     auctionTopic,
		Partition: 0,
		MaxBytes:  10e6,
		MinBytes:  10e3,
	})

	go func() {
		err := StartReaders()
		if err != nil {
			log.Panic().Err(err).Msgf("error consuming")
		}
	}()

	return nil
}

func Disconnect() error {
	return conn.Close()
}

func StartReaders() error {

	defer func() {
		if err := auctionReader.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close reader:")
			return
		}

		log.Info().Msgf("closed kafka reader gracefully")
	}()

	runs := 0
	for {
		err := ReadAuctions()
		if err != nil {
			log.Panic().Err(err).Msgf("error consuming new auctions")
		}

		if runs >= 20000 {
			log.Info().Msgf("inserted %d messages", runs)
			runs = 0
		}
		runs++
	}
}
