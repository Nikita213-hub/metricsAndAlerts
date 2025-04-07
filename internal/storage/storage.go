package storage

type Storage interface {
	GetGauge(key string) (float64, error)
	GetCounter(key string) (int64, error)
	UpadateGauge(key string, val float64) error
	UpadateCounter(key string, val int64) error
	GetAllMetrics() (map[string]float64, map[string]int64, error)
}
