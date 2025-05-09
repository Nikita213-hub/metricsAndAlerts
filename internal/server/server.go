package server

import (
	"net/http"

	gzip_local "github.com/Nikita213-hub/metricsAndAlerts/cmd/gzip"
	hashsign "github.com/Nikita213-hub/metricsAndAlerts/cmd/hashSign"
	"github.com/Nikita213-hub/metricsAndAlerts/handlers"
	"github.com/Nikita213-hub/metricsAndAlerts/internal/logger"
)

type Server struct {
	server  *http.Server
	router  *http.ServeMux
	address string
	port    string
}

func NewServer(address, port string) *Server {
	return &Server{
		address: address,
		port:    port,
	}
}

func (s *Server) Start(handlers *handlers.StorageHandlers) error {
	s.router = http.NewServeMux()
	s.router.HandleFunc("POST /update/", hashsign.WithHash(handlers.UpdateMetricHandler))
	s.router.HandleFunc("POST /value/", handlers.GetMetricHandler)
	s.router.Handle("GET /", logger.WithLogger(gzip_local.WithGzip(handlers.GetAllMetricsHandler)))
	s.server = &http.Server{
		Addr:    s.address + s.port,
		Handler: s.router,
	}

	return s.server.ListenAndServe()
}
