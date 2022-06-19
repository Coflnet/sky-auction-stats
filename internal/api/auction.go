package api

import (
	"github.com/Coflnet/auction-stats/internal/prometheus"
	"github.com/Coflnet/auction-stats/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"time"
)

// NewAuctions api endpoint to retrieve the amount of the new auctions
// @Summary Get recent auctions
// @Description returns the amount of auctions within in the last x minutes
// @Tags Stats
// @Accept json
// @Produce json
// @Param duration query int false "duration in minutes" minimum(1) maximum(2880)
// @Router /new-auctions [get]
func NewAuctions(c *gin.Context) {
	d := c.DefaultQuery("duration", "5")
	if d == "" {
		d = "5"
	}

	log.Debug().Msgf("api input: %v", d)

	m, err := strconv.Atoi(d)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid duration, must be an integer",
		})
		return
	}

	duration := time.Duration(m) * time.Minute

	data, err := usecase.AuctionCount(duration)
	if err != nil {
		if e, ok := err.(*usecase.InvalidDurationError); ok {
			log.Warn().Msgf("searched new auctions, but the duration given is invalid: %s", d)
			c.JSON(400, gin.H{
				"error": e.Error(),
			})
			return
		}

		log.Error().Err(err).Msgf("searched new auctions, but an internal error occurred, given duration: %s", d)
		c.JSON(500, gin.H{
			"error": "an internal error occurred",
		})
		return
	}

	prometheus.RequestsProcessed()
	log.Info().Msgf("searched new auctions found: %d", data)
	c.JSON(http.StatusOK, data)
}
