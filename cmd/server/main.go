package main

import (
	"fmt"
	"github.com/FogusB/metrics-alerts-svc/internal/handlers"
	"github.com/FogusB/metrics-alerts-svc/internal/storages"
	"net/http"
)

type Storage interface {
	UpdateMetric(name string, mType storages.MetricType, value storages.MetricValue)
}

func main() {
	var storage Storage = storages.NewMemStorage()
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", handlers.PostHandler(storage))

	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
