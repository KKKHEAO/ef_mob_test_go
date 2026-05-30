package main

import (
	"ef_mob_test_go/config"
	"ef_mob_test_go/pkg/logger"
	"ef_mob_test_go/pkg/postgres"
	"log"
	"os"
)

func main() {
	cfg, err := config.GetConfigByFilename(os.Getenv("config"))
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}
	logger, err := logger.NewLogger(cfg)
	if err != nil {
		log.Fatalf("Error init logger %w", err)
	}
	logger.Info("Service started")

	psqlDB, err := postgres.NewSqlDB(cfg)
	if err != nil {
		logger.Fatalf("Error init postgres %w", err)
	}
	defer psqlDB.Close()
}
