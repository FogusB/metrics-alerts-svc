package main

import (
	"github.com/FogusB/metrics-alerts-svc/internal/handlers"
	"github.com/FogusB/metrics-alerts-svc/internal/routers"
	"github.com/FogusB/metrics-alerts-svc/internal/storages"
)

func main() {
	memStorage := storages.NewMemStorage()
	metricHandler := &handlers.MetricHandler{Storage: memStorage}
	routers.Run(metricHandler)
}
