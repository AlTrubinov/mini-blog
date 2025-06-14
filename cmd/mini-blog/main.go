package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"mini-blog/internal/config"
	"mini-blog/internal/lib/auth"
	"mini-blog/internal/mini-blog/handlers/users/login"
	"mini-blog/internal/mini-blog/handlers/users/notes/create"
	"mini-blog/internal/mini-blog/handlers/users/notes/delete"
	"mini-blog/internal/mini-blog/handlers/users/notes/get"
	"mini-blog/internal/mini-blog/handlers/users/notes/list"
	"mini-blog/internal/mini-blog/handlers/users/notes/update"
	"mini-blog/internal/mini-blog/handlers/users/registration"
	"mini-blog/internal/mini-blog/handlers/users/token"
	"mini-blog/internal/storage/postgres"
	"mini-blog/pkg/logger"
)

func main() {
	cfg := config.NewConfig()

	setupLogger(cfg.Env)
	slog.Info("logger initialized")

	ctx := context.Background()
	storagePool, err := postgres.NewStorage(ctx, cfg.Database)
	if err != nil {
		slog.Error("storage initialize error:", err)
		os.Exit(1)
	}
	defer storagePool.Close()
	slog.Info("storage initialized")

	authManager := auth.NewTokenManager(cfg.Auth)
	slog.Info("auth initialized")

	router := chi.NewRouter()
	router.Use(logger.ApiInfo)

	router.Post("/users", registration.New(storagePool, authManager))
	router.Post("/login", login.New(storagePool, authManager))
	router.Post("/refresh", token.New(authManager))
	router.Group(func(r chi.Router) {
		r.Use(authManager.Middleware)

		r.Route("/users/{user_id}/notes", func(r chi.Router) {
			r.Post("/", create.New(storagePool))
			r.Get("/", list.New(storagePool))
			r.Get("/{note_id}", get.New(storagePool))
			r.Put("/{note_id}", update.New(storagePool))
			r.Delete("/{note_id}", delete.New(storagePool))
		})
	})

	slog.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		slog.Error("failed to start server")
	}

	slog.Error("server stopped")
}

func setupLogger(env string) {
	var lgr *slog.Logger

	switch env {
	case "local":
		lgr = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "product":
		lgr = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log.Fatalf("unknown env name: %s", env)
	}

	slog.SetDefault(lgr)
}
