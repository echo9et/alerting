package storage

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
