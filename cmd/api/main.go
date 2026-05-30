// @title           Subscription API
// @version         1.0
// @description     API для управления подписками пользователей
// @host            localhost:3000
// @BasePath        /
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ef_mob_test_go/config"
	"ef_mob_test_go/internal/middleware"
	"ef_mob_test_go/internal/subscriptions/handler"
	"ef_mob_test_go/internal/subscriptions/repository"
	"ef_mob_test_go/internal/subscriptions/service"
	"ef_mob_test_go/pkg/logger"
	"ef_mob_test_go/pkg/postgres"

	_ "ef_mob_test_go/docs"

	httpswagger "github.com/swaggo/http-swagger/v2"
)

func main() {
	cfg, err := config.GetConfigByFilename(os.Getenv("config"))
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	appLog, err := logger.NewLogger(cfg)
	if err != nil {
		log.Fatalf("Error init logger: %v", err)
	}
	defer appLog.Sync()

	appLog.Info("Service started")

	psqlDB, err := postgres.NewSqlDB(cfg)
	if err != nil {
		appLog.Fatalf("Error init postgres: %v", err)
	}
	defer psqlDB.Close()

	subRepo := repository.NewSubRepository(psqlDB)
	subService := service.NewSubService(subRepo, appLog)
	subHandler := handler.NewSubHandler(subService, appLog)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /subscriptions", subHandler.CreateSub)
	mux.HandleFunc("GET /subscriptions", subHandler.ListSubs)
	mux.HandleFunc("GET /subscriptions/cost", subHandler.CalculateCost)
	mux.HandleFunc("GET /subscriptions/{id}", subHandler.GetSubByID)
	mux.HandleFunc("PUT /subscriptions/{id}", subHandler.UpdateSubByID)
	mux.HandleFunc("DELETE /subscriptions/{id}", subHandler.DeleteSubByID)
	mux.Handle("GET /swagger/{path...}", httpswagger.Handler())
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		if err := psqlDB.PingContext(r.Context()); err != nil {
			http.Error(w, "database unavailable", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	srv := &http.Server{
		Addr:         cfg.Server.Port,
		Handler:      middleware.RequestID(middleware.Logging(appLog)(mux)),
		ReadTimeout:  cfg.Server.ReadTimeout * time.Second,
		WriteTimeout: cfg.Server.WriteTimeout * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		appLog.Infof("Server listening on %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLog.Fatalf("Server error: %v", err)
		}
	}()

	sig := <-quit
	appLog.Infof("Received signal: %v. Shutting down...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLog.Fatalf("Server forced to shutdown: %v", err)
	}

	appLog.Info("Server exited gracefully")
}
