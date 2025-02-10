package entities

type CounterStorage interface {
	GetCounter(string) (string, bool)
	SetCounter(string, int64)
}

type GaugeStorage interface {
	GetGauge(string) (string, bool)
	SetGauge(string, float64)
}

type WorkerJSON interface {
	AllMetricsJSON() []MetricsJSON
	SetMetrics([]MetricsJSON) error
}

type Storage interface {
	WorkerJSON
	CounterStorage
	GaugeStorage

	AllMetrics() map[string]string
	Ping() bool
}
