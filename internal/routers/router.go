// Package routers содержит настройки маршрутизатора для веб-сервера.
package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/FogusB/metrics-alerts-svc/internal/handlers"
	"github.com/FogusB/metrics-alerts-svc/internal/middleware"
)

// Run запускает HTTP-сервер с заданными маршрутами.
func Run(metricHandler *handlers.MetricHandler, addr string) {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.Use(middleware.RequestLogger(logger))
	router.LoadHTMLGlob("templates/*")
	ServeRoutes(router, metricHandler)

	fullAddr := "http://" + addr
	zap.L().Info("Listening and serving HTTP on " + fullAddr)
	if err := router.Run(addr); err != nil {
		zap.L().Fatal("Could not start server", zap.Error(err))
	}
}

// ServeRoutes настраивает маршруты для веб-сервера.
func ServeRoutes(router *gin.Engine, metricHandler *handlers.MetricHandler) {
	router.POST("/update/", middleware.CheckContentTypeMiddleware("application/json", http.MethodPost), metricHandler.RestUpdateMetricValue)
	router.POST("/update/:type/:name/:value", middleware.CheckContentTypeMiddleware("text/plain", http.MethodPost), metricHandler.UpdateMetricValue)
	router.GET("/value/:type/:name", middleware.CheckContentTypeMiddleware("text/plain", http.MethodGet), metricHandler.GetMetricValue)
	router.GET("/", middleware.CheckContentTypeMiddleware("text/plain", http.MethodGet), metricHandler.GetAllMetrics)
}
