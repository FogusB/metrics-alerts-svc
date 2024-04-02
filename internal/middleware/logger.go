package middleware

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GlobalLogger() {
	// Настройка конфигурации энкодера
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // This line is illustrative and won't work without further customization

	// Настройка ядра логгера
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.Lock(os.Stdout),
		zap.NewAtomicLevelAt(zap.DebugLevel),
	)

	// Инициализация логгера с настроенным ядром
	logger := zap.New(core)
	zap.ReplaceGlobals(logger)

	zap.L().Info("Global logger initialized")
}

func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Запуск таймера
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Обработка запроса
		c.Next()

		// Подсчет времени
		end := time.Now()
		duration := end.Sub(start)

		// Логирование запроса
		logger.Info("request",
			zap.String("path", path),
			zap.String("method", method),
			zap.Duration("duration", duration),
		)

		// Логирование ответа
		statusCode := c.Writer.Status()
		responseSize := c.Writer.Size()
		logger.Info("response",
			zap.Int("status", statusCode),
			zap.Int("size", responseSize),
		)
	}
}
