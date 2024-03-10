package main

import (
	"fmt"
	"github.com/FogusB/metrics-alerts-svc/internal"
	"net/http"
	"strconv"
	"strings"
)

type Storage interface {
	UpdateMetric(name string, mType internal.MetricType, value internal.MetricValue)
}

func postHandler(storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// этот обработчик принимает только запросы, отправленные методом POST и ContentType только text/plain
		if r.Method != http.MethodPost {
			http.Error(w, "Only Post requests are allowed!", http.StatusMethodNotAllowed)
			return
		}
		for k, v := range r.Header {
			if k == "Content-Type" && v[0] != "text/plain" {
				http.Error(w, "Only text/plain contents are allowed!", http.StatusUnsupportedMediaType)
				return
			}
		}
		parts := strings.Split(r.URL.Path, "/")
		// Ожидаемый формат: /update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
		if len(parts) != 5 {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		mType := internal.MetricType(parts[2])
		name := parts[3]
		rawValue := parts[4]

		if mType != internal.Gauge && mType != internal.Counter {
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}
		if name == "" {
			http.Error(w, "Metric name is required", http.StatusNotFound)
			return
		}

		var value internal.MetricValue
		var err error
		if mType == internal.Gauge {
			value.GaugeValue, err = strconv.ParseFloat(rawValue, 64)
			if err != nil {
				http.Error(w, "Invalid value for gauge", http.StatusBadRequest)
				return
			}
			fmt.Printf("Gauge type: %s, value: %v\n", name, value)
		} else {
			value.CounterValue, err = strconv.ParseInt(rawValue, 10, 64)
			if err != nil {
				http.Error(w, "Invalid value for counter", http.StatusBadRequest)
				return
			}
			fmt.Printf("Counter type: %s, value: %v\n", name, value)
		}

		storage.UpdateMetric(name, mType, value)
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	var storage Storage = internal.NewMemStorage()
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", postHandler(storage))

	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
