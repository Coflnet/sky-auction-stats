package api

import (
	_ "github.com/Coflnet/auction-stats/docs"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func StartApi() error {
	r := setupRouter()
	return r.Run()
}

// @title Auction Stats API
// @version 1.0
// @description API for Auction Stats Service

// @contact.name Flou21
// @contact.email muehlhans.f@coflnet.com

// @license.name AGPL v3

// @host localhost:8080
// @BasePath /api/
func setupRouter() *gin.Engine {
	r := gin.Default()

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// auctions counts
	r.GET("/api/new-auctions", NewAuctions)

	// flippers counts

	// notifiers
	r.GET("/api/notifier/:userId", ListByNotifiersUser)
	r.POST("/api/notifier", CreateNotifier)
	r.PUT("/api/notifier", UpdateNotifier)
	r.DELETE("/api/notifier", DeleteNotifier)

	return r
}
