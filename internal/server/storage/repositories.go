package storage

import (
	"fmt"

	"github.com/echo9et/alerting/internal/entities"
)

type Store struct {
	Metrics map[string]entities.MetricsJSON
}

func NewStore() *Store {
	return &Store{
		Metrics: make(map[string]entities.MetricsJSON),
	}
}

func (s *Store) GetCounter(name string) (string, bool) {
	metric, ok := s.Metrics[name]
	if ok {
		if metric.MType == entities.Counter {
			return fmt.Sprint(*metric.Delta), true
		}
	}
	return "", false
}

func (s *Store) SetCounter(name string, iValue int64) {
	if metric, ok := s.Metrics[name]; ok {
		newValue := *(metric.Delta) + iValue
		metric.Delta = &newValue
		s.Metrics[name] = metric
	} else {
		s.Metrics[name] = entities.MetricsJSON{
			ID:    name,
			MType: entities.Counter,
			Delta: &iValue,
		}
	}
}

func (s *Store) GetGauge(name string) (string, bool) {
	metric, ok := s.Metrics[name]
	if ok {
		if metric.MType == entities.Gauge {
			return fmt.Sprint(*metric.Value), true
		}
	}
	return "", false
}

func (s *Store) SetGauge(name string, fValue float64) {
	if metric, ok := s.Metrics[name]; ok {
		metric.Value = &fValue
		s.Metrics[name] = metric
	} else {
		s.Metrics[name] = entities.MetricsJSON{
			ID:    name,
			MType: entities.Gauge,
			Value: &fValue,
		}
	}
}

func (s *Store) AllMetrics() map[string]string {
	out := make(map[string]string)
	for k, v := range s.Metrics {
		out[k] = fmt.Sprint(v)
	}
	return out
}

func (s *Store) AllMetricsJSON() []entities.MetricsJSON {

	metricsJSON := make([]entities.MetricsJSON, 0)

	for _, metric := range s.Metrics {
		metricsJSON = append(metricsJSON, metric)
	}
	return metricsJSON
}
