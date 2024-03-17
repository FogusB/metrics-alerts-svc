// Package routers содержит настройки маршрутизатора для веб-сервера.
package routers

import (
	"github.com/FogusB/metrics-alerts-svc/internal/handlers"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Run запускает HTTP-сервер с заданными маршрутами.
func Run(metricHandler *handlers.MetricHandler, addr string) {
	router := gin.Default()

	// Загрузка HTML-шаблонов
	router.LoadHTMLGlob("templates/*")

	ServeRoutes(router, metricHandler)

	log.Printf("Starting server on %s\n", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

// ServeRoutes настраивает маршруты для веб-сервера.
func ServeRoutes(router *gin.Engine, metricHandler *handlers.MetricHandler) {
	router.POST("/update/:type/:name/:value", metricHandler.UpdateMetricValue)
	router.GET("/value/:type/:name", metricHandler.GetMetricValue)
	router.GET("/", metricHandler.GetAllMetrics)
}
