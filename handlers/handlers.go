package handlers

import (
	"fmt"
	"net/http"

	"github.com/Nikita213-hub/metricsAndAlerts/internal/helpers"
	"github.com/Nikita213-hub/metricsAndAlerts/internal/storage"
)

type StorageHandlers struct {
	strg storage.Storage
}

func NewStorageHandlers(strg storage.Storage) *StorageHandlers {
	return &StorageHandlers{
		strg: strg,
	}
}

func (s *StorageHandlers) GaugeHandler(w http.ResponseWriter, r *http.Request) {
	k, v, err := helpers.GetMetricDataFromUri[float64](r.URL.Path)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = s.strg.UpadateGauge(k, v)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	updatedValue, err := s.strg.GetGauge(k)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	fmt.Println(k, " was updated successfully to val = ", updatedValue)
	w.WriteHeader(http.StatusOK)
}

func (s *StorageHandlers) CounterHandler(w http.ResponseWriter, r *http.Request) {
	k, v, err := helpers.GetMetricDataFromUri[int64](r.URL.Path)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = s.strg.UpadateCounter(k, v)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	updatedValue, err := s.strg.GetCounter(k)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	fmt.Println(k, " was updated successfully to val = ", updatedValue)
	w.WriteHeader(http.StatusOK)
}
