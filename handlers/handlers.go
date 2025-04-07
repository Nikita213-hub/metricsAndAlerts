package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

func (s *StorageHandlers) UpdateGaugeHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *StorageHandlers) UpdateCounterHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *StorageHandlers) GetGaugeHandler(w http.ResponseWriter, r *http.Request) {
	k, err := helpers.GetMetricKeyFromUrl(r.URL.Path)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	v, err := s.strg.GetGauge(k)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatFloat(v, 'f', -1, 64)))
}

func (s *StorageHandlers) GetCounterHandler(w http.ResponseWriter, r *http.Request) {
	k, err := helpers.GetMetricKeyFromUrl(r.URL.Path)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	v, err := s.strg.GetCounter(k)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatInt(v, 10)))
}

func (s *StorageHandlers) GetAllMetricsHandler(w http.ResponseWriter, r *http.Request) {
	gauge, counter, err := s.strg.GetAllMetrics()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	htmlForm := "<html><body><h1>Metrics</h1><ul><h2>Gauge:</h2>"
	var sb strings.Builder
	sb.WriteString(htmlForm)
	for k, v := range gauge {
		sb.WriteString(fmt.Sprintf("<li>%s: %f</li>", k, v))
	}
	sb.WriteString("<h2>Counter:</h2><ul>")
	for k, v := range counter {
		sb.WriteString(fmt.Sprintf("<li>%s: %d</li>", k, v))
	}
	sb.WriteString("</ul></body></html>")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(sb.String()))
}
