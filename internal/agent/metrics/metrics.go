package metrics

import "runtime"

type Metrics struct {
	Memory         runtime.MemStats
	SupportMetrics map[string]interface{}
	PollCount      int64
	RandomValue    float64
}

func NewMetrics() Metrics {
	metrics := Metrics{}
	metrics.SupportMetrics = make(map[string]interface{})
	metrics.SupportMetrics["Alloc"] = &metrics.Memory.Alloc
	metrics.SupportMetrics["BuckHashSys"] = &metrics.Memory.BuckHashSys
	metrics.SupportMetrics["Frees"] = &metrics.Memory.Frees
	metrics.SupportMetrics["GCCPUFraction"] = &metrics.Memory.GCCPUFraction
	metrics.SupportMetrics["GCSys"] = &metrics.Memory.GCSys
	metrics.SupportMetrics["HeapAlloc"] = &metrics.Memory.HeapAlloc
	metrics.SupportMetrics["HeapIdle"] = &metrics.Memory.HeapIdle
	metrics.SupportMetrics["HeapInuse"] = &metrics.Memory.HeapInuse
	metrics.SupportMetrics["HeapObjects"] = &metrics.Memory.HeapObjects
	metrics.SupportMetrics["HeapReleased"] = &metrics.Memory.HeapReleased
	metrics.SupportMetrics["HeapSys"] = &metrics.Memory.HeapSys
	metrics.SupportMetrics["LastGC"] = &metrics.Memory.LastGC
	metrics.SupportMetrics["Lookups"] = &metrics.Memory.Lookups
	metrics.SupportMetrics["MCacheInuse"] = &metrics.Memory.MCacheInuse
	metrics.SupportMetrics["MCacheSys"] = &metrics.Memory.MCacheSys
	metrics.SupportMetrics["MSpanInuse"] = &metrics.Memory.MSpanInuse
	metrics.SupportMetrics["MSpanSys"] = &metrics.Memory.MSpanSys
	metrics.SupportMetrics["Mallocs"] = &metrics.Memory.Mallocs
	metrics.SupportMetrics["NextGC"] = &metrics.Memory.NextGC
	metrics.SupportMetrics["NumForcedGC"] = &metrics.Memory.NumForcedGC
	metrics.SupportMetrics["NumGC"] = &metrics.Memory.NumGC
	metrics.SupportMetrics["OtherSys"] = &metrics.Memory.OtherSys
	metrics.SupportMetrics["PauseTotalNs"] = &metrics.Memory.PauseTotalNs
	metrics.SupportMetrics["StackInuse"] = &metrics.Memory.StackInuse
	metrics.SupportMetrics["StackSys"] = &metrics.Memory.StackSys
	metrics.SupportMetrics["Sys"] = &metrics.Memory.Sys
	metrics.SupportMetrics["TotalAlloc"] = &metrics.Memory.TotalAlloc

	return metrics
}
