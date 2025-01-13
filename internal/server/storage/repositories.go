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

func (m *MemStorage) GetCounter(name string) (string, bool) {
	value, ok := m.Counters[name]
	if ok {
		return fmt.Sprint(value), true
	}
	return "", false
}

func (m *MemStorage) SetCounter(name string, iValue int64) {
	m.Counters[name] += iValue
}

func (m *MemStorage) GetGauge(name string) (string, bool) {
	value, ok := m.Gauges[name]
	if ok {
		return fmt.Sprint(value), true
	}
	return "", false
}

func (m *MemStorage) SetGauge(name string, fValue float64) {
	m.Gauges[name] = fValue
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
