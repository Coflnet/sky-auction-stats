package usecase

import (
	"fmt"
	"github.com/Coflnet/auction-stats/internal/model"
	"github.com/Coflnet/auction-stats/internal/redis"
	"github.com/rs/zerolog/log"
	"time"
)

func AuctionCount(duration time.Duration) (*model.AuctionResponse, error) {

	if !DurationValid(duration) {
		return nil, &InvalidDurationError{duration: duration}
	}

	start := time.Now().Add(-duration)
	end := time.Now()

	val, err := redis.ReceiveAuctionCount(time.Now(), duration)
	if err != nil {
		log.Error().Err(err).Msgf("error counting redis auctions")
		return nil, fmt.Errorf("error when counting auctions within redis")
	}

	return &model.AuctionResponse{
		Count: val,
		From:  start,
		To:    end,
	}, nil
}

type InvalidDurationError struct {
	duration time.Duration
}

func (e *InvalidDurationError) Error() string {
	return fmt.Sprintf("the duration %v is not supported, duration has to be more than 1 minute and less than 1 year", e.duration.Minutes())
}

func DurationValid(duration time.Duration) bool {

	if duration.Hours() < 24*30*12 {
		log.Warn().Msgf("more than 1 year is not supported")
		return false
	}

	if duration.Minutes() < 1 {
		log.Warn().Msgf("less than 1 minute is not supported")
		return false
	}

	return true
}
