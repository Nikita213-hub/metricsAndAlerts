package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func getMetricDataFromUri[V int64 | float64](url string) (key string, val V, err error) {
	splittedUrl := strings.Split(url, "/")
	if len(splittedUrl) < 5 {
		return "", 0, errors.New("incorrect data has been provided")
	}
	metricType := splittedUrl[2]
	metricKey := splittedUrl[3]
	metricVal := splittedUrl[4]

	if metricKey == "" {
		return "", 0, errors.New("provide metric's key")
	}

	if metricType == "counter" {
		metricValInt, err := strconv.Atoi(metricVal)
		if err != nil {
			return "", 0, err
		}
		return metricKey, V(metricValInt), nil
	} else if metricType == "gauge" {
		metricValFloat, err := strconv.ParseFloat(metricVal, 64)
		if err != nil {
			return "", 0, err
		}
		return metricKey, V(metricValFloat), nil
	} else {
		return "", 0, err
	}
}

func (m *MemStorage) UpadateGaugeHandler(key string, val float64) error {
	m.gauge[key] = val
	_, ok := m.gauge[key]
	if !ok {
		return errors.New("error while updating gauge metric type occured")
	}
	return nil
}
func (m *MemStorage) UpadateCounterHandler(key string, val int64) error {
	m.counter[key] += val
	_, ok := m.counter[key]
	if !ok {
		return errors.New("error while updating counter metric type occured")
	}
	return nil
}

func main() {
	mux := http.NewServeMux()
	ms := NewMemStorage()
	mux.HandleFunc("/update/gauge/", func(w http.ResponseWriter, r *http.Request) {
		k, v, err := getMetricDataFromUri[float64](r.URL.Path)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = ms.UpadateGaugeHandler(k, v)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println(k, " was updated successfully to val = ", v)
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/update/counter/", func(w http.ResponseWriter, r *http.Request) {
		k, v, err := getMetricDataFromUri[int64](r.URL.Path)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = ms.UpadateCounterHandler(k, v)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println(k, " was updated successfully to val = ", v)
		w.WriteHeader(http.StatusOK)
	})
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
