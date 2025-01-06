package storage

import "fmt"

type MemStorage struct {
	Counters map[string]int64
	Gauges   map[string]float64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Counters: make(map[string]int64),
		Gauges:   make(map[string]float64),
	}
}

func (m *MemStorage) AllMetrics() map[string]string {
	out := make(map[string]string)
	for k, v := range m.Counters {
		out[k] = fmt.Sprint(v)
	}
	for k, v := range m.Gauges {
		out[k] = fmt.Sprint(v)
	}
	return out
}
