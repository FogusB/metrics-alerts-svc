package main

import (
	"bytes"
	"github.com/FogusB/metrics-alerts-svc/internal/handlers"
	"github.com/FogusB/metrics-alerts-svc/internal/storages"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MetricType string
type MetricValue float64

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) UpdateMetric(name string, mType MetricType, value MetricValue) {
	m.Called(name, mType, value)
}

func TestUpdateMetric(t *testing.T) {
	mockStorage := new(MockStorage)
	mockStorage.On("UpdateMetric", "testMetric", MetricType("test"), MetricValue(100))
	mockStorage.UpdateMetric("testMetric", "test", MetricValue(100))
	mockStorage.AssertExpectations(t)
}

func TestPostHandler(t *testing.T) {
	storage := storages.NewMemStorage()
	handler := handlers.PostHandler(storage)

	t.Run("ValidRequest", func(t *testing.T) {
		body := bytes.NewBufferString(`valid`)
		req, err := http.NewRequest("POST", "/update/counter/testMetric/100", body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("InvalidRequest", func(t *testing.T) {
		body := bytes.NewBufferString(`invalid`)
		req, err := http.NewRequest("POST", "/update/", body)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code for invalid request: got %v want %v", status, http.StatusBadRequest)
		}
	})
}
