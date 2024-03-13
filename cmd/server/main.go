package main

import (
	"flag"
	"github.com/FogusB/metrics-alerts-svc/internal/handlers"
	"github.com/FogusB/metrics-alerts-svc/internal/routers"
	"github.com/FogusB/metrics-alerts-svc/internal/storages"
)

func main() {
	var addr string
	flag.StringVar(&addr, "a", "localhost:8080", "HTTP server address")
	flag.Parse()
	memStorage := storages.NewMemStorage()
	metricHandler := &handlers.MetricHandler{Storage: memStorage}
	routers.Run(metricHandler, addr)
}
