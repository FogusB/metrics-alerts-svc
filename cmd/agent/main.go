package main

import (
	"fmt"
	"github.com/FogusB/metrics-alerts-svc/internal/collectMetrics"
	"github.com/FogusB/metrics-alerts-svc/internal/senders"
	"time"
)

func main() {
	serverAddress := "http://localhost:8080"
	pollInterval := 2 * time.Second
	reportInterval := 10 * time.Second

	tickerPoll := time.NewTicker(pollInterval)
	tickerReport := time.NewTicker(reportInterval)

	metrics := make(map[string]interface{})

	go func() {
		for {
			select {
			case <-tickerPoll.C:
				for key, value := range collectMetrics.CollectMetrics() {
					metrics[key] = value
				}
				// Увеличиваем PollCount
				if val, ok := metrics["PollCount"].(int64); ok {
					metrics["PollCount"] = val + 1
				} else {
					metrics["PollCount"] = int64(1)
				}
			case <-tickerReport.C:
				fmt.Println("==============Metrics====================")
				for key, value := range metrics {
					if key == "GCCPUFraction" || key == "RandomValue" {
						fmt.Printf("%s : %f\n", key, value)
					} else {
						fmt.Printf("%s : %d\n", key, value)
					}
				}
				fmt.Println("=========================================")
				err := senders.SendMetrics(metrics, serverAddress)
				if err != nil {
					fmt.Printf("Error SendMetrics - %s", err)
				}
			}
		}
	}()

	select {} // Бесконечный цикл
}
