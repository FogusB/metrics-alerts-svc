package handlers

import (
	"bytes"
	"github.com/FogusB/metrics-alerts-svc/internal/storages"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockStorage - моковая реализация интерфейса Storage для тестирования
type MockStorage struct {
	UpdateMetricFunc func(name string, mType storages.MetricType, value storages.MetricValue)
}

func (m *MockStorage) UpdateMetric(name string, mType storages.MetricType, value storages.MetricValue) {
	if m.UpdateMetricFunc != nil {
		m.UpdateMetricFunc(name, mType, value)
	}
}

// TestPostHandler - тестирование PostHandler
func TestPostHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		contentType    string
		url            string
		body           io.Reader
		wantStatusCode int
		wantErr        bool
	}{
		{
			name:           "valid gauge update",
			method:         http.MethodPost,
			contentType:    "text/plain",
			url:            "/update/gauge/test_metric/123.456",
			body:           nil,
			wantStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "valid counter update",
			method:         http.MethodPost,
			contentType:    "text/plain",
			url:            "/update/counter/test_metric/789",
			body:           nil,
			wantStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "invalid method",
			method:         http.MethodGet,
			contentType:    "text/plain",
			url:            "/update/gauge/test_metric/123",
			body:           nil,
			wantStatusCode: http.StatusMethodNotAllowed,
			wantErr:        true,
		},
		{
			name:           "invalid content type",
			method:         http.MethodPost,
			contentType:    "application/json",
			url:            "/update/gauge/test_metric/123",
			body:           bytes.NewBufferString("{\"value\":123}"),
			wantStatusCode: http.StatusUnsupportedMediaType,
			wantErr:        true,
		},
		{
			name:           "invalid metric type",
			method:         http.MethodPost,
			contentType:    "text/plain",
			url:            "/update/unknown/test_metric/123",
			body:           nil,
			wantStatusCode: http.StatusBadRequest,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.url, tt.body)
			request.Header.Set("Content-Type", tt.contentType)
			recorder := httptest.NewRecorder()

			mockStorage := &MockStorage{}
			handler := PostHandler(mockStorage)
			handler.ServeHTTP(recorder, request)

			if recorder.Code != tt.wantStatusCode {
				t.Errorf("expected status %v, got %v", tt.wantStatusCode, recorder.Code)
			}
		})
	}
}
