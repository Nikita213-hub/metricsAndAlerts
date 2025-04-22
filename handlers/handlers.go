package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Nikita213-hub/metricsAndAlerts/internal/models"
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

func (s *StorageHandlers) UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {
	var reqData models.Metrics
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		slog.Error(err.Error())
		w.WriteHeader(501)
		return
	}
	resData := models.Metrics{}
	switch reqData.MType {
	case "gauge":
		v, err := s.strg.UpadateGauge(reqData.ID, *reqData.Value)
		if err != nil {
			slog.Error("Error while updating metric", "metric_type", reqData.MType, "metric_id", reqData.ID)
			w.WriteHeader(501)
			return
		}
		resData.Value = &v
	case "counter":
		v, err := s.strg.UpadateCounter(reqData.ID, int64(*reqData.Delta))
		if err != nil {
			slog.Error("Error while updating metric", "metric_type", reqData.MType, "metric_id", reqData.ID)
			w.WriteHeader(501)
			return
		}
		vf := float64(v)
		resData.Value = &vf
	default:
		slog.Info("Incorrect metric type", "metric_type", reqData.MType)
		w.WriteHeader(501)
		return
	}
	resData.ID = reqData.ID
	resData.MType = reqData.MType
	if err := json.NewEncoder(w).Encode(resData); err != nil {
		slog.Info("Encoding", "metric_type", reqData.MType)
		w.WriteHeader(501)
		return
	}
}

func (s *StorageHandlers) GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	var reqData models.Metrics
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		slog.Error(err.Error())
		w.WriteHeader(501)
		return
	}
	resData := models.Metrics{}
	switch reqData.MType {
	case "gauge":
		v, err := s.strg.GetGauge(reqData.ID)
		if err != nil {
			slog.Error("Error while updating metric", "metric_type", reqData.MType, "metric_id", reqData.ID)
			w.WriteHeader(501)
			return
		}
		resData.Value = &v
	case "counter":
		v, err := s.strg.GetCounter(reqData.ID)
		if err != nil {
			slog.Error("Error while updating metric", "metric_type", reqData.MType, "metric_id", reqData.ID)
			w.WriteHeader(501)
			return
		}
		vf := float64(v)
		resData.Value = &vf
	default:
		slog.Info("Incorrect metric type", "metric_type", reqData.MType)
		w.WriteHeader(501)
		return
	}
	resData.ID = reqData.ID
	resData.MType = reqData.MType
	if err := json.NewEncoder(w).Encode(resData); err != nil {
		slog.Info("Encoding error", "metric_type", reqData.MType)
		w.WriteHeader(501)
		return
	}
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
