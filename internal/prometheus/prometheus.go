package prometheus

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartPrometheus() error {
	return runPrometheus()
}

func setupPrometheus() *gin.Engine {
	r := gin.Default()
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	return r
}

func runPrometheus() error {
	return setupPrometheus().Run(":2112")
}
