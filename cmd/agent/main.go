package main

//go:generate go build -o=agent

import (
	"flag"
	"github.com/FogusB/metrics-alerts-svc/internal/collects"
	"github.com/FogusB/metrics-alerts-svc/internal/senders"
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
	"net/url"
	"time"
)

type Config struct {
	AddressEnv        url.URL `env:"ADDRESS"`
	AddressSrv        string
	ReportIntervalEnv int `env:"REPORT_INTERVAL"`
	ReportInterval    int
	PollIntervalEnv   int `env:"POLL_INTERVAL"`
	PollInterval      int
}

func parseFlags() (string, time.Duration, time.Duration) {
	var cfg Config
	urlSchema := "http://"

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Info(cfg)

	flag.StringVar(&cfg.AddressSrv, "a", "localhost:8080", "HTTP server address")
	flag.IntVar(&cfg.ReportInterval, "r", 10, "Report interval (s)")
	flag.IntVar(&cfg.PollInterval, "p", 2, "Poll interval (s)")
	flag.Parse()

	if cfg.AddressEnv.String() != "" {
		cfg.AddressSrv = cfg.AddressEnv.String()
	}
	cfg.AddressSrv = urlSchema + cfg.AddressSrv

	if cfg.ReportIntervalEnv != 0 {
		cfg.ReportInterval = cfg.ReportIntervalEnv
	}

	if cfg.PollIntervalEnv != 0 {
		cfg.PollInterval = cfg.PollIntervalEnv
	}

	return cfg.AddressSrv, time.Duration(cfg.ReportInterval) * time.Second, time.Duration(cfg.PollInterval) * time.Second
}

func main() {
	runAddress, reportInterval, pollInterval := parseFlags()

	log.Infof("Server address: %s\n", runAddress)
	log.Infof("Report interval: %v\n", reportInterval)
	log.Infof("Poll interval: %v\n", pollInterval)

	tickerPoll := time.NewTicker(reportInterval)
	tickerReport := time.NewTicker(pollInterval)

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
