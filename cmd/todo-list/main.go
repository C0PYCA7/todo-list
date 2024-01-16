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
	"todo-list/internal/http-server/handlers/list"
	"todo-list/internal/http-server/handlers/signIn"
	"todo-list/internal/http-server/handlers/task"
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

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Post("/auth", auth.New(log, database))
	router.Post("/sign_in", signIn.New(log, database))
	router.Post("/newTask", task.New(log, database))
	router.Post("/list", list.New(log, database))

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
