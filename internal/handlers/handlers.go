package handlers

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"

	"github.com/FogusB/metrics-alerts-svc/internal/storages"
	"github.com/gin-gonic/gin"
)

type MetricHandler struct {
	Storage Storage
}

type Storage interface {
	UpdateMetric(name string, mType storages.MetricType, value storages.MetricValue) error
	GetMetric(name string) (storages.MetricValue, bool)
	GetAllMetrics() (map[string]storages.MetricValue, error)
}

func (h *MetricHandler) UpdateMetricValue(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.String(http.StatusMethodNotAllowed, "Method Not Allowed")
		log.Warning("method not allowed")
		return
	}

	contentType := c.GetHeader("Content-Type")
	if contentType != "" && contentType != "text/plain" {
		c.String(http.StatusUnsupportedMediaType, "Unsupported Media Type")
		log.Warning("unsupported media type")
		return
	}

	mType := storages.MetricType(c.Param("type"))
	name := c.Param("name")
	rawValue := c.Param("value")

	if mType != storages.Gauge && mType != storages.Counter {
		c.String(http.StatusBadRequest, "Invalid metric type")
		log.Warning("invalid metric type")
		return
	}

	var value storages.MetricValue
	var err error
	if mType == storages.Gauge {
		value.GaugeValue, err = strconv.ParseFloat(rawValue, 64)
	} else {
		value.CounterValue, err = strconv.ParseUint(rawValue, 10, 64)
	}
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid value for %s", mType)
		log.Error(err)
		return
	}

	err = h.Storage.UpdateMetric(name, mType, value)
	if err != nil {
		log.Error(err)
		c.String(http.StatusInternalServerError, "Error updating metric")
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

	c.String(http.StatusOK, fmt.Sprintf("%v", value))
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
