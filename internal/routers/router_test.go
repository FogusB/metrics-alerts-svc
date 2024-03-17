package routers

import (
	"github.com/FogusB/metrics-alerts-svc/internal/storages"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type MockMetricHandler struct {
	mock.Mock
}

func (m *MockMetricHandler) UpdateMetricValue(name string, mType storages.MetricType, value storages.MetricValue) error {
	args := m.Called(name, mType, value)
	return args.Error(0)
}

func (m *MockMetricHandler) GetMetricValue(name string) (storages.MetricValue, bool) {
	args := m.Called(name)
	return args.Get(0).(storages.MetricValue), args.Bool(1)
}

func (m *MockMetricHandler) GetAllMetrics() (map[string]storages.MetricValue, error) {
	args := m.Called()
	return args.Get(0).(map[string]storages.MetricValue), args.Error(1)
}

func TestServeRoutes(t *testing.T) {
	mockHandler := new(MockMetricHandler)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.LoadHTMLGlob("../../templates/*")

	t.Run("TestUpdateRoute", func(t *testing.T) {
		mockHandler.On("UpdateMetricValue", "cpu_usage", storages.Gauge, storages.MetricValue{GaugeValue: 42.0}).Return(nil)

		r.POST("/update/:type/:name/:value", func(c *gin.Context) {
			mType := storages.MetricType(c.Param("type"))
			name := c.Param("name")
			rawValue := c.Param("value")
			var value storages.MetricValue
			var err error
			value.GaugeValue, err = strconv.ParseFloat(rawValue, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": "invalid value"})
				return
			}

			err = mockHandler.UpdateMetricValue(name, mType, value)
			if err != nil {
				log.Error(err)
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating metric"})
				return
			}
			c.JSON(http.StatusOK, value)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/update/gauge/cpu_usage/42.0", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockHandler.AssertExpectations(t)
	})

	t.Run("TestValueRoute", func(t *testing.T) {
		mockHandler.On("GetMetricValue", "cpu_usage").Return(storages.MetricValue{GaugeValue: 42.0}, true)

		r.GET("/value/:type/:name", func(c *gin.Context) {
			name := c.Param("name")
			value, found := mockHandler.GetMetricValue(name)
			if !found {
				c.JSON(http.StatusNotFound, gin.H{"message": "metric not found"})
				return
			}
			c.JSON(http.StatusOK, value)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/value/gauge/cpu_usage", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockHandler.AssertExpectations(t)
	})

	t.Run("TestRootRoute", func(t *testing.T) {
		mockHandler.On("GetAllMetrics").Return(map[string]storages.MetricValue{"cpu_usage": {GaugeValue: 42.0}}, nil)

		r.GET("/", func(c *gin.Context) {
			metrics, err := mockHandler.GetAllMetrics()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Error getting metrics"})
				return
			}

			c.HTML(http.StatusOK, "metrics.tmpl", gin.H{
				"metrics": metrics,
			})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockHandler.AssertExpectations(t)
	})

}
