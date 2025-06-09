package coreserver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/echo9et/alerting/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type want struct {
	code        int
	url         string
	method      string
	response    string
	contentType string
}

func TestStatusHandler(t *testing.T) {
	s := storage.NewMemStore()
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive method test #1",
			want: want{
				url:         "/update/counter/test/1",
				method:      http.MethodPost,
				code:        200,
				response:    "",
				contentType: "text/plain",
			},
		},
		{
			name: "negative method test #2",
			want: want{
				url:         "/update/counter/test/1",
				method:      http.MethodPost,
				code:        200,
				response:    "",
				contentType: "media",
			},
		},
		{
			name: "negative method test #3",
			want: want{
				url:         "/update/counter/test/1",
				method:      http.MethodGet,
				code:        405,
				response:    "",
				contentType: "text/plain",
			},
		},

		{
			name: "positive method test #4",
			want: want{
				url:         "/",
				method:      http.MethodGet,
				code:        200,
				response:    "",
				contentType: "text/plain",
			},
		},
		{
			name: "negative method test #5",
			want: want{
				url:         "/",
				method:      http.MethodPost,
				code:        405,
				response:    "",
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		ts := httptest.NewServer(GetRouter("", s, "", nil, nil))
		defer ts.Close()
		t.Run(test.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, test.want)
			assert.Equal(t, test.want.code, resp.StatusCode)
			resp.Body.Close()
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, testData want) (*http.Response, string) {
	req, err := http.NewRequest(testData.method, ts.URL+testData.url, nil)
	req.Header.Add("Content-Type", testData.contentType)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

type MockHandlerFunc struct {
	mock.Mock
}

func (m *MockHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func TestApplyRequestLogger(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test postive case #1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHandler := new(MockHandlerFunc)

			handler := applyRequestLogger(mockHandler.ServeHTTP)
			assert.NotNil(t, handler)
		})
	}
}

func TestApplyHashMiddleware_NoSecretKey(t *testing.T) {
	textResponse := "test response"
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(textResponse))
	}
	httpHandler := http.HandlerFunc(handler)

	middlewareHandler := applyHashMiddleware(httpHandler, "")

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	middlewareHandler.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		assert.Error(t, err)
	}
	answer := string(respBody)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, textResponse, answer)
	assert.Empty(t, resp.Header.Get("Hashsha256"))
}

func TestApplyGzipMiddleware(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test response"))
	}
	httpHandler := http.HandlerFunc(handler)
	middlewareHandler := applyGzipMiddleware(httpHandler)

	// запрос c сжатием
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	middlewareHandler.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, rec.Header().Get("Content-Encoding"), "gzip")

	// запрос без сжатия
	reqNoGzip := httptest.NewRequest("GET", "/", nil)
	recNoGzip := httptest.NewRecorder()
	middlewareHandler.ServeHTTP(recNoGzip, reqNoGzip)

	rec = httptest.NewRecorder()
	resp = rec.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Empty(t, rec.Header().Get("Content-Encoding"))
}

func TestApplyHashMiddleware(t *testing.T) {
	textResponse := "test response"
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(textResponse))
	}
	httpHandler := http.HandlerFunc(handler)

	secretKey := "my-secret-key"
	middlewareHandler := applyHashMiddleware(httpHandler, secretKey)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	middlewareHandler.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, textResponse, rec.Body.String())
	assert.Equal(t, resp.Header.Get("Hashsha256"), "ebf31a7d817d2091f7238be75431e05dd831ceaa349253b7eb2cd6c71ecbae65")
}
