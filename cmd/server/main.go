package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Nikita213-hub/metricsAndAlerts/cmd/flags"
	"github.com/Nikita213-hub/metricsAndAlerts/handlers"
	"github.com/Nikita213-hub/metricsAndAlerts/internal/server"
	memstorage "github.com/Nikita213-hub/metricsAndAlerts/internal/storage/memStorage"
)

type AddressVal struct {
	Host string
	Port string
}

func main() {
	strg := memstorage.NewMemStorage()
	handlers := handlers.NewStorageHandlers(strg)
	var address AddressVal
	if envAddr, ok := os.LookupEnv("ADDRESS"); ok {
		splt := strings.Split(envAddr, ":")
		if len(splt) != 2 {
			panic(errors.New("incorrect env var"))
		}
		address.Host = splt[0]
		address.Port = splt[1]
	} else {
		addr := flags.NewAddress("localhost", "8080")
		_ = flag.Value(addr)
		flag.Var(addr, "a", "Address in host:port fmt")
		flag.Parse()
		address.Host = addr.GetHost()
		address.Port = addr.GetPort()
	}
	fmt.Println(address)
	server := server.NewServer(address.Host, ":"+address.Port)
	err := server.Start(handlers)
	if err != nil {
		panic(err)
	}
}
