package logger

import (
	"github.com/hmsidm/internal/config"
	"github.com/rs/zerolog"
)

func InitLogger(cfg *config.Config) {
	if cfg == nil {
		panic("'cfg' cannot be nil")
	}

	lvl, err := zerolog.ParseLevel(cfg.Logging.Level)
	if err != nil {
		panic(err)
	}
	zerolog.SetGlobalLevel(lvl)
}
