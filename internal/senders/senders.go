package senders

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FogusB/metrics-alerts-svc/internal/models"
)

// SendMetrics отправляет метрики на сервер
func SendMetrics(metrics []models.Metrics, serverAddress string) error {
	client := &http.Client{}
	var request *http.Request
	for _, metric := range metrics {
		jsonData, err := json.Marshal(metric)
		if err != nil {
			return fmt.Errorf("marshaling json for key %v: %w", metric, err)
		}
		request, err = http.NewRequest("POST", fmt.Sprintf("%s/update/", serverAddress), bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("creating request for key %v: %w", metric, err)
		}

		request.Header.Set("Content-Type", "application/json")

		response, err := client.Do(request)
		if err != nil {
			return fmt.Errorf("sending request for key %v: %w", metric, err)
		}
		err = response.Body.Close()
		if err != nil {
			return fmt.Errorf("closing response body for key %v: %w", metric, err)
		}
	}
	return nil
}
