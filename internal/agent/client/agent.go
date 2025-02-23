package client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/echo9et/alerting/internal/agent/metrics"
	"github.com/echo9et/alerting/internal/entities"
	"github.com/echo9et/alerting/internal/hashing"
)

type Agent struct {
	metrics   *metrics.MetricsRuntime
	outServer string
}

func NewAgent(addressServer string) *Agent {
	return &Agent{metrics: metrics.NewMetrics(),
		outServer: addressServer,
	}
}

func (a Agent) UpdateMetrics(reportInterval time.Duration, pollInterval time.Duration, key string, rateLimit int64) {
	in := make(chan []entities.MetricsJSON)
	defer close(in)
	go generatorMetric(in, metrics.NewMetrics(), pollInterval, reportInterval)
	go generatorMetric(in, metrics.NewMetricsMem(), pollInterval, reportInterval)
	a.poll(in, key)
}

func (a *Agent) poll(in chan []entities.MetricsJSON, secretKey string) {
	for metric := range in {
		data, err := json.Marshal(metric)
		if err != nil {
			slog.Error(fmt.Sprintln(err))
			return
		}
		cd, err := CompressGzip(data)
		if err != nil {
			slog.Error(fmt.Sprintln(err))
			return
		}
		a.SendToServer(cd, secretKey)
	}
}

func (a *Agent) SendToServer(data []byte, secretKey string) error {
	slog.Info("SendToServer")
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

func generatorMetric(in chan []entities.MetricsJSON, m metrics.Metricer, reportInterval time.Duration, pollInterval time.Duration) {
	counter := time.Duration(0)
	for {
		m.Update()
		time.Sleep(reportInterval)
		counter += reportInterval
		if counter >= pollInterval {
			counter = time.Duration(0)
			in <- m.ToJSON()
		}
	}
}
