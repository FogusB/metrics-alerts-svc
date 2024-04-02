package tickers

import (
	"time"

	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"

	"github.com/FogusB/metrics-alerts-svc/internal/collects"
	"github.com/FogusB/metrics-alerts-svc/internal/senders"
)

func Tickers(runAddress string, pollInterval time.Duration, reportInterval time.Duration) {
	zap.L().Sugar().Info("Agent started")
	zap.L().Sugar().Infof("Server address: %s", runAddress)
	zap.L().Sugar().Infof("Report interval: %v", reportInterval)
	zap.L().Sugar().Infof("Poll interval: %v", pollInterval)

	tickerPoll := time.NewTicker(pollInterval)
	tickerReport := time.NewTicker(reportInterval)
	metrics := make(map[string]interface{})

	go func() {
		for {
			select {
			case <-tickerPoll.C:
				for key, value := range collects.CollectMetrics() {
					metrics[key] = value
				}
				// Увеличиваем PollCount
				if val, ok := metrics["PollCount"].(int64); ok {
					metrics["PollCount"] = val + 1
				} else {
					metrics["PollCount"] = int64(1)
				}
			case <-tickerReport.C:
				log.Info("==============Metrics====================")
				for key, value := range metrics {
					if key == "GCCPUFraction" || key == "RandomValue" {
						log.Infof("%s : %f\n", key, value)
					} else {
						log.Infof("%s : %d\n", key, value)
					}
				}
				log.Info("=========================================")
				err := senders.SendMetrics(metrics, runAddress)
				if err != nil {
					log.Errorf("Error SendMetrics - %s", err)
				}
			}
		}
	}()

	select {} // Бесконечный цикл
}
