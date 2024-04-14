package main

//go:generate go build -o=../../bin/server

import (
	"github.com/FogusB/metrics-alerts-svc/internal/flags"
	"github.com/FogusB/metrics-alerts-svc/internal/handlers"
	"github.com/FogusB/metrics-alerts-svc/internal/middleware"
	"github.com/FogusB/metrics-alerts-svc/internal/routers"
	"github.com/FogusB/metrics-alerts-svc/internal/storages"
)

func main() {
	middleware.GlobalLogger()
	memStorage := storages.NewMemStorage()
	metricHandler := &handlers.MetricHandler{Storage: memStorage}
	runAddress, _, _ := flags.ParseFlags("server")
	routers.Run(metricHandler, runAddress)
}
