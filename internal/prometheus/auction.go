package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	auctionsAdded = promauto.NewCounter(prometheus.CounterOpts{
		Name: "auction_stats_auctions_added",
		Help: "The total number of processed auctions",
	})

	requestsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "auction_stats_requests_processed",
		Help: "The total number of processed requests",
	})

	notifierEvaluated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "auction_stats_notifiers_evaluated",
		Help: "The total number of evaluated notifiers",
	})
)

func AuctionAdded() {
	auctionsAdded.Inc()
}

func RequestsProcessed() {
	requestsProcessed.Inc()
}

func NotifiersEvaluated() {
	notifierEvaluated.Inc()
}
