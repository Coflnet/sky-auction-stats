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
)

func ObserveAuctionAt() {
	auctionsAdded.Inc()
}
