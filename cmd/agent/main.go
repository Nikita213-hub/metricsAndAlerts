package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"time"
)

const (
	serverURL = "http://localhost:8080"
	// Using the same private key as in the server's hash middleware
	privateKey = "secretkey"
)

type Metric struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func generateHash(data []byte) string {
	h := hmac.New(sha256.New, []byte(privateKey))
	h.Write(data)
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}

func sendMetric(metric Metric) error {
	jsonData, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("error marshaling metric: %w", err)
	}

	hash := generateHash(jsonData)
	url := fmt.Sprintf("%s/update/", serverURL)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HashSHA256", hash)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func main() {
	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Parse command line flags
	pollInterval := flag.Duration("poll", 2*time.Second, "Poll interval")
	reportInterval := flag.Duration("report", 10*time.Second, "Report interval")
	flag.Parse()

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Create metrics
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

	pollCount := int64(0)
	ticker := time.NewTicker(*pollInterval)
	defer ticker.Stop()

	for {
		select {
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
			pollCount++
			metrics[27].Delta = &pollCount

			// Update random value
			randomValue := rand.Float64() * 100
			metrics[28].Value = &randomValue

			// Send metrics
			for _, metric := range metrics {
				if err := sendMetric(metric); err != nil {
					slog.Error("Failed to send metric", "error", err, "metric", metric.ID)
				}
			}

			// Sleep for report interval
			time.Sleep(*reportInterval)
		}
	}
}

func float64Ptr(v float64) *float64 {
	return &v
}
