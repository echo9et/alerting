package handlers

import (
	"strconv"
	"testing"

	"github.com/echo9et/alerting/internal/server/storage"
)

func TestHandlerCounters(t *testing.T) {
	s := storage.NewMemStorage()
	tests := []struct {
		name  string
		value string
		want  error
	}{
		{
			name:  "test float",
			value: "3.14",
			want:  strconv.ErrSyntax,
		},
		{
			name:  "test ok",
			value: "3",
			want:  nil,
		},
		{
			name:  "test string",
			value: "x",
			want:  strconv.ErrSyntax,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handlerCounters(s, tt.name, tt.value)
			if (got != nil || tt.want != nil) && got == tt.want {
				t.Errorf("handlerGauge(%s, %s) = %s, want: %s", tt.name, tt.value, got, tt.want)
			}
		})
	}
}

func TestHandlerGauge(t *testing.T) {
	s := storage.NewMemStorage()
	tests := []struct {
		name  string
		value string
		want  error
	}{
		{
			name:  "test ok float",
			value: "3.14",
			want:  nil,
		},
		{
			name:  "test ok int",
			value: "3",
			want:  nil,
		},
		{
			name:  "test string",
			value: "x",
			want:  strconv.ErrSyntax,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handlerGauge(s, tt.name, tt.value)
			if (got != nil || tt.want != nil) && got == tt.want {
				t.Errorf("handlerGauge(%s, %s) = %s, want: %s", tt.name, tt.value, got, tt.want)
			}
		})
	}
}
