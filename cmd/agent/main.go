package main

//go:generate go build -o=agent

import (
	"flag"
	"github.com/FogusB/metrics-alerts-svc/internal/collects"
	"github.com/FogusB/metrics-alerts-svc/internal/senders"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	var serverAddrFlag string
	url_schema := "http://"
	var reportInterval, pollInterval time.Duration

	flag.StringVar(&serverAddrFlag, "a", "localhost:8080", "HTTP server address")
	flag.DurationVar(&reportInterval, "r", 10*time.Second, "Report interval (s)")
	flag.DurationVar(&pollInterval, "p", 2*time.Second, "Poll interval (s)")
	flag.Parse()

	serverAddr := url_schema + serverAddrFlag

	log.Infof("Server address: %s\n", serverAddr)
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
				err := senders.SendMetrics(metrics, serverAddr)
				if err != nil {
					log.Errorf("Error SendMetrics - %s", err)
				}
			}
		}
	}()

	select {} // Бесконечный цикл
}
