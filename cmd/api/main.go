package main

import (
	"log"
	"net/http"
	"os"

	"ef_mob_test_go/config"
	"ef_mob_test_go/internal/subscriptions/handler"
	"ef_mob_test_go/internal/subscriptions/repository"
	"ef_mob_test_go/internal/subscriptions/service"
	"ef_mob_test_go/pkg/logger"
	"ef_mob_test_go/pkg/postgres"
)

func main() {
	cfg, err := config.GetConfigByFilename(os.Getenv("config"))
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	appLog, err := logger.NewLogger(cfg)
	if err != nil {
		log.Fatalf("Error init logger %w", err)
	}
	defer appLog.Sync()

	appLog.Info("Service started")

	psqlDB, err := postgres.NewSqlDB(cfg)
	if err != nil {
		appLog.Fatalf("Error init postgres %w", err)
	}
	defer psqlDB.Close()

	subRepo := repository.NewSubRepository(psqlDB)
	subService := service.NewSubService(subRepo, appLog)
	subHandler := handler.NewSubHandler(subService, appLog)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /subscriptions", subHandler.CreateSub)

	// TODO: сделать graceful shutdown
	appLog.Infof("Server listening on %s", cfg.Server.Port)
	if err := http.ListenAndServe(cfg.Server.Port, mux); err != nil {
		appLog.Fatalf("Server error: %v", err)
	}
}
