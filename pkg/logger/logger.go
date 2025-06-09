package logger

import (
	"log"
	"log/slog"
	"os"
)

func Init(env string) {
	var logger *slog.Logger

	switch env {
	case "local":
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "product":
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log.Fatalf("unknown env name: %s", env)
	}

	slog.SetDefault(logger)
}
