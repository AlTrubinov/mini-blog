package main

import (
	"log/slog"

	"mini-blog/internal/config"
	"mini-blog/pkg/logger"
)

func main() {
	cfg := config.NewConfig()

	logger.Init(cfg.Env)
	slog.Info("logger initialized")

	// TODO: init storage

	// TODO: init handler

	// TODO: run server
}
