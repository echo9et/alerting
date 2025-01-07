package coreserver

import (
	"fmt"
	"net/http"

	"github.com/echo9et/alerting/cmd/server/handlers"
	"github.com/go-chi/chi/v5"
)

func GetRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", metricsHandle)
	router.Post("/update/{type}/{name}/{value}", setMetricHandle)
	router.Get("/value/{type}/{name}", metricHandle)
	return router
}

func Run(addr string) error {
	return http.ListenAndServe(addr, GetRouter())
}

func metricHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
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
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintln(value)))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func metricsHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
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

	// if r.Header.Get("Content-Type") != "text/plain" {
	// 	w.WriteHeader(http.StatusUnsupportedMediaType)
	// 	return
	// }

	if err := handlers.WriteMetric(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
