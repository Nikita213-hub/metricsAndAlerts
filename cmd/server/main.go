package main

import (
	"github.com/Nikita213-hub/metricsAndAlerts/handlers"
	"github.com/Nikita213-hub/metricsAndAlerts/internal/server"
	memstorage "github.com/Nikita213-hub/metricsAndAlerts/internal/storage/memStorage"
)

func main() {
	strg := memstorage.NewMemStorage()
	handlers := handlers.NewStorageHandlers(strg)
	server := server.NewServer("localhost", ":8080")
	err := server.Start(handlers)
	if err != nil {
		panic(err)
	}
}
