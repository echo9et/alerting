package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type MetricsJson struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type Storage interface {
	GetCounter(string) (string, bool)
	SetCounter(string, int64)
	GetGauge(string) (string, bool)
	SetGauge(string, float64)
	AllMetrics() map[string]string
}

const (
	Gauge   = "gauge"
	Counter = "counter"
)

var supportMetrics = map[string]func(Storage, string, string) error{
	Gauge:   handlerGauge,
	Counter: handlerCounters,
}

func handlerCounters(s Storage, name, sValue string) error {
	iValue, err := strconv.ParseInt(sValue, 10, 64)
	if err != nil {
		return err
	}
	s.SetCounter(name, iValue)
	return nil
}

func handlerGauge(s Storage, name, sValue string) error {
	fValue, err := strconv.ParseFloat(sValue, 64)
	if err != nil {
		return err
	}
	s.SetGauge(name, fValue)
	return nil
}

func WriteMetric(w http.ResponseWriter, r *http.Request, s Storage) error {
	handlerMetric, ok := supportMetrics[chi.URLParam(r, "type")]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}
	name, value := chi.URLParam(r, "name"), chi.URLParam(r, "value")
	err := handlerMetric(s, name, value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func WriteJsonMetric(w http.ResponseWriter, r *http.Request, s Storage) error {
	var metricsJson MetricsJson
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(buf.Bytes(), &metricsJson); err != nil {
		return err
	}

	if err = saveJsonMetrics(s, metricsJson); err != nil {
		return err
	}

	out, err := json.Marshal(metricsJson)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
	return nil
}

func ReadJsonMetric(w http.ResponseWriter, r *http.Request, s Storage) error {
	var metricsJson MetricsJson
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(buf.Bytes(), &metricsJson); err != nil {
		return err
	}
	if err = getJsonMetrics(s, &metricsJson); err != nil {
		return err
	}

	out, err := json.Marshal(metricsJson)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)

	return nil
}

func getJsonMetrics(s Storage, mj *MetricsJson) error {
	switch mj.MType {
	case Counter:
		value, status := s.GetCounter(mj.ID)
		if !status {
			return errors.New("counter not found")
		}
		iValue, _ := strconv.ParseInt(value, 10, 64)
		mj.Delta = &iValue
	case Gauge:
		value, status := s.GetGauge(mj.ID)
		if !status {
			return errors.New("gauge not found")
		}
		dValue, _ := strconv.ParseFloat(value, 64)
		mj.Value = &dValue
	}

	return nil
}

func saveJsonMetrics(s Storage, mj MetricsJson) error {
	switch mj.MType {
	case Counter:
		s.SetCounter(mj.ID, *mj.Delta)
	case Gauge:
		s.SetGauge(mj.ID, *mj.Value)
	default:
		return errors.New("Unknow type")
	}

	return nil
}
