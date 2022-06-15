package main

import (
	"github.com/Coflnet/auction-stats/internal/api"
	"github.com/Coflnet/auction-stats/internal/kafka"
	"github.com/Coflnet/auction-stats/internal/prometheus"
	"github.com/Coflnet/auction-stats/internal/redis"
	"github.com/rs/zerolog/log"
)

func main() {

	if err := redis.Init(); err != nil {
		log.Panic().Err(err).Msg("can not init redis")
	}

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

	err = prometheus.InitPrometheus()
	if err != nil {
		log.Panic().Err(err).Msg("failed to init prometheus")
		return
	}

	err = api.StartApi()
	log.Panic().Err(err).Msgf("api stopped")
}
