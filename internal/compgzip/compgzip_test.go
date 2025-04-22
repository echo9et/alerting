package compgzip

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Gzip!"))
}

func TestGzipMiddleware_Decompression(t *testing.T) {
	handler := GzipMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(body)
	}))
	server := httptest.NewServer(handler)
	defer server.Close()

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write([]byte("Compressed Request"))
	assert.NoError(t, err)
	assert.NoError(t, gz.Close())

	req, err := http.NewRequest("POST", server.URL, &buf)
	assert.NoError(t, err)

	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "Compressed Request", string(body))
}

// проверка ошибки при неверных данных в запросе
func TestGzipMiddleware_ErrorHandling(t *testing.T) {
	handler := GzipMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.ReadAll(r.Body)
	}))
	server := httptest.NewServer(handler)
	defer server.Close()

	req, err := http.NewRequest("POST", server.URL, strings.NewReader("invalid gzip data"))
	assert.NoError(t, err)

	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

// проверка, что middleware не изменяет ответ, если заголовок Content-Type не является текстовым
func TestGzipMiddleware_ContentTypes(t *testing.T) {
	handler := GzipMiddleware(http.HandlerFunc(testHandler))
	server := httptest.NewServer(handler)
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	assert.NoError(t, err)

	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "image/png")
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, "", resp.Header.Get("Content-Encoding"))

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, Gzip!", string(body))
}
