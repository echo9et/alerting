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

func (a *Agent) compressMetrics() error {
	metricsJSON := make([]entities.MetricsJSON, 0)
	for key, value := range a.metrics.SupportMetrics {
		metric := entities.MetricsJSON{}
		var fValue float64
		switch v := value.(type) {
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

		metricsJSON = append(metricsJSON, metric)
	}
	metricsJSON = append(metricsJSON, entities.MetricsJSON{
		ID:    "PollCount",
		MType: entities.Counter,
		Delta: &a.metrics.PollCount,
	})

	metricsJSON = append(metricsJSON, entities.MetricsJSON{
		ID:    "RandomValue",
		MType: entities.Gauge,
		Value: &a.metrics.RandomValue,
	})

	data, err := json.Marshal(metricsJSON)
	if err != nil {
		fmt.Println("ERROR:", err)
		return err
	}
	cd, err := CompressGzip(data)
	if err != nil {
		return err
	}
	return a.SendGzipToServer(cd)
}

func (a *Agent) SendGzipToServer(data []byte) error {
	body := bytes.NewReader(data)

	url := fmt.Sprintf("http://%s/updates", a.outServer)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
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
	// err = a.SendDataToServer(data, false)
	err = a.SendDataToServerGzip(data)
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
			a.compressMetrics()
			// a.pollMetrics()
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

func (a *Agent) SendDataToServerGzip(data []byte) error {
	cd, err := CompressGzip(data)
	if err != nil {
		return err
	}
	return a.SendDataToServer(cd, true)
}

func (a *Agent) SendDataToServer(data []byte, isGzip bool) error {
	body := bytes.NewReader(data)

	url := fmt.Sprintf("http://%s/update", a.outServer)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	if isGzip {
		req.Header.Set("Content-Encoding", "gzip")
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
