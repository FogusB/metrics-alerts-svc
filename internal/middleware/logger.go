package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/FogusB/metrics-alerts-svc/internal/models"
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
		jsonDataBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			zap.L().Error("Error reading request body: ", zap.Error(err))
		}

		// Подсчет времени
		end := time.Now()
		duration := end.Sub(start)

		// Обработка запроса
		c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonDataBytes))
		c.Next()

		var jsonData models.Metrics
		err = json.Unmarshal(jsonDataBytes, &jsonData)
		if err != nil {
			zap.L().Error("Error unmarshalling request body: ", zap.Error(err))
		}

		// Логирование запроса
		logger.Info("request",
			zap.String("path", path),
			zap.String("method", method),
			zap.Any("body", jsonData),
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
