package main

//go:generate go build -o=server

import (
	"flag"
	"github.com/FogusB/metrics-alerts-svc/internal/handlers"
	"github.com/FogusB/metrics-alerts-svc/internal/routers"
	"github.com/FogusB/metrics-alerts-svc/internal/storages"
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	AddressEnv  string `env:"ADDRESS"`
	AddressFlag string
}

func parseFlags() string {
	var cfg Config

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Info(cfg)

	flag.StringVar(&cfg.AddressFlag, "a", ":8080", "HTTP server address")
	flag.Parse()

	if cfg.AddressEnv != "" {
		return cfg.AddressEnv
	} else {
		return cfg.AddressFlag
	}
}

func main() {
	memStorage := storages.NewMemStorage()
	metricHandler := &handlers.MetricHandler{Storage: memStorage}
	runAddress := parseFlags()
	routers.Run(metricHandler, runAddress)
}
