package metrics

import (
	"strings"
	"testing"

	"github.com/echo9et/alerting/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestMetricsRuntime_Update(t *testing.T) {
	metrics := NewMetrics()

	// Вызываем метод Update
	metrics.Update()

	// Проверяем, что данные были обновлены
	assert.NotZero(t, metrics.data.Counters["PollCount"], "PollCount should be updated")
	assert.NotZero(t, metrics.data.Gauges["RandomValue"], "RandomValue should be updated")
	assert.NotZero(t, metrics.data.Gauges["GCCPUFraction"], "GCCPUFraction should be updated")
}

func TestMetricsRuntime_ToJSON(t *testing.T) {
	metrics := NewMetrics()
	metrics.Update()

	// Получаем JSON-представление метрик
	jsonMetrics := metrics.ToJSON()

	// Проверяем, что JSON содержит ожидаемые метрики
	foundPollCount := false
	foundRandomValue := false
	foundGCCPUFraction := false

	for _, metric := range jsonMetrics {
		if metric.ID == "PollCount" && metric.MType == entities.Counter {
			foundPollCount = true
			assert.NotNil(t, metric.Delta, "Delta for PollCount should not be nil")
		}
		if metric.ID == "RandomValue" && metric.MType == entities.Gauge {
			foundRandomValue = true
			assert.NotNil(t, metric.Value, "Value for RandomValue should not be nil")
		}
		if metric.ID == "GCCPUFraction" && metric.MType == entities.Gauge {
			foundGCCPUFraction = true
			assert.NotNil(t, metric.Value, "Value for GCCPUFraction should not be nil")
		}
	}

	assert.True(t, foundPollCount, "PollCount metric should be present in JSON")
	assert.True(t, foundRandomValue, "RandomValue metric should be present in JSON")
	assert.True(t, foundGCCPUFraction, "GCCPUFraction metric should be present in JSON")
}

func TestMetricsMem_Update(t *testing.T) {
	metrics := NewMetricsMem()

	// Вызываем метод Update
	metrics.Update()

	// Проверяем, что данные были обновлены
	assert.NotZero(t, metrics.data.Counters["TotalMemory"], "TotalMemory should be updated")
	assert.NotZero(t, metrics.data.Counters["FreeMemory"], "FreeMemory should be updated")

	// Проверяем, что CPU utilization метрики присутствуют
	cpuMetricFound := false
	for key := range metrics.data.Gauges {
		if strings.HasPrefix(key, "CPUutilization") {
			cpuMetricFound = true
			break
		}
	}
	assert.True(t, cpuMetricFound, "At least one CPUutilization metric should be present")
}

func TestMetricsMem_ToJSON(t *testing.T) {
	metrics := NewMetricsMem()
	metrics.Update()

	// Получаем JSON-представление метрик
	jsonMetrics := metrics.ToJSON()

	// Проверяем, что JSON содержит ожидаемые метрики
	foundTotalMemory := false
	foundFreeMemory := false
	foundCPUUtilization := false

	for _, metric := range jsonMetrics {
		if metric.ID == "TotalMemory" && metric.MType == entities.Counter {
			foundTotalMemory = true
			assert.NotNil(t, metric.Delta, "Delta for TotalMemory should not be nil")
		}
		if metric.ID == "FreeMemory" && metric.MType == entities.Counter {
			foundFreeMemory = true
			assert.NotNil(t, metric.Delta, "Delta for FreeMemory should not be nil")
		}
		if strings.HasPrefix(metric.ID, "CPUutilization") && metric.MType == entities.Gauge {
			foundCPUUtilization = true
			assert.NotNil(t, metric.Value, "Value for CPUutilization should not be nil")
		}
	}

	assert.True(t, foundTotalMemory, "TotalMemory metric should be present in JSON")
	assert.True(t, foundFreeMemory, "FreeMemory metric should be present in JSON")
	assert.True(t, foundCPUUtilization, "At least one CPUutilization metric should be present in JSON")
}
