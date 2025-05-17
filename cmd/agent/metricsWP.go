package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type MetricsWorkerPool struct {
	jobs       chan Metric
	results    chan *Result
	workersNum int
	Stop       context.CancelFunc
}

func NewMetricsWP(rm int) *MetricsWorkerPool {
	return &MetricsWorkerPool{
		jobs:       make(chan Metric, 10),
		results:    make(chan *Result, 10),
		workersNum: rm,
	}
}

func generateHash(data []byte, privateKey string) string {
	h := hmac.New(sha256.New, []byte(privateKey))
	h.Write(data)
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}

func sendMetric(metric Metric, privateKey string, serverURL string) (*http.Response, error) {
	jsonData, err := json.Marshal(metric)
	if err != nil {
		return nil, fmt.Errorf("error marshaling metric: %w", err)
	}

	hash := generateHash(jsonData, privateKey)
	url := fmt.Sprintf("%s/update/", serverURL)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HashSHA256", hash)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, err
}

func (mwp *MetricsWorkerPool) Run(privateKey string, serverURL string, ctx context.Context) {
	var wg sync.WaitGroup
	go func() {
		for v := range mwp.results {
			_ = v
		}
	}()
	for v := range mwp.workersNum {
		_ = v
		wg.Add(1)
		go func() {
			defer wg.Done()
			for mj := range mwp.jobs {
				select {
				case <-ctx.Done():
					return
				default:
					res, err := sendMetric(mj, privateKey, serverURL)
					if err != nil {
						fmt.Println(err) // add error handling
					} else {
						mwp.results <- &Result{
							res,
							err,
						}
					}
				}
			}
		}()
	}
	defer close(mwp.jobs)
	defer close(mwp.results)
	wg.Wait()
}
