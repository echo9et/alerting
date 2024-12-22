package main

import (
	"net/http"
	"strconv"
	"strings"
)

func handlerCounters(name, value string) error {
	iValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	storage.counters[name] = iValue
	return nil
}

func handlerGauge(name, value string) error {
	iValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	storage.gauge[name] = iValue
	return nil
}

var (
	supportMetrices = map[string]func(string, string) error{
		Gauge:   handlerGauge,
		Counter: handlerCounters,
	}
	storage = MemStorage{make(map[string]float64), make(map[string]int64)}
)

func main() {
	if error := run(); error != nil {
		panic("error")
	}
}

const (
	Gauge   = "gauge"
	Counter = "counter"
)

type MemStorage struct {
	gauge    map[string]float64
	counters map[string]int64
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, webhook)
	return http.ListenAndServe(`:8080`, mux)
}

func webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	r.URL.Path = r.URL.Path[len("/update/"):]
	handlerMetrics(w, r)
}

func handlerMetrics(w http.ResponseWriter, r *http.Request) error {
	param := strings.Split(r.URL.Path, "/")

	if len(param) < 2 || len(param) > 3 {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}
	typeMetric := param[0]
	handlerMetric, ok := supportMetrices[typeMetric]
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
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
