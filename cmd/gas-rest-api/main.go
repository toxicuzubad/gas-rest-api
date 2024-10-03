package main

import (
	"fmt"
	"gas-rest-api/internal/config"
	"gas-rest-api/internal/lib/logger/sl"
	"gas-rest-api/internal/storage/sqlite"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// Reading config file.
	cfg := config.MustLoad()

	// Setup logger.
	log := setupLogger(cfg.Env)
	log.Info("starting gas-rest-api service", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// Init storage.
	fmt.Printf(cfg.StoragePath)
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	id, err := storage.SaveGuitar("Ibanez", "RG2770", "Made in Japan. Superstrat", "IBZ12321")
	if err != nil {
		log.Error("failed to save new guitar", sl.Err(err))
		os.Exit(1)
	}

	log.Info("saved guitar", slog.Int64("id", id))

	_ = storage

	// TODO: init router

	// TODO: run server

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
