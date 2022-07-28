package kafka

import (
	"context"
	"encoding/json"
	"github.com/Coflnet/auction-stats/internal/model"
	"github.com/Coflnet/auction-stats/internal/redis"
	"github.com/segmentio/kafka-go"
	"time"
)

func ReadFlipSummaries() error {

	// if offset is big use a bigger batch size
	offset := flipSummaryReader.Offset()
	batchSize := 5
	if offset >= 100000 {
		batchSize = 1000
	}

	ctx := context.Background()

	var flips []*model.Flip
	var msgs []kafka.Message

	for i := 0; i < batchSize; i++ {
		msg, err := flipSummaryReader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		var flip model.Flip
		err = json.Unmarshal(msg.Value, &flip)
		if err != nil {
			return err
		}

		flips = append(flips, &flip)
		msgs = append(msgs, msg)
	}

	err := processFlips(flips)
	if err != nil {
		return err
	}

	// write flips into kafka topic for discord bot to read
	err = writeFlipsToDiscordChat(flips)

	err = flipSummaryReader.CommitMessages(ctx, msgs...)

	return err
}

func processFlips(flips []*model.Flip) error {

	errCh := make(chan error)

	for _, flip := range flips {
		go func(f *model.Flip, ch chan error) {
			ch <- processFlip(f)
		}(flip, errCh)
	}

	for i := 0; i < len(flips); i++ {
		err := <-errCh
		if err != nil {
			return err
		}
	}

	return nil
}

func processFlip(flip *model.Flip) error {
	err := redis.CountFlipSummary(flip)
	if err != nil {
		return err
	}

	// update the flip buyer count
	return redis.UpdateFlipBuyerCount(flip)
}

func writeFlipsToDiscordChat(flips []*model.Flip) error {

	messages := make([]kafka.Message, len(flips))
	for i, flip := range flips {
		msg, err := json.Marshal(flip)
		if err != nil {
			return err
		}
		messages[i] = kafka.Message{
			Key:   []byte(flip.Sell.UUID),
			Value: msg,
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return flipSummaryWriter.WriteMessages(ctx, messages...)
}
