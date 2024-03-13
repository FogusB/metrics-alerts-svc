package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FogusB/metrics-alerts-svc/internal/storages"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) UpdateMetric(name string, mType storages.MetricType, value storages.MetricValue) error {
	args := m.Called(name, mType, value)
	return args.Error(0)
}

func (m *MockStorage) GetMetric(name string) (storages.MetricValue, bool) {
	args := m.Called(name)
	return args.Get(0).(storages.MetricValue), args.Bool(1)
}

func (m *MockStorage) GetAllMetrics() (map[string]storages.MetricValue, error) {
	args := m.Called()
	return args.Get(0).(map[string]storages.MetricValue), args.Error(1)
}

func TestUpdateMetricValue(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		contentType    string
		typeMetrics    string
		url            string
		body           io.Reader
		wantStatusCode int
		wantErr        bool
	}{
		{
			name:           "valid gauge update",
			method:         http.MethodPost,
			contentType:    "text/plain",
			typeMetrics:    "gauge",
			url:            "/update/gauge/test_metric/123.456",
			body:           nil,
			wantStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "valid counter update",
			method:         http.MethodPost,
			contentType:    "text/plain",
			typeMetrics:    "counter",
			url:            "/update/counter/test_metric/789",
			body:           nil,
			wantStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "invalid method",
			method:         http.MethodGet,
			contentType:    "text/plain",
			typeMetrics:    "gauge",
			url:            "/update/gauge/test_metric/123",
			body:           nil,
			wantStatusCode: http.StatusMethodNotAllowed,
			wantErr:        true,
		},
		{
			name:           "invalid content type",
			method:         http.MethodPost,
			contentType:    "application/json",
			typeMetrics:    "gauge",
			url:            "/update/gauge/test_metric/123",
			body:           bytes.NewBufferString("{\"value\":123}"),
			wantStatusCode: http.StatusUnsupportedMediaType,
			wantErr:        true,
		},
		{
			name:           "invalid metric type",
			method:         http.MethodPost,
			contentType:    "text/plain",
			typeMetrics:    "unknown",
			url:            "/update/unknown/test_metric/123",
			body:           nil,
			wantStatusCode: http.StatusBadRequest,
			wantErr:        true,
		},
	}

	mockStorage := new(MockStorage)
	handler := MetricHandler{Storage: mockStorage}

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/update/:type/:name/:value", handler.UpdateMetricValue)
	r.GET("/update/:type/:name/:value", handler.UpdateMetricValue)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.typeMetrics == "counter" {
				mockStorage.On("UpdateMetric", "test_metric", storages.Counter, mock.Anything).Return(nil)
			} else {
				mockStorage.On("UpdateMetric", "test_metric", storages.Gauge, mock.Anything).Return(nil)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.url, tt.body)
			req.Header.Set("Content-Type", tt.contentType)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatusCode, w.Code)
			mockStorage.AssertExpectations(t)
		})
	}
}

func TestGetMetricValue(t *testing.T) {
	mockStorage := new(MockStorage)
	handler := MetricHandler{Storage: mockStorage}

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/metric/:name", handler.GetMetricValue)

	t.Run("metric found", func(t *testing.T) {
		mockStorage.On("GetMetric", "test").Return(storages.MetricValue{GaugeValue: 123.45}, true)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/metric/test", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "123.45")
		mockStorage.AssertExpectations(t)
	})

	t.Run("metric not found", func(t *testing.T) {
		mockStorage.On("GetMetric", "unknown").Return(storages.MetricValue{}, false)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/metric/unknown", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestGetAllMetrics(t *testing.T) {
	mockStorage := new(MockStorage)
	handler := MetricHandler{Storage: mockStorage}

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.LoadHTMLGlob("../../templates/*")
	r.GET("/", handler.GetAllMetrics)

	t.Run("success", func(t *testing.T) {
		mockStorage.On("GetAllMetrics").Return(map[string]storages.MetricValue{"test": {GaugeValue: 123.45}}, nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockStorage.AssertExpectations(t)
	})
}
