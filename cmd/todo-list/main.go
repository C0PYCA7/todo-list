package main

import (
	"fmt"
	"log/slog"
	"os"
	"todo-list/internal/config"
	"todo-list/internal/database/postgres"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
	log := newLogger()
	log.Info("starting todo-list", slog.String("env", cfg.Env))
	database, err := postgres.New(cfg.Database)
	if err != nil {
		log.Error("failed to init database", err)
		os.Exit(1)
	}
	_ = database

}

func newLogger() *slog.Logger {
	var log *slog.Logger

	log = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	return log
}
