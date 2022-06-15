package redis

import (
	"context"
	"fmt"
	"github.com/Coflnet/auction-stats/internal/model"
	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

// CountAuction adds an auction to the counting system
func CountAuction(auction *model.Auction) error {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {

		key := MinuteKey(auction.Start)
		pipe.PFAdd(ctx, key, auction.UUID)

		key = HourKey(auction.Start)
		pipe.PFAdd(ctx, key, auction.UUID)

		key = DayKey(auction.Start)
		pipe.PFAdd(ctx, key, auction.UUID)

		key = MonthKey(auction.Start)
		pipe.PFAdd(ctx, key, auction.UUID)

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// ReceiveAuctionCount returns the amount of auctions for a given time
// time parameter is when the count happened
// duration is which counting duration should be chosen
// example time: 10.06.2022 14:49:52, duration: 1h -> Counts from 10.06.2022 13:49 - 14:49
// example time: 10.06.2022 14:49:52, duration: 1d -> Counts from 09.06.2022 14:49 - 10.06.2022 14:49
// example time: 10.06.2022 14:49:52, duration: 10m -> Counts from 10.06.2022 14:39 - 14:49
func ReceiveAuctionCount(time time.Time, duration time.Duration) (int64, error) {
	keys, err := keysForTimeDuration(time, duration)
	if err != nil {
		log.Error().Err(err).Msgf("error searching keys for counting, time: %v, duration: %v", time, duration)
		return 0, err
	}

	count, err := countKeys(keys)
	if err != nil {
		log.Error().Err(err).Msgf("error counting keys, time: %v, duration: %v", time, duration)
		return 0, err
	}

	return count, nil
}

func keysForTimeDuration(t time.Time, d time.Duration) ([]string, error) {
	start := t.Add(-d)

	log.Info().Msgf("searching keys for time: %v and duration %v", t, d)

	if d < time.Hour {
		return everyMinBetween(start, t, []string{}), nil
	}

	if d < time.Hour*24 {
		return everyHourBetween(start, t, []string{}), nil
	}

	if d < time.Hour*24*30 {
		return everyDayBetween(start, t, []string{}), nil
	}

	return nil, fmt.Errorf("duration is less than a min")
}

func everyMinBetween(start, end time.Time, keys []string) []string {
	if end.Sub(start) < 1*time.Minute {
		return keys
	}

	// add start key to array
	key := MinuteKey(start)

	// count start 1 m up
	newStart := start.Add(time.Minute * 1)

	keys = append(keys, key)
	return everyMinBetween(newStart, end, keys)
}

func everyHourBetween(start, end time.Time, keys []string) []string {
	if end.Sub(start) < 1*time.Hour {
		return keys
	}

	key := HourKey(start)
	newStart := start.Add(time.Hour * 1)

	keys = append(keys, key)
	return everyHourBetween(newStart, end, keys)
}

func everyDayBetween(start, end time.Time, keys []string) []string {
	if end.Sub(start) < 24*time.Hour {
		return keys
	}

	key := DayKey(start)
	newStart := start.Add(time.Hour * 24)

	keys = append(keys, key)
	return everyDayBetween(newStart, end, keys)
}

func MinuteKey(t time.Time) string {
	m := strconv.Itoa(t.Minute())
	return HourKey(t) + m
}

func HourKey(t time.Time) string {
	h := strconv.Itoa(t.Hour())
	return DayKey(t) + h
}

func DayKey(t time.Time) string {
	d := strconv.Itoa(t.Day())
	return MonthKey(t) + d
}

func MonthKey(t time.Time) string {
	y := strconv.Itoa(t.Year())
	m := int(t.Month())
	return y + strconv.Itoa(m)
}

func countKeys(keys []string) (int64, error) {

	if (len(keys)) == 0 {
		return 0, fmt.Errorf("can not count 0 keys")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	log.Info().Msgf("counting %d keys", len(keys))
	for _, k := range keys {
		log.Info().Msgf("key: %s", k)
	}

	return rdb.PFCount(ctx, keys...).Result()
}
