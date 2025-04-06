package main

import (
	"math/rand/v2"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var PollInterval = 2 * time.Second
var ReportInterval = 10 * time.Second

type MetWrapper struct {
	Rt             runtime.MemStats
	CPUutilization float64
	TotalMemory    float64
	FreeMemory     float64
	PollCount      int64
	RandomValue    float64
}

type Metrics struct {
	mem MetWrapper
	mx  sync.Mutex
}

func newMetrics() *Metrics {
	return &Metrics{
		mx: sync.Mutex{},
	}
}

func (m *Metrics) StartMetricsCollection() error {
	for {
		m.mx.Lock()
		runtime.ReadMemStats(&m.mem.Rt)
		m.mem.CPUutilization = 100 * (1 - float64(m.mem.Rt.HeapIdle)/float64(m.mem.Rt.HeapSys))
		m.mem.TotalMemory = float64(12)
		m.mem.FreeMemory = float64(23)
		m.mem.RandomValue = rand.Float64()
		m.mx.Unlock()
		time.Sleep(PollInterval)
	}
}

func (m *Metrics) Lock() {
	m.mx.Lock()
}

func (m *Metrics) Unlock() {
	m.mx.Unlock()
}

func (m *Metrics) GetMetrics() MetWrapper {
	m.mx.Lock()
	defer m.mx.Unlock()
	return m.mem
}

func (ms *Metrics) GetMap() map[string]float64 {
	ms.Lock()
	defer ms.Unlock()
	return map[string]float64{
		"Alloc":         float64(ms.mem.Rt.Alloc),
		"BuckHashSys":   float64(ms.mem.Rt.BuckHashSys),
		"Frees":         float64(ms.mem.Rt.Frees),
		"GCCPUFraction": ms.mem.Rt.GCCPUFraction,
		"GCSys":         float64(ms.mem.Rt.GCSys),
		"HeapAlloc":     float64(ms.mem.Rt.HeapAlloc),
		"HeapIdle":      float64(ms.mem.Rt.HeapIdle),
		"HeapInuse":     float64(ms.mem.Rt.HeapInuse),
		"HeapObjects":   float64(ms.mem.Rt.HeapObjects),
		"HeapReleased":  float64(ms.mem.Rt.HeapReleased),
		"HeapSys":       float64(ms.mem.Rt.HeapSys),
		"LastGC":        float64(ms.mem.Rt.LastGC),
		"Lookups":       float64(ms.mem.Rt.Lookups),
		"MCacheInuse":   float64(ms.mem.Rt.MCacheInuse),
		"MCacheSys":     float64(ms.mem.Rt.MCacheSys),
		"MSpanInuse":    float64(ms.mem.Rt.MSpanInuse),
		"MSpanSys":      float64(ms.mem.Rt.MSpanSys),
		"Mallocs":       float64(ms.mem.Rt.Mallocs),
		"NextGC":        float64(ms.mem.Rt.NextGC),
		"NumForcedGC":   float64(ms.mem.Rt.NumForcedGC),
		"NumGC":         float64(ms.mem.Rt.NumGC),
		"OtherSys":      float64(ms.mem.Rt.OtherSys),
		"PauseTotalNs":  float64(ms.mem.Rt.PauseTotalNs),
		"StackInuse":    float64(ms.mem.Rt.StackInuse),
		"StackSys":      float64(ms.mem.Rt.StackSys),
		"Sys":           float64(ms.mem.Rt.Sys),
		"TotalAlloc":    float64(ms.mem.Rt.TotalAlloc),
		"RandomValue":   ms.mem.RandomValue,
	}
}

func sendReport(url string) error {
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func strartRepoerting(metrics *Metrics, host string) {
	for {
		m := metrics.GetMap()
		for k, v := range m {
			url := host + "/update/gauge/" + k + "/" + strconv.FormatFloat(v, 'f', -1, 64)
			sendReport(url)
		}
		sendReport(host + "/update/counter/PollCount/1")
		time.Sleep(ReportInterval)
	}
}

func main() {
	metrics := newMetrics()
	var wg sync.WaitGroup
	wg.Add(2)
	go metrics.StartMetricsCollection()
	go strartRepoerting(metrics, "http://localhost:8080")
	go func() {
		time.Sleep(60 * time.Second)
		wg.Done()
		wg.Done()
	}()
	wg.Wait()
}
