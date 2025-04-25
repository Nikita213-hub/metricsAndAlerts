package memstorage

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/Nikita213-hub/metricsAndAlerts/internal/storage"
)

// add stopSaves chan, by singnaling in this chan we can terminate goroutine
// in enable saves, (it must be done via method StopSavingLogs)
type MemStorage struct {
	gmx     sync.Mutex
	cmx     sync.Mutex
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemStorage() storage.Storage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (ms *MemStorage) EnableSaves(dest string, interval time.Duration) error {
	fp, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	var writtenSize int64
	go func() {
		for {
			data, err := json.Marshal(struct {
				Counter map[string]int64   `json:"counter"`
				Gauge   map[string]float64 `json:"gauge"`
			}{
				Counter: ms.counter,
				Gauge:   ms.gauge,
			})
			if err != nil {
				slog.Info("Save metrics into file failure", "error", err.Error())
				continue
			}
			if writtenSize > 0 {
				fmt.Println("HELLO")
				err := fp.Truncate(0)
				if err != nil {
					slog.Info("Same metrics into file failure", "error", err.Error())
					continue
				}
				fp.Seek(0, 0)
			}
			n, err := fp.Write(data)
			if err != nil {
				slog.Info("Save metrics into file failure", "error", err.Error())
			}
			slog.Info("Metrics were saved into file", "bytes written", n)
			writtenSize = int64(n)
			time.Sleep(interval)
		}
	}()
	return nil
}

func (ms *MemStorage) GetGauge(key string) (float64, error) {
	ms.gmx.Lock()
	defer ms.gmx.Unlock()
	v, ok := ms.gauge[key]
	if !ok {
		return 0, errors.New("invalid metric")
	}
	return v, nil
}

func (ms *MemStorage) GetCounter(key string) (int64, error) {
	ms.cmx.Lock()
	defer ms.cmx.Unlock()
	v, ok := ms.counter[key]
	if !ok {
		return 0, errors.New("invalid metric")
	}
	return v, nil
}

func (m *MemStorage) UpadateGauge(key string, val float64) (float64, error) {
	m.gmx.Lock()
	defer m.gmx.Unlock()
	m.gauge[key] = val
	updated, ok := m.gauge[key]
	if !ok {
		return 0, errors.New("error while updating gauge metric type occured")
	}
	return updated, nil
}

func (m *MemStorage) UpadateCounter(key string, val int64) (int64, error) {
	m.cmx.Lock()
	defer m.cmx.Unlock()
	m.counter[key] += val
	updated, ok := m.counter[key]
	if !ok {
		return 0, errors.New("error while updating counter metric type occured")
	}
	return updated, nil
}

func (m *MemStorage) GetAllMetrics() (map[string]float64, map[string]int64, error) {
	m.gmx.Lock()
	m.cmx.Lock()
	defer m.gmx.Unlock()
	defer m.cmx.Unlock()
	return m.gauge, m.counter, nil
}
