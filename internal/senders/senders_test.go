package senders

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FogusB/metrics-alerts-svc/internal/models"
)

func TestSendMetrics(t *testing.T) {
	// Создание тестового сервера
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверка метода запроса
		if r.Method != http.MethodPost {
			t.Errorf("Expected 'POST' request, got '%s'", r.Method)
		}
		_, err := fmt.Fprintln(w, "OK")
		if err != nil {
			return
		}
	}))
	defer testServer.Close()

	tempFloat64 := 5.5
	tempInt64 := int64(55)
	tests := []struct {
		name    string
		metrics []models.Metrics
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid gauge metric",
			metrics: []models.Metrics{
				{"memory", "gauge", nil, &tempFloat64},
				{"cpu", "gauge", nil, &tempFloat64},
			},
			wantErr: false,
		},
		{
			name: "valid counter metric",
			metrics: []models.Metrics{
				{"disk", "counter", &tempInt64, nil},
				{"requests", "counter", &tempInt64, nil},
			},
			wantErr: false,
		},
		{
			name: "valid all metric",
			metrics: []models.Metrics{
				{"PollCount", "counter", &tempInt64, nil},
				{"NumForcedGC", "counter", &tempInt64, nil},
				{"HeapSys", "counter", &tempInt64, nil},
				{"MSpanInuse", "counter", &tempInt64, nil},
				{"HeapReleased", "counter", &tempInt64, nil},
				{"HeapIdle", "counter", &tempInt64, nil},
				{"Sys", "counter", &tempInt64, nil},
				{"BuckHashSys", "counter", &tempInt64, nil},
				{"NextGC", "counter", &tempInt64, nil},
				{"MSpanSys", "counter", &tempInt64, nil},
				{"PauseTotalNs", "counter", &tempInt64, nil},
				{"HeapAlloc", "counter", &tempInt64, nil},
				{"StackInuse", "counter", &tempInt64, nil},
				{"HeapObjects", "counter", &tempInt64, nil},
				{"GCCPUFraction", "gauge", nil, &tempFloat64},
				{"RandomValue", "gauge", nil, &tempFloat64},
				{"MCacheSys", "counter", &tempInt64, nil},
				{"OtherSys", "counter", &tempInt64, nil},
				{"Alloc", "counter", &tempInt64, nil},
				{"Lookups", "counter", &tempInt64, nil},
				{"LastGC", "counter", &tempInt64, nil},
				{"GCSys", "counter", &tempInt64, nil},
				{"Mallocs", "counter", &tempInt64, nil},
				{"TotalAlloc", "counter", &tempInt64, nil},
				{"HeapInuse", "counter", &tempInt64, nil},
				{"MCacheInuse", "counter", &tempInt64, nil},
				{"Frees", "counter", &tempInt64, nil},
				{"StackSys", "counter", &tempInt64, nil},
				{"NumGC", "counter", &tempInt64, nil},
			},
			wantErr: false,
		},
		{
			name: "unsupported metric type",
			metrics: []models.Metrics{
				{"usersOnline", "usersOnline", &tempInt64, nil},
			},
			wantErr: false, // Функция продолжает работу, несмотря на неподдерживаемый тип
		},
		{
			name:    "empty metrics",
			metrics: []models.Metrics{},
			wantErr: false, // Пустые метрики не вызывают ошибку
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SendMetrics(tt.metrics, testServer.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
