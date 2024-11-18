package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime/debug"
	"strings"

	cw "github.com/RedHatInsights/cloudwatch"
	"github.com/aws/aws-sdk-go/aws"
	aws_creds "github.com/aws/aws-sdk-go/aws/credentials"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/podengo-project/idmsvc-backend/internal/config"
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

var (
	cloudwatchCloser io.Closer    = nil
	oldSlog          *slog.Logger = nil
)

// Early logging setup so we can use slog, this will just log to stderr.
// This will change once the configuration has been parsed and we setup the
// logger accordingly.
func init() {
	h := slog.NewTextHandler(os.Stderr, nil)
	slog.SetDefault(slog.New(h))
}

func InitLogger(cfg *config.Config, componentName string) {
	var ok bool

	if cfg == nil {
		panic("'cfg' cannot be nil")
	}

	// This will be always logged to console on service startup
	slog.Info(
		"idmsvc-backend",
		slog.Group("Logging",
			slog.Bool("Console", cfg.Logging.Console),
			slog.String("Type", cfg.Logging.Type),
			slog.String("Level", cfg.Logging.Level),
			slog.Bool("Location", cfg.Logging.Location),
		),
	)

	globalLevel := new(slog.LevelVar)
	// set default to warning
	globalLevel.Set(LevelWarn)

	LevelNames := map[slog.Leveler]string{
		LevelTrace:  "TRACE",
		LevelNotice: "NOTICE",
	}

	opts := slog.HandlerOptions{
		AddSource: cfg.Logging.Location,
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

	metaHandler := NewSlogMetaHandler()
	if cfg.Logging.Type == "cloudwatch" {
		// Stderr json handler
		stderr_handler := slog.NewJSONHandler(
			os.Stderr,
			&opts,
		)
		metaHandler.Add(stderr_handler)

		// Cloudwatch json handler
		awsConf := aws.NewConfig().WithCredentials(
			aws_creds.NewStaticCredentials(
				cfg.Logging.Cloudwatch.Key,
				cfg.Logging.Cloudwatch.Secret,
				"",
			),
		).WithRegion(cfg.Logging.Cloudwatch.Region)

		// This will panic on error
		aws_sess := aws_session.Must(aws_session.NewSession(awsConf))

		group := cw.NewGroup(
			cfg.Logging.Cloudwatch.Group,
			cloudwatchlogs.New(aws_sess),
		)
		w, err := group.Create(cfg.Logging.Cloudwatch.Stream)
		if err != nil {
			panic(fmt.Errorf("cloudwatch: Failed to create group: %w", err))
		}

		cloudwatch_handler := slog.NewJSONHandler(
			w,
			&opts,
		)
		if cloudwatchCloser, ok = w.(io.Closer); !ok {
			slog.Warn("cloudwatchCloser is nil; log output could stay un-flushed and be lost")
		}
		metaHandler.Add(cloudwatch_handler)
	} else if cfg.Logging.Console {
		h := slog.NewTextHandler(
			os.Stderr,
			&opts,
		)
		metaHandler.Add(h)
	} else {
		h := slog.NewJSONHandler(
			os.Stderr,
			&opts,
		)
		metaHandler.Add(h)
	}
	if oldSlog == nil {
		oldSlog = slog.Default()
	} else {
		slog.Warn("InitLogger called twice")
	}
	slog.SetDefault(slog.New(metaHandler))

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

	slog.SetDefault(slog.Default().With(slog.String("component", componentName)))
	cfg.Log(slog.Default())
	slog.Log(
		context.Background(),
		LevelNotice,
		"Logging successfully initialized",
	)
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

// DoneLogger should be called before exiting any process that
// invokes InitLogger.  e.g. via `defer logger.DoneLogger()`.
func DoneLogger() {
	slog.SetDefault(oldSlog)
	if cloudwatchCloser != nil {
		if err := cloudwatchCloser.Close(); err != nil {
			slog.Error(err.Error())
		}
	}
}
