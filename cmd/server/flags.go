package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/bobgromozeka/metrics/internal/server"
)

var startupConfig server.StartupConfig

const (
	Address         = "ADDRESS"
	StoreInterval   = "STORE_INTERVAL"
	FileStoragePath = "FILE_STORAGE_PATH"
	Restore         = "RESTORE"
	DatabaseDsn     = "DATABASE_DSN"
)

func parseFlags() {
	flag.StringVar(&startupConfig.ServerAddr, "a", ":8080", "address and port to run server")
	flag.UintVar(&startupConfig.StoreInterval, "i", 300, "Interval of storing metrics to file")
	flag.StringVar(&startupConfig.FileStoragePath, "f", "/tmp/metrics-db.json", "Metrics file storage path")
	flag.BoolVar(&startupConfig.Restore, "r", true, "Restore metrics from file on server start or not")
	flag.StringVar(
		&startupConfig.DatabaseDsn, "d", "",
		"Postgresql data source name (connection string like postgres://username:password@localhost:5432/database_name)",
	)

	flag.Parse()
}

func parseEnv() {
	if addr := os.Getenv(Address); addr != "" {
		startupConfig.ServerAddr = addr
	}

	if interval := os.Getenv(StoreInterval); interval != "" {
		parsedInterval, err := strconv.Atoi(interval)
		if err != nil {
			log.Fatalln(StoreInterval+" parsing error ", err)
		}
		if parsedInterval < 0 {
			log.Fatalln(StoreInterval + " must be greater or equal 0")
		}
		startupConfig.StoreInterval = uint(parsedInterval)
	}

	if path := os.Getenv(FileStoragePath); path != "" {
		startupConfig.FileStoragePath = path
	}

	if r := os.Getenv(Restore); r == "false" || r == "0" {
		startupConfig.Restore = false
	}

	if dsn := os.Getenv(DatabaseDsn); dsn != "" {
		startupConfig.DatabaseDsn = dsn
	}
}

func setupConfiguration() {
	parseFlags()
	parseEnv()
}
