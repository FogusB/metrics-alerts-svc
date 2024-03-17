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
	var serverAddr string
	urlSchema := "http://"
	var reportInterval, pollInterval int

	flag.StringVar(&serverAddr, "a", "127.0.0.1:8080", "HTTP server address")
	flag.IntVar(&reportInterval, "r", 10, "Report interval (s)")
	flag.IntVar(&pollInterval, "p", 2, "Poll interval (s)")
	flag.Parse()

	log.Infof("Server address: %s\n", urlSchema+serverAddr)
	log.Infof("Report interval: %v\n", time.Duration(reportInterval)*time.Second)
	log.Infof("Poll interval: %v\n", time.Duration(pollInterval)*time.Second)

	tickerPoll := time.NewTicker(time.Duration(reportInterval) * time.Second)
	tickerReport := time.NewTicker(time.Duration(pollInterval) * time.Second)

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
				err := senders.SendMetrics(metrics, urlSchema+serverAddr)
				if err != nil {
					log.Errorf("Error SendMetrics - %s", err)
				}
			}
		}
	}()

	select {} // Бесконечный цикл
}
