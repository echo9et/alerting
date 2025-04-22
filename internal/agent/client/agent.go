package client

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/echo9et/alerting/internal/agent/metrics"
	"github.com/echo9et/alerting/internal/entities"
	"github.com/echo9et/alerting/internal/hashing"
)

type Agent struct {
	metrics   *metrics.MetricsRuntime
	outServer string
}

// NewAgent конструктор для создания объекта агента
func NewAgent(addressServer string) *Agent {
	return &Agent{metrics: metrics.NewMetrics(),
		outServer: addressServer,
	}
}

// UpdateMetrics запуск сбора метрик и отправки их на сервер.
func (a Agent) UpdateMetrics(reportInterval time.Duration, pollInterval time.Duration, key string, rateLimit int64) {
	queueMetrics := make(chan []entities.MetricsJSON)
	defer close(queueMetrics)
	go generatorMetric(queueMetrics, metrics.NewMetrics(), pollInterval, reportInterval)
	go generatorMetric(queueMetrics, metrics.NewMetricsMem(), pollInterval, reportInterval)

	var wg sync.WaitGroup
	for range rateLimit - 1 {
		wg.Add(1)
		go a.poll(queueMetrics, key, &wg)
	}
	wg.Wait()
}

// poll отправка метрик на сервер.
func (a *Agent) poll(in chan []entities.MetricsJSON, secretKey string, wg *sync.WaitGroup) {
	defer wg.Done()
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
		entities.Retry(func() error {
			return a.SendToServer(cd, secretKey)
		})
	}
}

// SendToServer отправка метрик на сервер.
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

// CompressGzip сжатие метрик перед отправкой на сервер.
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

// generatorMetric сбор метрик с заданным интервалом времени.
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
