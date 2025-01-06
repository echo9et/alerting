package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/echo9et/alerting/cmd/server/storage"
	"github.com/go-chi/chi/v5"
)

var storageInstance = storage.NewMemStorage()

const (
	Gauge   = "gauge"
	Counter = "counter"
)

var supportMetrics = map[string]func(string, string) error{
	Gauge:   handlerGauge,
	Counter: handlerCounters,
}

func handlerCounters(name, value string) error {
	iValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	storageInstance.Counters[name] += iValue // Увеличиваем значение счетчика
	return nil
}

func handlerGauge(name, value string) error {
	fValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	storageInstance.Gauges[name] = fValue // Устанавливаем значение гаужа
	return nil
}

func GetCounterValue(name string) (string, bool) {
	value, ok := storageInstance.Counters[name]
	if ok {
		return fmt.Sprint(value), true
	}
	return "", false
}

func GetGaugeValue(name string) (string, bool) {
	value, ok := storageInstance.Gauges[name]
	if ok {
		return fmt.Sprint(value), true
	}
	return "", false
}

func GetMetrics() map[string]string {
	return storageInstance.AllMetrics()
}

func WriteMetric(w http.ResponseWriter, r *http.Request) error {
	handlerMetric, ok := supportMetrics[chi.URLParam(r, "type")]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}
	name, value := chi.URLParam(r, "name"), chi.URLParam(r, "value")
	err := handlerMetric(name, value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}

	w.WriteHeader(http.StatusOK)
	print(http.StatusOK, name, value, "\n")
	return nil
}
