package memstorage

import (
	"errors"

	"github.com/Nikita213-hub/metricsAndAlerts/internal/storage"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemStorage() storage.Storage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (ms *MemStorage) GetGauge(key string) (float64, error) {
	v, ok := ms.gauge[key]
	if !ok {
		return 0, errors.New("invalid metric")
	}
	return v, nil
}

func (ms *MemStorage) GetCounter(key string) (int64, error) {
	v, ok := ms.counter[key]
	if !ok {
		return 0, errors.New("invalid metric")
	}
	return v, nil
}

func (m *MemStorage) UpadateGauge(key string, val float64) error {
	m.gauge[key] = val
	_, ok := m.gauge[key]
	if !ok {
		return errors.New("error while updating gauge metric type occured")
	}
	return nil
}

func (m *MemStorage) UpadateCounter(key string, val int64) error {
	m.counter[key] += val
	_, ok := m.counter[key]
	if !ok {
		return errors.New("error while updating counter metric type occured")
	}
	return nil
}
