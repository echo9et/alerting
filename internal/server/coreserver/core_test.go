package coreserver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/echo9et/alerting/internal/server/storage"
	"github.com/stretchr/testify/assert"
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
	s := storage.NewMemStorage()
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
		ts := httptest.NewServer(GetRouter(s))
		defer ts.Close()
		t.Run(test.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, test.want)
			assert.Equal(t, test.want.code, resp.StatusCode)
			resp.Body.Close()
			// print(get)
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
