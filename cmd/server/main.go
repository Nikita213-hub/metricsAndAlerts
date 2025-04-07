package main

import (
	"flag"
	"fmt"

	"github.com/Nikita213-hub/metricsAndAlerts/cmd/flags"
	"github.com/Nikita213-hub/metricsAndAlerts/handlers"
	"github.com/Nikita213-hub/metricsAndAlerts/internal/server"
	memstorage "github.com/Nikita213-hub/metricsAndAlerts/internal/storage/memStorage"
)

func main() {
	strg := memstorage.NewMemStorage()
	handlers := handlers.NewStorageHandlers(strg)

	addr := flags.NewAddress("localhost", "8080")
	_ = flag.Value(addr)
	flag.Var(addr, "a", "Address in host:port fmt")
	flag.Parse()
	fmt.Println(addr)
	server := server.NewServer(addr.GetHost(), ":"+addr.GetPort())
	err := server.Start(handlers)
	if err != nil {
		panic(err)
	}
}
