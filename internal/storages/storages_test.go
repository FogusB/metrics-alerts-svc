package storages

import (
	"testing"
)

func TestNewMemStorage(t *testing.T) {
	storage := NewMemStorage()
	if storage == nil {
		t.Errorf("NewMemStorage() не должен возвращать nil")
	}
	//if storage.metrics != nil {
	//	t.Errorf("Новый экземпляр MemStorage должен быть пустым")
	//}
}

func TestMemStorage_UpdateMetric(t *testing.T) {
	storage := NewMemStorage()

	// Позитивный тест для Gauge
	storage.UpdateMetric("test_gauge", Gauge, MetricValue{GaugeValue: 10.5})
	if val, ok := storage.metrics["test_gauge"]; !ok || val.GaugeValue != 10.5 {
		t.Errorf("Ожидалось значение Gauge 10.5, получено %v", val.GaugeValue)
	}

	// Позитивный тест для Counter
	storage.UpdateMetric("test_counter", Counter, MetricValue{CounterValue: 5})
	storage.UpdateMetric("test_counter", Counter, MetricValue{CounterValue: 3})
	if val, ok := storage.metrics["test_counter"]; !ok || val.CounterValue != 8 {
		t.Errorf("Ожидалось значение Counter 8, получено %v", val.CounterValue)
	}

	//// Негативный тест: обновление несуществующего типа метрики
	//defer func() {
	//	if r := recover(); r == nil {
	//		t.Errorf("Ожидалась паника при обновлении метрики неизвестного типа")
	//	}
	//}()
	//storage.UpdateMetric("test_unknown", "unknown", MetricValue{})
}
