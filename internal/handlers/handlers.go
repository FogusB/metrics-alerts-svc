package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/FogusB/metrics-alerts-svc/internal/storages"
)

type MetricHandler struct {
	Storage Storage
}

type Storage interface {
	UpdateMetric(name string, mType storages.MetricType, value storages.MetricValue) error
	GetMetric(name string) (storages.Value, bool)
	GetAllMetrics() (map[string]storages.MetricValue, error)
}

// MetricUpdateRequest определяет структуру входящего запроса на обновление метрики.
type MetricUpdateRequest struct {
	Type  storages.MetricType `uri:"type" binding:"required,oneof=gauge counter"`
	Name  string              `uri:"name" binding:"required"`
	Value string              `uri:"value" binding:"required"`
}

func (h *MetricHandler) UpdateMetricValue(c *gin.Context) {
	var request MetricUpdateRequest
	println(request.Type)
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Error("Ошибка привязки URI: ", err)
		return
	}

	var value storages.MetricValue
	var parseErr error
	switch request.Type {
	case storages.Gauge:
		value.GaugeValue, parseErr = strconv.ParseFloat(request.Value, 64)
	case storages.Counter:
		value.CounterValue, parseErr = strconv.ParseUint(request.Value, 10, 64)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric type"})
		return
	}

	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid value for metric"})
		log.Error("Ошибка преобразования значения метрики: ", parseErr)
		return
	}
	if err := h.Storage.UpdateMetric(request.Name, request.Type, value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating metric"})
		log.Error("Ошибка обновления метрики: ", err)
		return
	}
	if c.Request.Method != http.MethodPost {
		c.String(http.StatusMethodNotAllowed, "Method Not Allowed")
		log.Warning("Метод не разрешен")
		return
	}

	c.Status(http.StatusOK)
}

func (h *MetricHandler) GetMetricValue(c *gin.Context) {
	name := c.Param("name")
	value, found := h.Storage.GetMetric(name)
	if !found {
		c.String(http.StatusNotFound, "Metric not found")
		log.Error("metric not found")
		return
	}
	c.JSON(http.StatusOK, value)
}

func (h *MetricHandler) GetAllMetrics(c *gin.Context) {
	metrics, err := h.Storage.GetAllMetrics()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting metrics")
		log.Error(err)
		return
	}

	c.HTML(http.StatusOK, "metrics.tmpl", gin.H{
		"metrics": metrics,
	})
}
