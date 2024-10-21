package main

import (
	"log/slog"

	"github.com/podengo-project/idmsvc-backend/cmd/db-tool/cmd"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/logger"
)

func initLogSystem(name string, cfg *config.Config) {
	logger.LogBuildInfo("db-tool")
	logger.InitLogger(cfg)
	slog.SetDefault(slog.Default().With(slog.String("component", "db-tool")))
	cfg.Log(slog.Default())
}

func main() {
	cfg := config.Get()
	initLogSystem("db-tool", cfg)
	cmd.Execute()
}
