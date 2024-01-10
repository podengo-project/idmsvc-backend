package logger

import (
	"os"
	"runtime/debug"
	"strings"

	"github.com/podengo-project/idmsvc-backend/internal/config"
	"golang.org/x/exp/slog"
)

// If you want to learn more about slog visit:
// https://betterstack.com/community/guides/logging/logging-in-go/

const (
	LevelTrace  = slog.Level(-8)
	LevelDebug  = slog.LevelDebug
	LevelInfo   = slog.LevelInfo
	LevelNotice = slog.Level(2)
	LevelWarn   = slog.LevelWarn
	LevelError  = slog.LevelError
)

type AppHandler struct {
	handler slog.Handler
}

type Clonable interface {
	Clone() interface{}
}

// Early logging setup so we can use slog, this will just log to stderr.
// This will change once the configuration has been parsed and we setup the
// logger accordingly.
func init() {
	h := slog.NewTextHandler(os.Stderr, nil)
	slog.SetDefault(slog.New(h))
}

func InitLogger(cfg *config.Config) {
	if cfg == nil {
		panic("'cfg' cannot be nil")
	}

	var h slog.Handler

	globalLevel := new(slog.LevelVar)
	// set default to warning
	globalLevel.Set(LevelWarn)

	LevelNames := map[slog.Leveler]string{
		LevelTrace:  "TRACE",
		LevelNotice: "NOTICE",
	}

	opts := slog.HandlerOptions{
		AddSource: false,
		Level:     globalLevel,
		// This will print TRACE and NOTICE in logs nicely
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := LevelNames[level]
				if !exists {
					levelLabel = level.String()
				}

				a.Value = slog.StringValue(levelLabel)
			}

			return a
		},
	}

	if cfg.Logging.Console {
		h = slog.NewTextHandler(
			os.Stderr,
			&opts,
		)
	} else {
		h = slog.NewJSONHandler(
			os.Stderr,
			&opts,
		)
	}
	slog.SetDefault(slog.New(h))

	// set global log level
	lvl := strings.ToUpper(cfg.Logging.Level)

	switch {
	case lvl == "TRACE":
		globalLevel.Set(LevelTrace)
	case lvl == "DEBUG":
		globalLevel.Set(LevelDebug)
	case lvl == "INFO":
		globalLevel.Set(LevelInfo)
	case lvl == "NOTICE":
		globalLevel.Set(LevelNotice)
	case lvl == "WARN":
		globalLevel.Set(LevelWarn)
	case lvl == "ERROR":
		globalLevel.Set(LevelError)
	default:
		globalLevel.Set(LevelWarn)
	}
}

func LogBuildInfo(msg string) {
	var (
		version    string = "unknown"
		revision   string = "unknown"
		commitTime string = "unknown"
		dirty      bool   = false
	)
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	// version is set to "(devel)" when building from a git checkout
	if info.Main.Version != "" {
		version = info.Main.Version
	}

	for _, s := range info.Settings {
		if s.Value == "" {
			continue
		}
		switch s.Key {
		case "vcs.modified":
			dirty = s.Value == "true"
		case "vcs.revision":
			revision = s.Value
		case "vcs.time":
			commitTime = s.Value
		}
	}
	slog.Info(
		msg,
		slog.String("Version", version),
		slog.String("Commit", revision),
		slog.String("CommitTime", commitTime),
		slog.Bool("Dirty", dirty),
	)
}
