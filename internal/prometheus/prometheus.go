package prometheus

import (
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/rs/zerolog/log"
	"os"
)

var (
	client api.Client
)

func InitPrometheus() error {
	var err error
	addr := os.Getenv("PROMETHEUS_ADDRESS")
	if addr == "" {
		addr = "http://prometheus-prometheus.prometheus:9090"
	}
	client, err = api.NewClient(api.Config{
		Address: addr,
	})

	if err != nil {
		log.Panic().Err(err).Msgf("Error creating prometheus client")
		return err
	}

	// TODO use prometheus
	_ = v1.NewAPI(client)
	return nil
}
