package client

import (
	"bytes"
	"compress/gzip"
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
		mj.MType = entities.Counter
		mj.Delta = &v
	}
	data, err := json.Marshal(mj)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	err = a.SendDataToServerNoGzip(data)
	if err != nil {
		fmt.Println("ERROR:", err)
	}
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

func (a *Agent) SendDataToServer(data []byte) error {
	cd, err := CompressGzip(data)
	if err != nil {
		return err
	}

	body := bytes.NewReader(cd)

	url := fmt.Sprintf("http://%s/update", a.outServer)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Enecoding", "gzip")
	req.Header.Set("Content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (a *Agent) SendDataToServerNoGzip(data []byte) error {
	body := bytes.NewReader(data)

	url := fmt.Sprintf("http://%s/update", a.outServer)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func CompressGzip(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}
	err = gz.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}
	return b.Bytes(), nil
}
