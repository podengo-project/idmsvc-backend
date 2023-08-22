package logger

import (
	"os"

	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/rs/zerolog"
)

type apiZeroLog struct {
	logger zerolog.Logger
}

func NewApiLogger(cfg *config.Config) zerolog.Logger {
	var output *os.File
	if cfg.Logging.Console {
		output = os.Stdout
	} else {
		output = os.Stderr
	}
	lvl, err := zerolog.ParseLevel(cfg.Logging.Level)
	if err != nil {
		panic(err)
	}
	return zerolog.New(output).Level(lvl)
}
