package coreserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusHandler(t *testing.T) {
	type want struct {
		code        int
		url         string
		method      string
		response    string
		contentType string
	}
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
			name: "positive method test #1",
			want: want{
				url:         "/update/counter/test/1",
				method:      http.MethodPost,
				code:        415,
				response:    "",
				contentType: "media",
			},
		},
		{
			name: "negative method test #1",
			want: want{
				url:         "/update/counter/test/1",
				method:      http.MethodGet,
				code:        405,
				response:    "",
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.want.method, test.want.url, nil)
			request.Header.Add("Content-Type", test.want.contentType)
			w := httptest.NewRecorder()
			setMetricHandle(w, request)

			res := w.Result()
			println(res.Status)
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
		})
	}
}
