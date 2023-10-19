package main

import (
	"os"

	"golang.org/x/exp/slog"

	"github.com/boskuv/finance-control_expenses-service/internal/config"
)

// func init() {
// 	viper.AddConfigPath("config")
// 	viper.SetConfigName("config")
// 	viper.SetConfigType("yaml")

// 	if err := viper.ReadInConfig(); err != nil {
// 		log.Fatal("Error reading config file", err)
// 	}
// }

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	log.Info("initializing server", slog.String("address", "Not implemented"))
	log.Debug("logger debug mode enabled")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
