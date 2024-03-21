// Package routers содержит настройки маршрутизатора для веб-сервера.
package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/FogusB/metrics-alerts-svc/internal/handlers"
)

func CheckContentTypeMiddleware(expectedContentType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if contentType := c.GetHeader("Content-Type"); contentType != expectedContentType {
			c.String(http.StatusUnsupportedMediaType, "Unsupported Media Type")
			log.Warning("unsupported media type")
			c.Abort()
			return
		}
		c.Next()
	}
}

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
	router.POST("/update/:type/:name/:value", CheckContentTypeMiddleware("text/plain"), metricHandler.UpdateMetricValue)
	router.GET("/value/:type/:name", CheckContentTypeMiddleware("text/plain"), metricHandler.GetMetricValue)
	router.GET("/", CheckContentTypeMiddleware("text/plain"), metricHandler.GetAllMetrics)
}
