package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"runtime"
	"time"
)

func collectMetrics(jobs chan<- Metric, ctx context.Context, pollInterval *time.Duration) {
	metrics := []Metric{
		{ID: "Alloc", MType: "gauge"},
		{ID: "BuckHashSys", MType: "gauge"},
		{ID: "Frees", MType: "gauge"},
		{ID: "GCCPUFraction", MType: "gauge"},
		{ID: "GCSys", MType: "gauge"},
		{ID: "HeapAlloc", MType: "gauge"},
		{ID: "HeapIdle", MType: "gauge"},
		{ID: "HeapInuse", MType: "gauge"},
		{ID: "HeapObjects", MType: "gauge"},
		{ID: "HeapReleased", MType: "gauge"},
		{ID: "HeapSys", MType: "gauge"},
		{ID: "LastGC", MType: "gauge"},
		{ID: "Lookups", MType: "gauge"},
		{ID: "MCacheInuse", MType: "gauge"},
		{ID: "MCacheSys", MType: "gauge"},
		{ID: "MSpanInuse", MType: "gauge"},
		{ID: "MSpanSys", MType: "gauge"},
		{ID: "Mallocs", MType: "gauge"},
		{ID: "NextGC", MType: "gauge"},
		{ID: "NumForcedGC", MType: "gauge"},
		{ID: "NumGC", MType: "gauge"},
		{ID: "OtherSys", MType: "gauge"},
		{ID: "PauseTotalNs", MType: "gauge"},
		{ID: "StackInuse", MType: "gauge"},
		{ID: "StackSys", MType: "gauge"},
		{ID: "Sys", MType: "gauge"},
		{ID: "TotalAlloc", MType: "gauge"},
		{ID: "PollCount", MType: "counter"},
		{ID: "RandomValue", MType: "gauge"},
	}

	ticker := time.NewTicker(*pollInterval)
	go func() {
		pollCount := int64(0)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// Update runtime metrics
				var m runtime.MemStats
				runtime.ReadMemStats(&m)

				// Update metrics with current values
				metrics[0].Value = float64Ptr(float64(m.Alloc))
				metrics[1].Value = float64Ptr(float64(m.BuckHashSys))
				metrics[2].Value = float64Ptr(float64(m.Frees))
				metrics[3].Value = float64Ptr(m.GCCPUFraction)
				metrics[4].Value = float64Ptr(float64(m.GCSys))
				metrics[5].Value = float64Ptr(float64(m.HeapAlloc))
				metrics[6].Value = float64Ptr(float64(m.HeapIdle))
				metrics[7].Value = float64Ptr(float64(m.HeapInuse))
				metrics[8].Value = float64Ptr(float64(m.HeapObjects))
				metrics[9].Value = float64Ptr(float64(m.HeapReleased))
				metrics[10].Value = float64Ptr(float64(m.HeapSys))
				metrics[11].Value = float64Ptr(float64(m.LastGC))
				metrics[12].Value = float64Ptr(float64(m.Lookups))
				metrics[13].Value = float64Ptr(float64(m.MCacheInuse))
				metrics[14].Value = float64Ptr(float64(m.MCacheSys))
				metrics[15].Value = float64Ptr(float64(m.MSpanInuse))
				metrics[16].Value = float64Ptr(float64(m.MSpanSys))
				metrics[17].Value = float64Ptr(float64(m.Mallocs))
				metrics[18].Value = float64Ptr(float64(m.NextGC))
				metrics[19].Value = float64Ptr(float64(m.NumForcedGC))
				metrics[20].Value = float64Ptr(float64(m.NumGC))
				metrics[21].Value = float64Ptr(float64(m.OtherSys))
				metrics[22].Value = float64Ptr(float64(m.PauseTotalNs))
				metrics[23].Value = float64Ptr(float64(m.StackInuse))
				metrics[24].Value = float64Ptr(float64(m.StackSys))
				metrics[25].Value = float64Ptr(float64(m.Sys))
				metrics[26].Value = float64Ptr(float64(m.TotalAlloc))

				// Update poll count
				pollCount = 1
				metrics[27].Delta = &pollCount

				// Update random value
				randomValue := rand.Float64() * 100
				metrics[28].Value = &randomValue
				fmt.Println(pollCount)
				go func() {
					select {
					case <-ctx.Done():
						return
					default:
						for _, v := range metrics {
							select {
							case <-ctx.Done():
								return
							default:
								jobs <- v
							}
						}
					}
				}()
			}
		}
	}()
}

func float64Ptr(v float64) *float64 {
	return &v
}
