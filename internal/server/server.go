package server

import (
	"net/http"

	"github.com/Nikita213-hub/metricsAndAlerts/handlers"
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
	s.router.HandleFunc("POST /update/gauge/", handlers.GaugeHandler)
	s.router.HandleFunc("POST /update/counter/", handlers.CounterHandler)
	s.server = &http.Server{
		Addr:    s.address + s.port,
		Handler: s.router,
	}

	return s.server.ListenAndServe()
}
