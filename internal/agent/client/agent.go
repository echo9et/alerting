package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/echo9et/alerting/internal/agent/metrics"
	"github.com/echo9et/alerting/internal/entities"
)

type Agent struct {
	metrics   *metrics.Metrics
	outServer string
}

func NewAgent(addressServer string) *Agent {
	return &Agent{metrics: metrics.NewMetrics(),
		outServer: addressServer,
	}
}

func (a *Agent) SendMetric(name string, value interface{}) {

	var mj entities.MetricsJSON
	mj.ID = name
	switch v := value.(type) {
	case float64:
		mj.MType = entities.Gauge
		mj.Value = &v
	case int64:
		print("======")
		mj.MType = entities.Counter
		mj.Delta = &v
		print("-----")
	}
	data, err := json.Marshal(mj)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	r := bytes.NewReader(data)
	url := fmt.Sprintf("http://%s/update", a.outServer)
	resp, err := http.Post(url, "application/json", r)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	defer resp.Body.Close()
}
func (a Agent) UpdateMetrics(reportInterval time.Duration, pollInterval time.Duration) {
	runtime.GC()
	counter := time.Duration(0)
	for {
		runtime.ReadMemStats(&a.metrics.Memory)
		a.metrics.PollCount += 1
		a.metrics.RandomValue = rand.Float64()
		time.Sleep(reportInterval)

		counter += reportInterval
		if counter >= pollInterval {
			counter = time.Duration(0)
			a.pollMetrics()
		}
	}
}

func (a *Agent) pollMetrics() {
	for key, value := range a.metrics.SupportMetrics {

		var sendValue float64
		switch v := value.(type) {
		case *uint64:
			sendValue = float64(*v)
		case *uint32:
			sendValue = float64(*v)
		case *float64:
			sendValue = *v
		}
		a.SendMetric(key, sendValue)
	}
	a.SendMetric("PollCount", a.metrics.PollCount)
	a.SendMetric("RandomValue", a.metrics.RandomValue)
}
