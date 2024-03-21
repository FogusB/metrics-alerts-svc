// Package routers содержит настройки маршрутизатора для веб-сервера.
package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/FogusB/metrics-alerts-svc/internal/handlers"
)

func CheckContentTypeMiddleware(expectedContentType string, expectedMethod string) gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.GetHeader("Content-Type")
		if contentType != expectedContentType && contentType != "" {
			c.String(http.StatusUnsupportedMediaType, "Unsupported Media Type")
			log.Warning("Не поддерживаемый тип контента")
			c.Abort()
			return
		}
		if c.Request.Method != expectedMethod {
			c.String(http.StatusMethodNotAllowed, "Method Not Allowed")
			log.Warning("Метод не поддерживается")
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
	router.POST("/update/:type/:name/:value", CheckContentTypeMiddleware("text/plain", http.MethodPost), metricHandler.UpdateMetricValue)
	router.GET("/value/:type/:name", CheckContentTypeMiddleware("text/plain", http.MethodGet), metricHandler.GetMetricValue)
	router.GET("/", CheckContentTypeMiddleware("text/plain", http.MethodGet), metricHandler.GetAllMetrics)
}
