package coreserver

import (
	"fmt"
	"net/http"

	"github.com/echo9et/alerting/internal/compgzip"
	"github.com/echo9et/alerting/internal/entities"
	"github.com/echo9et/alerting/internal/logger"
	"github.com/echo9et/alerting/internal/server/handlers"
	"github.com/go-chi/chi/v5"
)

func middleware(h http.HandlerFunc) http.HandlerFunc {
	return logger.RequestLogger(compgzip.GzipMiddleware(h))
}

func GetRouter(addrDatabase string, storage entities.Storage) *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", middleware(func(w http.ResponseWriter, r *http.Request) {
		metricsHandle(w, r, storage)
	}))
	router.Post("/update/", middleware(func(w http.ResponseWriter, r *http.Request) {
		WriteMetricJSONHandle(w, r, storage)
	}))
	router.Post("/updates/", middleware(func(w http.ResponseWriter, r *http.Request) {
		WriteMetricsJSONHandle(w, r, storage)
	}))
	router.Post("/update/{type}/{name}/{value}", middleware(func(w http.ResponseWriter, r *http.Request) {
		setMetricHandle(w, r, storage)
	}))
	router.Post("/value/", middleware(func(w http.ResponseWriter, r *http.Request) {
		ReadMetricJSONHandle(w, r, storage)
	}))
	router.Get("/value/{type}/{name}", middleware(func(w http.ResponseWriter, r *http.Request) {
		metricHandle(w, r, storage)
	}))
	router.Get("/ping", middleware(func(w http.ResponseWriter, r *http.Request) {
		PingDatabase(w, r, addrDatabase, storage)
	}))
	return router
}

func Run(addr, addrDatabase string, storage entities.Storage) error {
	return http.ListenAndServe(addr, GetRouter(addrDatabase, storage))
}

func metricHandle(w http.ResponseWriter, r *http.Request, s entities.Storage) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	t := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	var value string
	status := false

	switch t {
	case handlers.Gauge:
		value, status = s.GetGauge(name)
	case handlers.Counter:
		value, status = s.GetCounter(name)
	}

	if status {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintln(value)))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func metricsHandle(w http.ResponseWriter, r *http.Request, s entities.Storage) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Accept") == "text/html" {
		w.Header().Set("Content-Type", "text/html")
	}

	metrics := s.AllMetrics()
	for k, v := range metrics {
		w.Write([]byte(fmt.Sprintln(k, v)))
	}
	w.WriteHeader(http.StatusOK)
}

func setMetricHandle(w http.ResponseWriter, r *http.Request, s entities.Storage) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := handlers.WriteMetric(w, r, s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func WriteMetricJSONHandle(w http.ResponseWriter, r *http.Request, s entities.Storage) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := handlers.WriteMetricJSON(w, r, s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func WriteMetricsJSONHandle(w http.ResponseWriter, r *http.Request, s entities.Storage) {
	if r.Method != http.MethodPost {
		fmt.Println("=== Error: WriteMetricsJSONHandle", 405)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-type") != "application/json" {
		fmt.Println("=== Error: WriteMetricsJSONHandle", 400)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := handlers.WriteMetricsJSON(w, r, s); err != nil {
		fmt.Println(" === Error: WriteMetricsJSONHandle", 505)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func ReadMetricJSONHandle(w http.ResponseWriter, r *http.Request, s entities.Storage) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := handlers.ReadMetricJSON(w, r, s); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func PingDatabase(w http.ResponseWriter, r *http.Request, addr string, s entities.Storage) {
	if !s.Ping() {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-type", "text/plain")
		out := []byte(" ")
		w.Write(out)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
