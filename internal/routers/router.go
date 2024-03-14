// Package routers содержит настройки маршрутизатора для веб-сервера.
package routers

import (
	"fmt"

	"github.com/FogusB/metrics-alerts-svc/internal/handlers"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Run запускает HTTP-сервер с заданными маршрутами.
func Run(metricHandler *handlers.MetricHandler) {
	router := gin.Default()

	// Загрузка HTML-шаблонов
	router.LoadHTMLGlob("templates/*")

	ServeRoutes(router, metricHandler)

	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	fmt.Println("Server is running on http://localhost:8080")
}

// ServeRoutes настраивает маршруты для веб-сервера.
func ServeRoutes(router *gin.Engine, metricHandler *handlers.MetricHandler) {
	router.POST("/update/:type/:name/:value", metricHandler.UpdateMetricValue)
	router.GET("/value/:type/:name", metricHandler.GetMetricValue)
	router.GET("/", metricHandler.GetAllMetrics)
}
