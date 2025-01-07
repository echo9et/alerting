package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"time"
)

const (
	reportInterval = 1
	pullInterval   = 2
)

type metrics struct {
	memory         runtime.MemStats
	supportMetrics map[string]interface{}
	pullCount      int64
	randomValue    float64
}

func (m *metrics) pullMetrics() {
	for key, value := range m.supportMetrics {
		var sendValue float64
		switch v := value.(type) {
		case *uint64:
			sendValue = float64(*v)
		case *uint32:
			sendValue = float64(*v)
		case *float64:
			sendValue = *v
		}
		sendMetric(key, sendValue)
	}
	sendMetric("PollCount", m.pullCount)
	sendMetric("RandomValue", m.randomValue)
}

func sendMetric(name string, value interface{}) {
	var url string
	switch v := value.(type) {
	case float64:
		url = fmt.Sprintf("http://%s/update/gauge/%s/%v", *addrServer, name, v)
	case int64:
		url = fmt.Sprintf("http://%s/update/counter/%s/%v", *addrServer, name, v)
	}

	r := bytes.NewReader([]byte(``))
	resp, err := http.Post(url, "text/plain", r)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	defer resp.Body.Close()
}
func updateMetrics(m *metrics, reportInterval int, pullInterval int) {
	runtime.GC()
	counter := 0
	for {
		runtime.ReadMemStats(&m.memory)
		m.pullCount += 1
		m.randomValue = rand.Float64()

		time.Sleep(time.Second * time.Duration(reportInterval))

		counter += reportInterval
		if counter >= pullInterval {
			counter = 0
			m.pullMetrics()
		}
	}
}
func main() {
	initAgent()
	metrics := NewMetrics()
	updateMetrics(&metrics, *reportTimeout, *pollTimeout)
}

func initAgent() {
	parseFlags()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		*addrServer = envRunAddr
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		*addrServer = envReportInterval
	}

	if envPoolInterval := os.Getenv("POLL_INTERVAL"); envPoolInterval != "" {
		*addrServer = envPoolInterval
	}

}

func NewMetrics() metrics {
	metrics := metrics{}
	metrics.supportMetrics = make(map[string]interface{})
	metrics.supportMetrics["Alloc"] = &metrics.memory.Alloc
	metrics.supportMetrics["BuckHashSys"] = &metrics.memory.BuckHashSys
	metrics.supportMetrics["Frees"] = &metrics.memory.Frees
	metrics.supportMetrics["GCCPUFraction"] = &metrics.memory.GCCPUFraction
	metrics.supportMetrics["GCSys"] = &metrics.memory.GCSys
	metrics.supportMetrics["HeapAlloc"] = &metrics.memory.HeapAlloc
	metrics.supportMetrics["HeapIdle"] = &metrics.memory.HeapIdle
	metrics.supportMetrics["HeapInuse"] = &metrics.memory.HeapInuse
	metrics.supportMetrics["HeapObjects"] = &metrics.memory.HeapObjects
	metrics.supportMetrics["HeapReleased"] = &metrics.memory.HeapReleased
	metrics.supportMetrics["HeapSys"] = &metrics.memory.HeapSys
	metrics.supportMetrics["LastGC"] = &metrics.memory.LastGC
	metrics.supportMetrics["Lookups"] = &metrics.memory.Lookups
	metrics.supportMetrics["MCacheInuse"] = &metrics.memory.MCacheInuse
	metrics.supportMetrics["MCacheSys"] = &metrics.memory.MCacheSys
	metrics.supportMetrics["MSpanInuse"] = &metrics.memory.MSpanInuse
	metrics.supportMetrics["MSpanSys"] = &metrics.memory.MSpanSys
	metrics.supportMetrics["Mallocs"] = &metrics.memory.Mallocs
	metrics.supportMetrics["NextGC"] = &metrics.memory.NextGC
	metrics.supportMetrics["NumForcedGC"] = &metrics.memory.NumForcedGC
	metrics.supportMetrics["NumGC"] = &metrics.memory.NumGC
	metrics.supportMetrics["OtherSys"] = &metrics.memory.OtherSys
	metrics.supportMetrics["PauseTotalNs"] = &metrics.memory.PauseTotalNs
	metrics.supportMetrics["StackInuse"] = &metrics.memory.StackInuse
	metrics.supportMetrics["StackSys"] = &metrics.memory.StackSys
	metrics.supportMetrics["Sys"] = &metrics.memory.Sys
	metrics.supportMetrics["TotalAlloc"] = &metrics.memory.TotalAlloc

	return metrics
}
