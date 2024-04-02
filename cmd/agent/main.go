package main

//go:generate go build -o=../../bin/agent

import (
	"github.com/FogusB/metrics-alerts-svc/internal/flags"
	"github.com/FogusB/metrics-alerts-svc/internal/middleware"
	"github.com/FogusB/metrics-alerts-svc/internal/tickers"
)

func main() {
	middleware.GlobalLogger()
	runAddress, reportInterval, pollInterval := flags.ParseFlags("agent")
	tickers.Tickers(runAddress, reportInterval, pollInterval)
}
