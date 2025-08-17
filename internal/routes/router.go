package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	handler "github.com/pulkit2910-bit/rate-limiter-service/internal/handlers"
)

func SetupRouter(handler handler.LimiterHandler) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"message": "Service is running",
		})
	})

	router.GET("/check", handler.CheckHandler)
	router.POST("/config", handler.ConfigHandler)

	return router
}