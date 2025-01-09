package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

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
