package internal

import (
	"fmt"
	"sync"
)

// MetricType определяет тип метрики: gauge или counter
type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

// MetricValue хранит значение метрики, которое может быть представлено как float64 или int64
type MetricValue struct {
	GaugeValue   float64
	CounterValue int64
}

// MemStorage структура для хранения метрик
type MemStorage struct {
	mu      sync.RWMutex
	metrics map[string]MetricValue
}

// NewMemStorage создает новый экземпляр MemStorage
func NewMemStorage() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]MetricValue),
	}
}

// UpdateMetric обновляет метрику в хранилище
func (s *MemStorage) UpdateMetric(name string, mType MetricType, value MetricValue) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if mType == Gauge {
		s.metrics[name] = value
	} else if mType == Counter {
		if existing, ok := s.metrics[name]; ok {
			value.CounterValue += existing.CounterValue
		}
		s.metrics[name] = value
	}
	fmt.Printf("Key: %s, value: %v\n", name, value)
}
