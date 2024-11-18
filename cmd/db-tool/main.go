package main

import (
	"github.com/podengo-project/idmsvc-backend/cmd/db-tool/cmd"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/logger"
)

const component = "db-tool"

func main() {
	logger.LogBuildInfo(component)
	cfg := config.Get()
	logger.InitLogger(cfg, component)
	defer logger.DoneLogger()

	cmd.Execute()
}
