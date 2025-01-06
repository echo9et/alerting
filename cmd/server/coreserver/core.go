package coreserver

import (
	"fmt"
	"net/http"

	"github.com/echo9et/alerting/cmd/server/handlers"
	"github.com/go-chi/chi/v5"
)

func Run() error {
	router := chi.NewRouter()
	router.Post("/", metricsHandle)
	router.Post("/update/{type}/{name}/{value}", setMetricHandle)
	router.Post("/value/{type}/{name}", metricHandle)
	return http.ListenAndServe(":8080", router)
}

func metricHandle(w http.ResponseWriter, r *http.Request) {
	typ := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	var value string
	status := false

	switch typ {
	case handlers.Gauge:
		value, status = handlers.GetGaugeValue(name)
	case handlers.Counter:
		value, status = handlers.GetCounterValue(name)
	}

	if status {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(fmt.Sprintln(name, value)))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func metricsHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	metrics := handlers.GetMetrics()
	for k, v := range metrics {
		w.Write([]byte(fmt.Sprintln(k, v)))
	}
	w.WriteHeader(http.StatusOK)
}

func setMetricHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	if err := handlers.WriteMetric(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
