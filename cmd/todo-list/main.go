package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log/slog"
	"net/http"
	"os"
	"todo-list/internal/config"
	"todo-list/internal/database/postgres"
	"todo-list/internal/http-server/handlers/auth"
)

func main() {
	cfg := config.MustLoad()

	log := newLogger()

	log.Info("starting todo-list", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	database, err := postgres.New(cfg.Database)
	if err != nil {
		log.Error("failed to init database", err)
		os.Exit(1)
	}
	_ = database

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	router.Post("/auth", auth.New(log, database))

	log.Info("starting server", slog.String("address", cfg.HttpServer.Address))

	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server is stopped")
}

func newLogger() *slog.Logger {
	var log *slog.Logger

	log = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	return log
}
