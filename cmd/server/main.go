package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Nikita213-hub/metricsAndAlerts/cmd/flags"
	"github.com/Nikita213-hub/metricsAndAlerts/handlers"
	"github.com/Nikita213-hub/metricsAndAlerts/internal/db"
	"github.com/Nikita213-hub/metricsAndAlerts/internal/helpers"
	"github.com/Nikita213-hub/metricsAndAlerts/internal/logger"
	"github.com/Nikita213-hub/metricsAndAlerts/internal/server"
	memstorage "github.com/Nikita213-hub/metricsAndAlerts/internal/storage/memStorage"
)

type AddressVal struct {
	Host string
	Port string
}

func main() {
	logger.Init()
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
	var sInterval uint64
	if envSInterval, ok := os.LookupEnv("STORE_INTERVAL"); ok {
		val, err := strconv.Atoi(envSInterval)
		if err != nil {
			panic(errors.New("incorrect env var"))
		}
		sInterval = uint64(val)
	} else {
		sIntervalFlag := flags.NewFlag("10")
		flag.Var(sIntervalFlag, "i", "Save metrics in file interval in seconds")
		flag.Parse()
		val, err := strconv.Atoi(sIntervalFlag.String())
		if err != nil {
			panic(errors.New("incorrect env var"))
		}
		sInterval = uint64(val)
	}
	// what must constructor return interface or structure
	w, ok := strg.(*memstorage.MemStorage)
	if ok {
		fmt.Println(w)
		w.EnableSaves("metrics.json", time.Duration(sInterval)*time.Second)
	}

	ctx, _ := context.WithCancel(context.Background()) // shouldnt ignore cancel sig, but now i dont have functional to process it :(
	db := db.NewDatabase("localhost", "postgres", "123", "metrics_and_alerts", false)
	err := helpers.WithRetry(ctx, 5, 3*time.Second, db.Run)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
	slog.Info("Server is started listening", "address", address.Host+":"+address.Port)
	server := server.NewServer(address.Host, ":"+address.Port)
	err = helpers.WithRetry(ctx, 5, 3*time.Second, func() error {
		return server.Start(handlers)
	})
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
}
