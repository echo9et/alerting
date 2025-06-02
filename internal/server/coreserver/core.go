package coreserver

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/echo9et/alerting/internal/compgzip"
	"github.com/echo9et/alerting/internal/entities"
	"github.com/echo9et/alerting/internal/hashing"
	"github.com/echo9et/alerting/internal/logger"
	"github.com/echo9et/alerting/internal/server/handlers"
	pb "github.com/echo9et/alerting/proto"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
)

// applyGzipMiddleware применяет GzipMiddleware к обработчику.
func applyDecryt(h http.HandlerFunc, privateKey *rsa.PrivateKey) http.HandlerFunc {
	if privateKey != nil {
		return DecryptMiddleware(h, privateKey)
	}
	return h
}

func applyGzipMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return compgzip.GzipMiddleware(h)
}

// applyHashMiddleware применяет HashMiddleware к обработчику, если secretKey указан.
func applyHashMiddleware(h http.HandlerFunc, secretKey string) http.HandlerFunc {
	if secretKey != "" {
		return HashMiddleware(h, secretKey)
	}
	return h
}

// applyRequestLogger применяет RequestLogger к обработчику.
func applyRequestLogger(h http.HandlerFunc) http.HandlerFunc {
	return logger.RequestLogger(h)
}

func applyTrustSubnet(h http.HandlerFunc, subnet *net.IPNet) http.HandlerFunc {
	if subnet != nil {
		return TrustedSubnetMiddleware(h, subnet)
	}
	return h
}

// Добавляет к обработчику протоколирование и сжатие в формате gzip.
// Если указан секретный ключ, оно также добавляет промежуточное программное обеспечение для хэширования.
func middleware(h http.HandlerFunc, secretKey string, privateKey *rsa.PrivateKey, trustedSubnet *net.IPNet) http.HandlerFunc {
	h = applyTrustSubnet(h, trustedSubnet)
	h = applyDecryt(h, privateKey)
	h = applyRequestLogger(h)
	h = applyHashMiddleware(h, secretKey)
	h = applyGzipMiddleware(h)

	return h
}

// Добавляет к обработчику протоколирование и сжатие в формате gzip.
func HashMiddleware(h http.HandlerFunc, secretKey string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ow := hashing.NewHashingWriter(w, secretKey)
		hash := r.Header.Get("HashSHA256")
		if hash != "" {
			var buf bytes.Buffer
			tee := io.TeeReader(r.Body, &buf)

			body, err := io.ReadAll(tee)
			if err != nil {
				slog.Error("не удалсть считать тело запроса")
				w.WriteHeader(http.StatusBadRequest)
			}

			r.Body = io.NopCloser(&buf)

			if hash != hashing.GetHash(body, secretKey) {
				w.WriteHeader(http.StatusBadRequest)
			}
		}
		h.ServeHTTP(ow, r)
	})
}

// Возвращает маршрутизатор сервера.
func GetRouter(addrDatabase string, storage entities.Storage, secretKey string, privateKey *rsa.PrivateKey, trustedSubnet *net.IPNet) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/", middleware(func(w http.ResponseWriter, r *http.Request) {
		metricsHandle(w, r, storage)
	}, secretKey, privateKey, trustedSubnet))

	router.Post("/update/", middleware(func(w http.ResponseWriter, r *http.Request) {
		WriteMetricJSONHandle(w, r, storage)
	}, secretKey, privateKey, trustedSubnet))

	router.Post("/updates/", middleware(func(w http.ResponseWriter, r *http.Request) {
		WriteMetricsJSONHandle(w, r, storage)
	}, secretKey, privateKey, trustedSubnet))

	router.Post("/update/{type}/{n, privateKey)ame}/{value}", middleware(func(w http.ResponseWriter, r *http.Request) {
		setMetricHandle(w, r, storage)
	}, secretKey, privateKey, trustedSubnet))

	router.Post("/value/", middleware(func(w http.ResponseWriter, r *http.Request) {
		ReadMetricJSONHandle(w, r, storage)
	}, secretKey, privateKey, trustedSubnet))

	router.Get("/value/{type}/{name}", middleware(func(w http.ResponseWriter, r *http.Request) {
		metricHandle(w, r, storage)
	}, secretKey, privateKey, trustedSubnet))

	router.Get("/ping", middleware(func(w http.ResponseWriter, r *http.Request) {
		PingDatabase(w, r, addrDatabase, storage)
	}, secretKey, privateKey, trustedSubnet))
	return router
}

