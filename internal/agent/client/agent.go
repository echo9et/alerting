package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/echo9et/alerting/internal/agent/metrics"
	"github.com/echo9et/alerting/internal/entities"
	"github.com/echo9et/alerting/internal/hashing"

	pb "github.com/echo9et/alerting/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	_ "google.golang.org/grpc/encoding/gzip"
)

type Agent struct {
	metrics   *metrics.MetricsRuntime
	outServer string
	selfIP    string
	useGRPC   bool
}

// NewAgent конструктор для создания объекта агента
func NewAgent(addressServer, selfIP string, useGRPC bool) *Agent {
	return &Agent{metrics: metrics.NewMetricsRuntime(),
		outServer: addressServer,
		selfIP:    selfIP,
		useGRPC:   useGRPC,
	}
}

// UpdateMetrics запуск сбора метрик и отправки их на сервер.
func (a Agent) UpdateMetrics(reportInterval time.Duration, pollInterval time.Duration, key string, rateLimit int64, pubKey *rsa.PublicKey) {

	queueMetrics := make(chan []entities.MetricsJSON)
	defer close(queueMetrics)
	go generatorMetric(queueMetrics, metrics.NewMetricsRuntime(), pollInterval, reportInterval)
	go generatorMetric(queueMetrics, metrics.NewMetricsMem(), pollInterval, reportInterval)

	var wg sync.WaitGroup
	for range rateLimit - 1 {
		wg.Add(1)
		go a.push(queueMetrics, key, &wg, pubKey)
	}
	wg.Wait()
}

// poll отправка метрик на сервер.
func (a *Agent) push(in chan []entities.MetricsJSON, secretKey string, wg *sync.WaitGroup, pubKey *rsa.PublicKey) {
	defer wg.Done()

	if a.useGRPC {
		if err := a.pushGRPC(in, pubKey); err != nil {
			slog.Error(fmt.Sprintln(err))
		}
		return
	}

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

		if pubKey != nil {
			cd, err = rsa.EncryptOAEP(
				sha256.New(),
				rand.Reader,
				pubKey,
				cd,
				nil,
			)
			if err != nil {
				slog.Error(fmt.Sprintln(err))
				return
			}
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
	req.Header.Set("X-Real-IP", a.selfIP)
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

func (a *Agent) pushGRPC(in chan []entities.MetricsJSON, pubKey *rsa.PublicKey) error {
	conn, err := grpc.NewClient(":3200",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewMetricsClient(conn)

	for jsonMetrics := range in {
		var metrics []*pb.Metric
		for _, metric := range jsonMetrics {
			if metric.MType == entities.Counter {
				metrics = append(metrics, &pb.Metric{
					Id:    metric.ID,
					Type:  pb.Metric_GOUNTER,
					Delta: *metric.Delta,
				})
			} else if metric.MType == entities.Gauge {
				metrics = append(metrics, &pb.Metric{Id: metric.ID,
					Type:  pb.Metric_GAUGE,
					Value: *metric.Value,
				})
			} else {
				slog.Error("Unknow type")
			}
		}

		if pubKey != nil {
			var data bytes.Buffer
			enc := gob.NewEncoder(&data)
			err := enc.Encode(metrics)
			if err != nil {
				slog.Error(fmt.Sprintln("Enecode :", err))
				continue
			}
			enecrypted, err := rsa.EncryptOAEP(
				sha256.New(),
				rand.Reader,
				pubKey,
				data.Bytes(),
				nil,
			)
			if err != nil {
				slog.Error(fmt.Sprintln("Enecode :", err))
				continue
			}
			resp, err := c.UpdateEncrypteMetrics(context.Background(), &pb.UpdateEncrypteMetricsRequest{
				Data: enecrypted,
			})
			if err != nil {
				slog.Error(fmt.Sprintln(err))
				continue
			}
			if resp.Error != "" {
				slog.Error(fmt.Sprintln("resp.Error", err))
				continue
			}
		} else {
			resp, err := c.UpdateMetrics(context.Background(), &pb.UpdateMetricsRequest{
				Metrics: metrics,
			})
			if err != nil {
				slog.Error(fmt.Sprintln(err))
				continue
			}
			if resp.Error != "" {
				slog.Error(fmt.Sprintln("resp.Error", resp.Error))
				continue
			}
		}
	}
	return nil
}
