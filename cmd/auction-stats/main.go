package main

import (
	"github.com/Coflnet/auction-stats/internal/api"
	"github.com/Coflnet/auction-stats/internal/kafka"
	"github.com/Coflnet/auction-stats/internal/mongo"
	"github.com/Coflnet/auction-stats/internal/prometheus"
	"github.com/Coflnet/auction-stats/internal/redis"
	"github.com/Coflnet/auction-stats/internal/usecase"
	"github.com/rs/zerolog/log"
)

func main() {

	if err := redis.Init(); err != nil {
		log.Panic().Err(err).Msg("can not init redis")
	}

	if err := mongo.Init(); err != nil {
		log.Panic().Err(err).Msg("can not init mongo")
	}
	defer func() {
		_ = mongo.Disconnect()
	}()

	err := kafka.Init()
	if err != nil {
		log.Panic().Err(err).Msgf("kafka init was not successfully")
		return
	}
	defer func() {
		err := kafka.Disconnect()
		if err != nil {
			log.Panic().Err(err).Msgf("can not close connection to kafka")
		}
	}()

	prometheus.StartPrometheus()
	usecase.StartNotifierSchedule()

	err = api.StartApi()
	log.Panic().Err(err).Msgf("api stopped")
}
