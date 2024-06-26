package tickers

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/FogusB/metrics-alerts-svc/internal/collects"
	"github.com/FogusB/metrics-alerts-svc/internal/senders"
)

func Tickers(runAddress string, pollInterval time.Duration, reportInterval time.Duration) {
	log.Infof("Server address: %s\n", runAddress)
	log.Infof("Report interval: %v\n", reportInterval)
	log.Infof("Poll interval: %v\n", pollInterval)

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
