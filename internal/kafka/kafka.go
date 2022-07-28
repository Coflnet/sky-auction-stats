package kafka

import (
	"fmt"
	"github.com/rs/zerolog/log"
	kafkago "github.com/segmentio/kafka-go"
	"os"
)

var (
	conn              *kafkago.Conn
	auctionTopic      string
	flipTopic         string
	auctionReader     *kafkago.Reader
	flipSummaryReader *kafkago.Reader
	flipSummaryWriter *kafkago.Writer
)

func Init() error {

	auctionTopic = os.Getenv("TOPIC_NEW_AUCTION")
	flipTopic = os.Getenv("TOPIC_NEW_FLIP")

	if flipTopic == "" || auctionTopic == "" {
		return fmt.Errorf("TOPIC env vars not set")
	}

	auctionReader = kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_HOST")},
		GroupID:  os.Getenv("KAFKA_CONSUMER_GROUP"),
		Topic:    auctionTopic,
		MaxBytes: 10e6,
		MinBytes: 10e3,
	})

	flipSummaryReader = kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_HOST")},
		GroupID:  "auction-stats-reader",
		Topic:    flipSummaryTopic(),
		MaxBytes: 10e6,
		MinBytes: 10e3,
	})

	flipSummaryWriter = &kafkago.Writer{
		Addr:                   kafkago.TCP(os.Getenv("KAFKA_HOST")),
		Topic:                  flipSummaryProcessedTopic(),
		AllowAutoTopicCreation: true,
	}

	go StartReaders()

	return nil
}

func Disconnect() error {
	return conn.Close()
}

func StartReaders() {

	defer func() {
		if err := auctionReader.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close auction reader")
			return
		}

		if err := flipSummaryReader.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close flipsummary reader")
			return
		}

		log.Info().Msgf("closed kafka reader gracefully")
	}()

	go func() {
		for {
			err := ReadAuctions()
			if err != nil {
				log.Error().Err(err).Msgf("error consuming new auctions")
			}
		}
	}()

	for {
		err := ReadFlipSummaries()
		if err != nil {
			log.Error().Err(err).Msgf("error consuming new flipsummaries")
		}
	}
}

func flipSummaryTopic() string {
	flipSummaryTopic := os.Getenv("TOPIC_FLIP_SUMMARY")

	if flipSummaryTopic == "" {
		log.Panic().Msg("TOPIC_FLIP_SUMMARY env var not set")
	}

	return flipSummaryTopic
}

func flipSummaryProcessedTopic() string {
	return "flip-summary-processed"
}
