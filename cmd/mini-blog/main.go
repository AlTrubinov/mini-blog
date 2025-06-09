package main

import (
	"context"
	"log/slog"
	"mini-blog/internal/storage/postgres"
	"os"

	"mini-blog/internal/config"
	"mini-blog/pkg/logger"
)

func main() {
	cfg := config.NewConfig()

	logger.Init(cfg.Env)
	slog.Info("logger initialized")

	ctx := context.Background()
	storagePool, err := postgres.NewStorage(ctx, cfg.Database)
	if err != nil {
		slog.Error("storage initialize error:", err)
		os.Exit(1)
	}
	defer storagePool.Close()
	slog.Info("storage initialized")

	// TODO: init handler

	// TODO: run server
}
