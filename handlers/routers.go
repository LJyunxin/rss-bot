package handlers

import "github.com/gin-gonic/gin"

func GinStart() *gin.Engine {
	e := gin.Default()
	e.GET("/subscription", GetSubscription)
	e.POST("/subscription", AddSubscription)
	e.DELETE("/subscription", DeleteSubscription)
	e.PUT("/webhooks", UpdateWebhook)
	return e
}
