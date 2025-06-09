package coreserver

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log/slog"

	"github.com/echo9et/alerting/internal/entities"
	pb "github.com/echo9et/alerting/proto"
)

type ServerGrpc struct {
	MetricsSever
	CryptoKey *rsa.PrivateKey
	Storage   entities.ManagerValues
}

// addMetric добавление метрики в хранилище
func (s *ServerGrpc) addMetric(m *pb.Metric) {
	switch m.Type {
	case pb.Metric_GAUGE:
		s.Storage.SetGauge(m.Id, m.Value)
	case pb.Metric_GOUNTER:
		s.Storage.SetCounter(m.Id, m.Delta)
	default:
		slog.Warn("Unkonow type metric")
	}
}

// Decryptor данных.
func (s *ServerGrpc) decryptData(encryptedData []byte) ([]byte, error) {
	data := []byte(encryptedData)

	decrypted, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, s.CryptoKey, data, nil)
	if err != nil {
		slog.Error(fmt.Sprintf("Ошибка при дешифрование информации %s", err))
		return decrypted, err
	}
	return decrypted, nil
}

type MetricsSever struct {
	pb.UnimplementedMetricsServer
}

// UpdateMetric реализует интерфейс добавления одной метрики.
func (s *ServerGrpc) UpdateMetric(ctx context.Context, in *pb.UpdateMetricRequest) (*pb.UpdateMetricResponse, error) {
	var response pb.UpdateMetricResponse
	s.addMetric(in.Metric)
	return &response, nil
}

// UpdateMetrics реализует интерфейс добавления метрик.
func (s *ServerGrpc) UpdateMetrics(ctx context.Context, in *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error) {
	var response pb.UpdateMetricsResponse

	for _, metric := range in.Metrics {
		s.addMetric(metric)
	}
	return &response, nil
}

// UpdateEncrypteMetrics реализует интерфейс добавления списка метрик.
func (s *ServerGrpc) UpdateEncrypteMetrics(ctx context.Context, in *pb.UpdateEncrypteMetricsRequest) (*pb.UpdateEncrypteMetricsResponse, error) {
	var response pb.UpdateEncrypteMetricsResponse
	decrypted, err := s.decryptData(in.Data)

	if err != nil {
		slog.Error(fmt.Sprintf("UpdateEncrypteMetricsr: %s", err))
		return &response, err
	}

	var metrics []*pb.Metric
	inReader := bytes.NewReader(decrypted)
	dec := gob.NewDecoder(inReader)

	if err := dec.Decode(&metrics); err != nil {
		slog.Error(fmt.Sprintf("UpdateEncrypteMetricsr: %s", err))
		return &response, err
	}

	for _, metric := range metrics {
		s.addMetric(metric)
	}

	return &response, nil
}
