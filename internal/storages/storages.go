package storages

import (
	"errors"
	"sync"
)

// MetricType определяет тип метрики: gauge или counter
type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

// MetricValue хранит значение метрики, которое может быть представлено как float64 или uint64
type MetricValue struct {
	GaugeValue   float64
	CounterValue uint64
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
func (s *MemStorage) UpdateMetric(name string, mType MetricType, value MetricValue) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch mType {
	case Gauge:
		s.metrics[name] = value
	case Counter:
		if existing, ok := s.metrics[name]; ok {
			value.CounterValue += existing.CounterValue
		}
		s.metrics[name] = value
	default:
		return errors.New("unknown metric type")
	}

	//log.Infof("Key: %s, value: %v\n", name, value)
	return nil
}

// GetMetric возвращает метрику из хранилища
func (s *MemStorage) GetMetric(name string) (MetricValue, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.metrics[name]
	return value, ok
}

// GetAllMetrics возвращает список всех метрик и значений из хранилища
func (s *MemStorage) GetAllMetrics() (map[string]MetricValue, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Создаем новую мапу для возврата значений
	metricsCopy := make(map[string]MetricValue, len(s.metrics))
	for key, value := range s.metrics {
		metricsCopy[key] = value
	}

	return metricsCopy, nil
}
