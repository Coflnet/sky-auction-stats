package redis

import (
	"context"
	"fmt"
	"github.com/Coflnet/auction-stats/internal/model"
	"github.com/Coflnet/auction-stats/internal/prometheus"
	"github.com/go-redis/redis/v9"
	"time"
)

// CountFlipSummary adds a flip summary to the user counting
func CountFlipSummary(flip *model.Flip) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		key := FlipSummaryPlayerHourKey(flip)
		pipe.PFAdd(ctx, key, flip.Sell.UUID)
		pipe.Expire(ctx, key, time.Hour*24*5)

		return nil
	})

	if err != nil {
		return err
	}

	prometheus.FlipSummaryAdded()

	return nil
}

func UpdateFlipBuyerCount(flip *model.Flip) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	cmds, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.PFCount(ctx, FlipSummaryPlayerHourKey(flip))
		return nil
	})
	if err != nil {
		panic(err)
	}

	sum := 0
	for _, cmd := range cmds {
		v := cmd.(*redis.IntCmd).Val()
		sum += int(v)
	}

	flip.AmountOfFlipsFromBuyerOfTheDay = sum
	return nil
}

func FlipSummaryPlayerHourKey(flip *model.Flip) string {
	return fmt.Sprintf("%s_%s", flip.Buy.ProfileID, HourKey(flip.Sell.End, true))
}