// Запуск сервера.
func Run(addr, addrDatabase string, storage entities.Storage, secretKey string, privateKey *rsa.PrivateKey, trustedSubnet *net.IPNet) error {
	var server = http.Server{Addr: addr, Handler: GetRouter(addrDatabase, storage, secretKey, privateKey, trustedSubnet)}
	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	go func() {
		<-sigint
		if err := server.Shutdown(context.Background()); err != nil {
			slog.Error(fmt.Sprintf("HTTP server Shutdown: %v", err))
		}
		close(idleConnsClosed)
	}()

	go func() {
		listen, err := net.Listen("tcp", ":3200")
		if err != nil {
			slog.Error(fmt.Sprintf("listent grps: %s", err))
			return
		}
		s := grpc.NewServer()
		serverGrpc := ServerGrpc{
			CryptoKey: privateKey,
			Storage:   storage,
		}
		pb.RegisterMetricsServer(s, &serverGrpc)

		slog.Info("Сервер gRPC начал работу")
		if err := s.Serve(listen); err != nil {
			slog.Error(fmt.Sprintf("listent grps: %s", err))
		}
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		slog.Error(fmt.Sprintf("server ListenAndServe:%v", err))
		return err
	}
	slog.Info("Server Shutdown")
	return nil
}

// Возвращает значения метрик по типу и имени.
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

// Возвращает все метрики.
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

// Обновляет метрику.
func setMetricHandle(w http.ResponseWriter, r *http.Request, s entities.Storage) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := handlers.WriteMetric(w, r, s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Записывает метрику в хранилище.
func WriteMetricJSONHandle(w http.ResponseWriter, r *http.Request, s entities.Storage) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := handlers.WriteMetricJSON(w, r, s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Записывает метрики в хранилище.
func WriteMetricsJSONHandle(w http.ResponseWriter, r *http.Request, s entities.Storage) {
	if r.Method != http.MethodPost {
		slog.Error(fmt.Sprintln("=== Error: WriteMetricsJSONHandle", 405))
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		slog.Error(fmt.Sprintln("=== Error: WriteMetricsJSONHandle", 400))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := handlers.WriteMetricsJSON(w, r, s); err != nil {
		slog.Error(fmt.Sprintln(" === Error: WriteMetricsJSONHandle", 505))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

// Считывает метрику из хранилища.
func ReadMetricJSONHandle(w http.ResponseWriter, r *http.Request, s entities.Storage) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := handlers.ReadMetricJSON(w, r, s); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

// Проверяет доступность хранилища.
func PingDatabase(w http.ResponseWriter, r *http.Request, addr string, s entities.Storage) {
	if !s.Ping() {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain")
		out := []byte(" ")
		w.Write(out)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// DecryptMiddleware декодирует полученные данные
func DecryptMiddleware(h http.HandlerFunc, privateKey *rsa.PrivateKey) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error(fmt.Sprintf("Ошибка при чтение информации из запроса %s", err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		decrypted, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, data, nil)
		if err != nil {
			slog.Error(fmt.Sprintf("Ошибка при дешифрование информации %s", err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(decrypted))
		h.ServeHTTP(w, r)
	})
}

// TrustedSubnetMiddleware проверяет, входит ли запрос в доверенную сеть
func TrustedSubnetMiddleware(h http.HandlerFunc, subnet *net.IPNet) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get("X-Real-IP")
		if ip == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		parsedIP := net.ParseIP(ip)
		if !subnet.Contains(parsedIP) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	})
}
