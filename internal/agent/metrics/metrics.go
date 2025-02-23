package metrics

import (
	"fmt"
	"math/rand/v2"
	"runtime"

	"github.com/echo9et/alerting/internal/entities"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type Metricer interface {
	ToJSON() []entities.MetricsJSON
	Update()
}

type MetricsRuntime struct {
	Memory      runtime.MemStats
	Storage     map[string]interface{}
	PollCount   int64
	RandomValue float64
}

func NewMetrics() *MetricsRuntime {
	runtime.GC()
	metrics := MetricsRuntime{}
	metrics.Memory.Alloc = 30
	metrics.Storage = make(map[string]interface{})
	metrics.Storage["Alloc"] = &metrics.Memory.Alloc
	metrics.Storage["BuckHashSys"] = &metrics.Memory.BuckHashSys
	metrics.Storage["Frees"] = &metrics.Memory.Frees
	metrics.Storage["GCCPUFraction"] = &metrics.Memory.GCCPUFraction
	metrics.Storage["GCSys"] = &metrics.Memory.GCSys
	metrics.Storage["HeapAlloc"] = &metrics.Memory.HeapAlloc
	metrics.Storage["HeapIdle"] = &metrics.Memory.HeapIdle
	metrics.Storage["HeapInuse"] = &metrics.Memory.HeapInuse
	metrics.Storage["HeapObjects"] = &metrics.Memory.HeapObjects
	metrics.Storage["HeapReleased"] = &metrics.Memory.HeapReleased
	metrics.Storage["HeapSys"] = &metrics.Memory.HeapSys
	metrics.Storage["LastGC"] = &metrics.Memory.LastGC
	metrics.Storage["Lookups"] = &metrics.Memory.Lookups
	metrics.Storage["MCacheInuse"] = &metrics.Memory.MCacheInuse
	metrics.Storage["MCacheSys"] = &metrics.Memory.MCacheSys
	metrics.Storage["MSpanInuse"] = &metrics.Memory.MSpanInuse
	metrics.Storage["MSpanSys"] = &metrics.Memory.MSpanSys
	metrics.Storage["Mallocs"] = &metrics.Memory.Mallocs
	metrics.Storage["NextGC"] = &metrics.Memory.NextGC
	metrics.Storage["NumForcedGC"] = &metrics.Memory.NumForcedGC
	metrics.Storage["NumGC"] = &metrics.Memory.NumGC
	metrics.Storage["OtherSys"] = &metrics.Memory.OtherSys
	metrics.Storage["PauseTotalNs"] = &metrics.Memory.PauseTotalNs
	metrics.Storage["StackInuse"] = &metrics.Memory.StackInuse
	metrics.Storage["StackSys"] = &metrics.Memory.StackSys
	metrics.Storage["Sys"] = &metrics.Memory.Sys
	metrics.Storage["TotalAlloc"] = &metrics.Memory.TotalAlloc

	return &metrics
}

func (m *MetricsRuntime) Update() {
	runtime.ReadMemStats(&m.Memory)
	m.PollCount += 1
	m.RandomValue = rand.Float64()
}

func (m *MetricsRuntime) ToJSON() []entities.MetricsJSON {
	metrics := make([]entities.MetricsJSON, 0)
	for key, value := range m.Storage {
		metric := entities.MetricsJSON{}
		var fValue float64
		switch v := value.(type) {
		case int:
			fValue = float64(v)
		case uint64:
			fValue = float64(v)
		case *uint64:
			fValue = float64(*v)
		case *uint32:
			fValue = float64(*v)
		case *float64:
			fValue = *v
		}
		metric.Value = &fValue
		metric.MType = entities.Gauge
		metric.ID = key

		metrics = append(metrics, metric)
	}
	metrics = append(metrics, entities.MetricsJSON{
		ID:    "PollCount",
		MType: entities.Counter,
		Delta: &m.PollCount,
	})

	metrics = append(metrics, entities.MetricsJSON{
		ID:    "RandomValue",
		MType: entities.Gauge,
		Value: &m.RandomValue,
	})

	return metrics
}

type MetricsMem struct {
	Storage map[string]interface{}
}

func NewMetricsMem() *MetricsMem {
	metrics := MetricsMem{}
	metrics.Storage = make(map[string]interface{})
	return &metrics
}

func (m *MetricsMem) Update() {
	v, _ := mem.VirtualMemory()
	m.Storage["TotalMemory"] = v.Total
	m.Storage["FreeMemory"] = v.Free
	c, _ := cpu.Percent(0, true)
	for i, percent := range c {
		m.Storage[fmt.Sprintf("CPUutilization%d", i+1)] = percent
	}
}

func (m *MetricsMem) ToJSON() []entities.MetricsJSON {
	metrics := make([]entities.MetricsJSON, 0)
	for key, value := range m.Storage {
		metric := entities.MetricsJSON{}
		var fValue float64
		switch v := value.(type) {
		case float64:
			fValue = float64(v)
		case uint64:
			fValue = float64(v)
		}
		metric.Value = &fValue
		metric.MType = entities.Gauge
		metric.ID = key

		metrics = append(metrics, metric)
	}
	return metrics
}
