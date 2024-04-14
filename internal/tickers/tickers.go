package tickers

import (
	"time"

	"go.uber.org/zap"

	"github.com/FogusB/metrics-alerts-svc/internal/collects"
	"github.com/FogusB/metrics-alerts-svc/internal/models"
	"github.com/FogusB/metrics-alerts-svc/internal/senders"
)

func Tickers(runAddress string, pollInterval time.Duration, reportInterval time.Duration) {
	zap.L().Sugar().Info("Agent started")
	zap.L().Sugar().Infof("Server address: %s", runAddress)
	zap.L().Sugar().Infof("Report interval: %v", reportInterval)
	zap.L().Sugar().Infof("Poll interval: %v", pollInterval)

	tickerPoll := time.NewTicker(pollInterval)
	tickerReport := time.NewTicker(reportInterval)
	var metricsToSend []models.Metrics
	var pollCount int64 = 0

	go func() {
		for {
			select {
			case <-tickerPoll.C:
				zap.L().Sugar().Infof("==============Metrics====================")
				collectedMetrics := collects.CollectMetrics()
				for key, value := range collectedMetrics {
					metric := models.Metrics{ID: key}
					switch v := value.(type) {
					case float64:
						metric.MType = "gauge"
						metric.Value = &v
					case uint64:
						metric.MType = "counter"
						temp := int64(v)
						metric.Delta = &temp
					case uint32:
						metric.MType = "counter"
						temp := int64(v)
						metric.Delta = &temp
					default:
						zap.L().Sugar().Errorf("Unsupported metric type for key %s: %T", key, v)
						continue
					}
					metricsToSend = append(metricsToSend, metric)
					if metric.MType == "gauge" {
						zap.L().Sugar().Infof("Processed metric - ID: %s, Type: %s, Value: %f", metric.ID, metric.MType, *metric.Value)
					} else {
						zap.L().Sugar().Infof("Processed metric - ID: %s, Type: %s, Delta: %d", metric.ID, metric.MType, *metric.Delta)
					}
				}
				// Увеличиваем PollCount
				pollCount++
				metricPollCount := models.Metrics{ID: "PollCount", MType: "counter", Delta: &pollCount}
				metricsToSend = append(metricsToSend, metricPollCount)
				zap.L().Sugar().Infof("Processed metric - ID: %s, Type: %s, Delta: %d", metricPollCount.ID, metricPollCount.MType, *metricPollCount.Delta)
				zap.L().Sugar().Infof("=========================================")

			case <-tickerReport.C:
				zap.L().Info("Sending metrics...")
				err := senders.SendMetrics(metricsToSend, runAddress)
				if err != nil {
					zap.L().Sugar().Errorf("Error sending metrics - %s", err)
				}
			}
		}
	}()

	select {} // Бесконечный цикл
}
