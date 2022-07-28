package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	flipSummaryAdded = promauto.NewCounter(prometheus.CounterOpts{
		Name: "auction_stats_flip_summary_added",
		Help: "The total number of added flip summaries",
	})

	flipSummaryRequested = promauto.NewCounter(prometheus.CounterOpts{
		Name: "auction_stats_flip_summary_requested",
		Help: "The total number of processed flip summary requests",
	})
)

func FlipSummaryAdded() {
	flipSummaryAdded.Inc()
}

func FlipSummaryRequested() {
	flipSummaryRequested.Inc()
}
