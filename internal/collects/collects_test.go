package collects

import (
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestCollectMetrics(t *testing.T) {
	metrics := CollectMetrics()

	// Проверяем, что возвращаемая map не пуста.
	assert.NotEmpty(t, metrics, "Метрики не должны быть пустыми")

	expectedKeys := []string{
		"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys",
		"HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased",
		"HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys",
		"MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC",
		"NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys",
		"Sys", "TotalAlloc", "RandomValue",
	}

	for _, key := range expectedKeys {
		assert.Contains(t, metrics, key, "Результат должен содержать ключ")
	}

	//for _, key := range expectedKeys {
	//	if _, exists := metrics[key]; !exists {
	//		t.Errorf("Expected key %s not found in metrics map", key)
	//	}
	//}

	// Проверяем типы некоторых ключей
	if _, ok := metrics["RandomValue"].(float64); !ok {
		t.Errorf("RandomValue is not of type float64")
	}

	if _, ok := metrics["Alloc"].(uint64); !ok {
		t.Errorf("Alloc is not of type uint64")
	}

	// Проверяем, что значение RandomValue находится в пределах от 0 до 1
	if val, ok := metrics["RandomValue"].(float64); ok {
		if val < 0 || val > 1 {
			t.Errorf("RandomValue is out of range: got %v", val)
		}
	} else {
		t.Errorf("RandomValue is not of type float64")
	}

	// Проверка на изменение значений при повторном вызове для динамически изменяющихся метрик.
	firstAlloc := metrics["Alloc"]
	runtime.GC() // Принудительно запускаем сборку мусора для изменения статистики памяти.
	metricsAfterGC := CollectMetrics()
	if firstAlloc == metricsAfterGC["Alloc"] {
		t.Errorf("Expected Alloc metric to change after GC, but it didn't")
	}
}
