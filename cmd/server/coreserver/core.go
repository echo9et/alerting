package coreserver

import (
	"net/http"

	"github.com/echo9et/alerting/cmd/server/handlers"
)

func Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", Webhook)
	return http.ListenAndServe(":8080", mux)
}

func Webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	r.URL.Path = r.URL.Path[len("/update/"):]

	if err := handlers.HandlerMetrics(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
