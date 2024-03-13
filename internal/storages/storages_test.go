package storages

import (
	"reflect"
	"testing"
)

func TestNewMemStorage(t *testing.T) {
	storage := NewMemStorage()
	if storage == nil {
		t.Errorf("NewMemStorage() не должен возвращать nil")
	}
	//if len(storage.metrics) != 0 {
	//	t.Errorf("Новый экземпляр MemStorage должен быть пустым")
	//}
}

func TestMemStorage_UpdateMetric(t *testing.T) {
	storage := NewMemStorage()
	tests := []struct {
		name    string
		mType   MetricType
		value   MetricValue
		wantErr bool
	}{
		{"testGauge", Gauge, MetricValue{GaugeValue: 10.5}, false},
		{"testCounter", Counter, MetricValue{CounterValue: 5}, false},
		{"testUnknown", MetricType("unknown"), MetricValue{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := storage.UpdateMetric(tt.name, tt.mType, tt.value); (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.UpdateMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemStorage_GetMetric(t *testing.T) {
	storage := NewMemStorage()
	// Установка тестовых значений
	storage.metrics["testGauge"] = MetricValue{GaugeValue: 10.5}
	storage.metrics["testCounter"] = MetricValue{CounterValue: 5}

	tests := []struct {
		name      string
		wantValue MetricValue
		wantOk    bool
	}{
		{"testGauge", MetricValue{GaugeValue: 10.5}, true},
		{"testCounter", MetricValue{CounterValue: 5}, true},
		{"unknown", MetricValue{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotOk := storage.GetMetric(tt.name)
			if !reflect.DeepEqual(gotValue, tt.wantValue) || gotOk != tt.wantOk {
				t.Errorf("MemStorage.GetMetric() = %v, %v, want %v, %v", gotValue, gotOk, tt.wantValue, tt.wantOk)
			}
		})
	}
}

func TestMemStorage_GetAllMetrics(t *testing.T) {
	storage := NewMemStorage()
	// Установка тестовых значений
	storage.metrics["testGauge"] = MetricValue{GaugeValue: 10.5}
	storage.metrics["testCounter"] = MetricValue{CounterValue: 5}

	want := map[string]MetricValue{
		"testGauge":   {GaugeValue: 10.5},
		"testCounter": {CounterValue: 5},
	}

	got, err := storage.GetAllMetrics()
	if err != nil {
		t.Errorf("MemStorage.GetAllMetrics() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MemStorage.GetAllMetrics() = %v, want %v", got, want)
	}
}
