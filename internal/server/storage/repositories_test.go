package storage

import (
	"testing"
	"time"
)

func TestStorageCounter(t *testing.T) {
	storage := NewMemStore()
	tests := []struct {
		name    string
		storage *MemStore
		value   int64
		want    string
	}{
		{
			name:    "test MemStore.SetCounter #1",
			storage: storage,
			value:   5,
			want:    "5",
		},
		{
			name:    "test MemStore.SetCounter #2",
			storage: storage,
			value:   50,
			want:    "55",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storage.SetCounter("test", tt.value)
			v, ok := tt.storage.GetCounter("test")
			if !ok || tt.want != v {
				t.Errorf("SetCounter(\"test\", %v) = %s, want: %s", tt.value, v, tt.want)
			}
		})
	}
}

func TestStorageGauge(t *testing.T) {
	tests := []struct {
		name    string
		storage *MemStore
		value   float64
		want    string
	}{
		{
			name:    "test MemStore.SetGauge #1",
			storage: NewMemStore(),
			value:   5.,
			want:    "5",
		},
		{
			name:    "test MemStore.SetGauge #1",
			storage: NewMemStore(),
			value:   50.,
			want:    "50",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storage.SetGauge("test", tt.value)
			v, ok := tt.storage.GetGauge("test")
			if !ok || tt.want != v {
				t.Errorf("SetGauge(\"test\", %v) = %s, want: %s", tt.value, v, tt.want)
			}
		})
	}
}

func TestSaver(t *testing.T) {
	saver, _ := NewSaver(NewMemStore(), "test.json", false, time.Second*2000)
	tests := []struct {
		name  string
		saver *Saver
		value int64
		want  string
	}{
		{
			name:  "test Saver.SetCounter #1",
			saver: saver,
			value: 5,
			want:  "5",
		},
		{
			name:  "test Saver.counter #2",
			saver: saver,
			value: 50,
			want:  "55",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.saver.SetCounter("test", tt.value)
			v, ok := tt.saver.GetCounter("test")
			if !ok || tt.want != v {
				t.Errorf("SetCounter(\"test\", %v) = %s, want: %s", tt.value, v, tt.want)
			}
		})
	}
}

func TestSavereGauge(t *testing.T) {
	saver, _ := NewSaver(NewMemStore(), "test.json", false, time.Second*2000)
	tests := []struct {
		name  string
		saver *Saver
		value float64
		want  string
	}{
		{
			name:  "test Saver.setGauge #1",
			saver: saver,
			value: 5.,
			want:  "5",
		},
		{
			name:  "test Saver.setGauge #2",
			saver: saver,
			value: 50.,
			want:  "50",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.saver.SetGauge("test", tt.value)
			v, ok := tt.saver.GetGauge("test")
			if !ok || tt.want != v {
				t.Errorf("SetGauge(\"test\", %v) = %s, want: %s", tt.value, v, tt.want)
			}
		})
	}
}
