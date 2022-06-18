package api

import (
	"github.com/Coflnet/auction-stats/internal/model"
	"github.com/Coflnet/auction-stats/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// ListByNotifiersUser list the notifiers by a user
// @Summary lists the notifiers of a user
// @Schemes http
// @Tags notifiers
// @Produce json
// @Param        userId   path      int  true  "User ID"
// @Router /notifier/{userId} [get]
func ListByNotifiersUser(c *gin.Context) {
	p := c.Param("userId")

	if p == "" {
		c.JSON(400, gin.H{
			"error": "userId is required",
		})
		return
	}

	userId, err := strconv.Atoi(p)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid userId, must be an integer",
		})
		return
	}

	data, err := usecase.ListUserNotifiers(userId)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "an internal error occurred",
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

// CreateNotifier create a notifier
// @Summary creates a notifier
// @Description states, next / last evaluation will be ignored
// @Schemes http
// @Tags notifiers
// @Produce json
// @Accept json
// @Param notifier body model.Notifier true "Notifier"
// @Router /notifier [post]
func CreateNotifier(c *gin.Context) {

	var notifier *model.Notifier
	err := c.ShouldBindJSON(&notifier)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request",
		})
		return
	}

	notifier.LastEvaluation = time.Time{}
	notifier.NextEvaluation = time.Time{}
	notifier.NotifierStates = make([]*model.NotifierState, 0)

	err = usecase.CreateNotifier(notifier)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "an internal error occurred",
		})
		return
	}

	c.JSON(http.StatusOK, notifier)
}

// UpdateNotifier updates a notifier
// @Summary updates a notifier
// @Description states, next / last evaluation will be ignored, replaces the notifier with the same ID
// @Schemes http
// @Tags notifiers
// @Produce json
// @Accept json
// @Param notifier body model.Notifier true "Notifier"
// @Router /notifier [put]
func UpdateNotifier(c *gin.Context) {
	var notifier *model.Notifier
	err := c.Bind(&notifier)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request",
		})
		return
	}

	notifier.LastEvaluation = time.Time{}
	notifier.NextEvaluation = time.Time{}
	notifier.NotifierStates = make([]*model.NotifierState, 0)

	err = usecase.UpdateNotifier(notifier)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "an internal error occurred",
		})
		return
	}

	c.JSON(http.StatusOK, notifier)
}

// DeleteNotifier deletes a notifier
// @Summary deletes a notifier with a specifc ID
// @Schemes http
// @Tags notifiers
// @Produce json
// @Accept json
// @Param notifier body model.Notifier true "Notifier"
// @Router /notifier [delete]
func DeleteNotifier(c *gin.Context) {

	var notifier *model.Notifier
	err := c.Bind(&notifier)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request",
		})
		return
	}

	err = usecase.DeleteNotifier(notifier)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "an internal error occurred",
		})
		return
	}

	c.JSON(http.StatusOK, notifier)
}
