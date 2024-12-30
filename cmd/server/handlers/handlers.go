package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/echo9et/alerting/cmd/server/storage"
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
	fmt.Println(name, value)
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

func HandlerMetrics(w http.ResponseWriter, r *http.Request) error {
	param := strings.Split(r.URL.Path, "/")

	if len(param) < 2 || len(param) > 3 {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}
	typeMetric := param[0]
	handlerMetric, ok := supportMetrics[typeMetric]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}

	if len(param) == 2 {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	name, value := param[1], param[2]
	err := handlerMetric(name, value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
