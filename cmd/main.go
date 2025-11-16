package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/example/avito_test_task_internship/configs"
	"github.com/example/avito_test_task_internship/internal/application"
	"github.com/example/avito_test_task_internship/internal/infrasctucture"
	presentation_http "github.com/example/avito_test_task_internship/internal/presentation/http"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := configs.NewConfig()
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	slog.SetDefault(logger)

	logger.Info("Starting merch store api", slog.String("env", cfg.Env))
	logger.Debug("Debug messages are enabled")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout*time.Second)
	defer cancel()

	connStr, err := cfg.BuildPGConnString()
	if err != nil {
		log.Fatalf("Error building connection to database string: %v", err)
	}
	database := infrasctucture.NewPostgresDb(ctx, connStr)
	defer database.Close()

	application.InitRegistry(database)

	r := chi.NewRouter()
	presentation_http.RegisterRoutes(r)

	http.ListenAndServe(cfg.HTTPServerConfig.Address, r)
}
