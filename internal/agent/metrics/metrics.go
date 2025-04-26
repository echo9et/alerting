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

type data struct {
	Counters map[string]uint64
	Gauges   map[string]float64
}

func (d *data) toJSON() []entities.MetricsJSON {
	metrics := make([]entities.MetricsJSON, 0)
	for key, value := range d.Gauges {
		metric := entities.MetricsJSON{}
		metric.ID = key
		metric.MType = entities.Gauge
		metric.Value = &value
		metrics = append(metrics, metric)
	}
	for key, value := range d.Counters {
		iValue := int64(value)
		metric := entities.MetricsJSON{}
		metric.ID = key
		metric.MType = entities.Counter
		metric.Delta = &iValue
		metrics = append(metrics, metric)
	}

	return metrics
}

func newData() data {
	return data{
		Counters: make(map[string]uint64),
		Gauges:   make(map[string]float64),
	}
}

type MetricsRuntime struct {
	Memory runtime.MemStats
	data   data
}

// NewMetrics возвращает структуру с метриками рантайма приложения
func NewMetrics() *MetricsRuntime {
	runtime.GC()
	return &MetricsRuntime{
		data: newData(),
	}
}

func (m *MetricsRuntime) Update() {
	runtime.ReadMemStats(&m.Memory)
	m.data.Counters["PollCount"] += 1
	m.data.Counters["Alloc"] = m.Memory.Alloc
	m.data.Counters["BuckHashSys"] = m.Memory.BuckHashSys
	m.data.Counters["Frees"] = m.Memory.Frees
	m.data.Counters["GCSys"] = m.Memory.GCSys
	m.data.Counters["HeapAlloc"] = m.Memory.HeapAlloc
	m.data.Counters["HeapIdle"] = m.Memory.HeapIdle
	m.data.Counters["HeapInuse"] = m.Memory.HeapInuse
	m.data.Counters["HeapObjects"] = m.Memory.HeapObjects
	m.data.Counters["HeapReleased"] = m.Memory.HeapReleased
	m.data.Counters["HeapSys"] = m.Memory.HeapSys
	m.data.Counters["LastGC"] = m.Memory.LastGC
	m.data.Counters["Lookups"] = m.Memory.Lookups
	m.data.Counters["MCacheInuse"] = m.Memory.MCacheInuse
	m.data.Counters["MCacheSys"] = m.Memory.MCacheSys
	m.data.Counters["MSpanInuse"] = m.Memory.MSpanInuse
	m.data.Counters["MSpanSys"] = m.Memory.MSpanSys
	m.data.Counters["Mallocs"] = m.Memory.Mallocs
	m.data.Counters["NextGC"] = m.Memory.NextGC
	m.data.Counters["NumForcedGC"] = uint64(m.Memory.NumForcedGC)
	m.data.Counters["NumGC"] = uint64(m.Memory.NumGC)
	m.data.Counters["OtherSys"] = m.Memory.OtherSys
	m.data.Counters["PauseTotalNs"] = m.Memory.PauseTotalNs
	m.data.Counters["StackInuse"] = m.Memory.StackInuse
	m.data.Counters["StackSys"] = m.Memory.StackSys
	m.data.Counters["Sys"] = m.Memory.Sys
	m.data.Counters["TotalAlloc"] = m.Memory.TotalAlloc

	m.data.Gauges["RandomValue"] = rand.Float64()
	m.data.Gauges["GCCPUFraction"] = m.Memory.GCCPUFraction
}

func (m *MetricsRuntime) ToJSON() []entities.MetricsJSON {
	return m.data.toJSON()
}

type MetricsMem struct {
	data data
}

// NewMetricsMem возвращает структуру для сбора информации о памяти.
func NewMetricsMem() *MetricsMem {
	return &MetricsMem{
		data: newData(),
	}
}

func (m *MetricsMem) Update() {
	v, _ := mem.VirtualMemory()
	m.data.Counters["TotalMemory"] = v.Total
	m.data.Counters["FreeMemory"] = v.Free

	c, _ := cpu.Percent(0, true)
	for i, percent := range c {
		m.data.Gauges[fmt.Sprintf("CPUutilization%d", i+1)] = percent
	}
}

func (m *MetricsMem) ToJSON() []entities.MetricsJSON {
	return m.data.toJSON()
}
