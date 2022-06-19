package prometheus

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
)

func StartPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":2112", nil)
		if err != nil {
			log.Panic().Err(err).Msgf("Error listening on port 2112")
			return
		}
	}()
}
