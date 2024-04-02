package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CheckContentTypeMiddleware(expectedContentType string, expectedMethod string) gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.GetHeader("Content-Type")
		if contentType != expectedContentType && contentType != "" {
			c.String(http.StatusUnsupportedMediaType, "Unsupported Media Type")
			zap.L().Warn("Не поддерживаемый тип контента")
			c.Abort()
			return
		}
		if c.Request.Method != expectedMethod {
			c.String(http.StatusMethodNotAllowed, "Method Not Allowed")
			zap.L().Warn("Метод не поддерживается")
			c.Abort()
			return
		}
		c.Next()
	}
}
