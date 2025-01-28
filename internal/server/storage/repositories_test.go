package storage

import (
	"fmt"
	"testing"
	"time"
)

func TestStorageCounter(t *testing.T) {
	tests := []struct {
		name    string
		storage *Store
		value   int64
	}{
		{
			name:    "test set counter",
			storage: NewStore(),
			value:   5, // errors.New("float"),
		},
		{
			name:    "test set counter",
			storage: NewStore(),
			value:   50, // errors.New("float"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storage.SetCounter("test", tt.value)
			fmt.Println(tt.storage.GetCounter("test"))
		})
	}
}

func TestStorageGauge(t *testing.T) {
	tests := []struct {
		name    string
		storage *Store
		value   float64
	}{
		{
			name:    "test set counter",
			storage: NewStore(),
			value:   5., // errors.New("float"),
		},
		{
			name:    "test set counter",
			storage: NewStore(),
			value:   50., // errors.New("float"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.storage.SetGauge("test", tt.value)
			fmt.Println(tt.storage.GetGauge("test"))
		})
	}
}

func TestSaver(t *testing.T) {
	saver, _ := NewSaver("test.json", false, time.Second*2000)
	tests := []struct {
		name  string
		saver *Saver
		value int64
	}{
		{
			name:  "test set counter",
			saver: saver,
			value: 5, // errors.New("float"),
		},
		{
			name:  "test set counter",
			saver: saver,
			value: 50, // errors.New("float"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.saver.SetCounter("test", tt.value)
			fmt.Println(tt.saver.GetCounter("test"))
		})
	}
}

func TestSavereGauge(t *testing.T) {
	saver, _ := NewSaver("test.json", false, time.Second*2000)
	tests := []struct {
		name  string
		saver *Saver
		value float64
	}{
		{
			name:  "test set counter",
			saver: saver,
			value: 5., // errors.New("float"),
		},
		{
			name:  "test set counter",
			saver: saver,
			value: 50., // errors.New("float"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.saver.SetGauge("test", tt.value)
			fmt.Println(tt.saver.GetGauge("test"))
		})
	}
}
