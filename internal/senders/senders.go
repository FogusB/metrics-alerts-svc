package senders

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func SendMetrics(metrics map[string]interface{}, serverAddress string) error {
	client := &http.Client{}
	for key, value := range metrics {
		var request *http.Request
		var err error
		//fmt.Printf("type: %T\n", value)
		switch v := value.(type) {
		case float64:
			request, err = http.NewRequest("POST", fmt.Sprintf("%s/update/gauge/%s/%f", serverAddress, key, v), nil)
		case uint64, uint32:
			request, err = http.NewRequest("POST", fmt.Sprintf("%s/update/gauge/%s/%d", serverAddress, key, v), nil)
		case int64:
			request, err = http.NewRequest("POST", fmt.Sprintf("%s/update/counter/%s/%d", serverAddress, key, v), nil)
		default:
			zap.L().Sugar().Warnf("Broken key: %s, broken value: %v, type: %T\n", key, v, v)
			continue // В этом случае продолжаем, так как это не критическая ошибка для всей операции
		}

		if err != nil {
			return fmt.Errorf("error creating request for key %s: %w", key, err)
		}

		request.Header.Set("Content-Type", "text/plain")

		response, err := client.Do(request)
		if err != nil {
			return fmt.Errorf("error sending request for key %s: %w", key, err)
		}
		err = response.Body.Close()
		if err != nil {
			return fmt.Errorf("error closing response body for key %s: %w", key, err)
		}
	}

	return nil // Возвращаем nil, если функция выполнена без ошибок
}
