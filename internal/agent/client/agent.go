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
	"github.com/echo9et/alerting/internal/hashing"
	"github.com/shirou/gopsutil/v4/mem"
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

func (a Agent) UpdateMetrics(reportInterval time.Duration, pollInterval time.Duration, key string, rateLimit int64) {
	runtime.GC()
	counter := time.Duration(0)
	for {
		runtime.ReadMemStats(&a.metrics.Memory)
		a.metrics.PollCount += 1
		a.metrics.RandomValue = rand.Float64()
		time.Sleep(reportInterval)
		v, _ := mem.VirtualMemory()
		a.metrics.SupportMetrics["TotalMemory"] = v.Total
		a.metrics.SupportMetrics["FreeMemory"] = v.Free
		a.metrics.SupportMetrics["CPUutilization1"] = runtime.NumCPU()
		counter += reportInterval
		if counter >= pollInterval {
			counter = time.Duration(0)
			a.pollMetrics(key)
		}
	}
}

func (a *Agent) pollMetrics(secretKey string) error {
	data, err := json.Marshal(a.dataJSON())
	if err != nil {
		fmt.Println("ERROR:", err)
		return err
	}
	cd, err := CompressGzip(data)
	if err != nil {
		return err
	}

	return a.SendToServer(cd, secretKey)
}

func (a *Agent) dataJSON() []entities.MetricsJSON {
	metrics := make([]entities.MetricsJSON, 0)
	for key, value := range a.metrics.SupportMetrics {
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
		Delta: &a.metrics.PollCount,
	})

	metrics = append(metrics, entities.MetricsJSON{
		ID:    "RandomValue",
		MType: entities.Gauge,
		Value: &a.metrics.RandomValue,
	})

	return metrics
}

func (a *Agent) SendToServer(data []byte, secretKey string) error {

	body := bytes.NewReader(data)
	url := fmt.Sprintf("http://%s/updates/", a.outServer)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	if secretKey != "" {
		req.Header.Set("HashSHA256", hashing.GetHash(data, secretKey))
	}
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
