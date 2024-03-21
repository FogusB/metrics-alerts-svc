package senders

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
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

	tests := []struct {
		name    string
		metrics map[string]interface{}
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid gauge metric",
			metrics: map[string]interface{}{
				"cpu":    0.95,
				"memory": uint64(1024),
				"disk":   int64(100),
			},
			wantErr: false,
		},
		{
			name: "valid counter metric",
			metrics: map[string]interface{}{
				"requests": int64(100),
			},
			wantErr: false,
		},
		{
			name: "valid all metric",
			metrics: map[string]interface{}{
				"PollCount":     int64(1),
				"NumForcedGC":   uint64(0),
				"HeapSys":       uint64(7929856),
				"MSpanInuse":    uint64(91896),
				"HeapReleased":  uint64(5677056),
				"HeapIdle":      uint64(6774784),
				"Sys":           uint64(12618768),
				"BuckHashSys":   uint64(3980),
				"NextGC":        uint64(4194304),
				"MSpanSys":      uint64(114072),
				"PauseTotalNs":  uint64(73233000),
				"HeapAlloc":     uint64(530400),
				"StackInuse":    uint64(458752),
				"HeapObjects":   uint64(4368),
				"GCCPUFraction": 0.000006,
				"RandomValue":   0.391760,
				"MCacheSys":     uint64(15600),
				"OtherSys":      uint64(713036),
				"Alloc":         uint64(530400),
				"Lookups":       uint64(0),
				"LastGC":        uint64(1710153068469314000),
				"GCSys":         uint64(3383472),
				"Mallocs":       uint64(2051553),
				"TotalAlloc":    uint64(172110256),
				"HeapInuse":     uint64(1155072),
				"MCacheInuse":   uint64(9600),
				"Frees":         uint64(2047185),
				"StackSys":      uint64(458752),
				"NumGC":         uint64(125),
			},
			wantErr: false,
		},
		{
			name: "unsupported metric type",
			metrics: map[string]interface{}{
				"usersOnline": "500", // строка вместо числа
			},
			wantErr: false, // Функция продолжает работу, несмотря на неподдерживаемый тип
		},
		{
			name:    "empty metrics",
			metrics: map[string]interface{}{},
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
