package entities

type ManagerValues interface {
	GetGauge(string) (string, bool)
	SetGauge(string, float64)
	GetCounter(string) (string, bool)
	SetCounter(string, int64)
}

type ManagerJSON interface {
	AllMetricsJSON() []MetricsJSON
	SetMetrics([]MetricsJSON) error
}

type Storage interface {
	ManagerJSON
	ManagerValues

	AllMetrics() map[string]string
	Ping() bool
}
