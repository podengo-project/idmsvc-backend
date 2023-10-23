package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/rs/zerolog"
	"gorm.io/gorm/logger"
)

type gormZerolog struct {
	logger zerolog.Logger
}

func NewGormLog(cfg *config.Config) logger.Interface {
	var o *os.File
	if cfg.Logging.Console {
		o = os.Stdout
	} else {
		o = os.Stderr
	}

	lvl, err := zerolog.ParseLevel(cfg.Logging.Level)
	if err != nil {
		panic(err)
	}
	l := zerolog.New(o).
		Level(lvl).
		With().
		Timestamp().
		Logger()

	return &gormZerolog{
		logger: l,
	}
}

func (l *gormZerolog) LogMode(level logger.LogLevel) logger.Interface {
	switch level {
	case logger.Silent:
		{
			l.logger.Level(zerolog.NoLevel)
		}
	case logger.Info:
		{
			l.logger.Level(zerolog.InfoLevel)
		}
	case logger.Warn:
		{
			l.logger.Level(zerolog.WarnLevel)
		}
	case logger.Error:
		{
			l.logger.Level(zerolog.ErrorLevel)
		}
	}
	return l
}

func (l *gormZerolog) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.logger.GetLevel() <= zerolog.InfoLevel {
		l.logger.Info().Msgf(msg, args...)
	}
}
func (l *gormZerolog) Warn(ctx context.Context, msg string, args ...interface{}) {
	if l.logger.GetLevel() <= zerolog.WarnLevel {
		l.logger.Warn().Msgf(msg, args...)
	}
}
func (l *gormZerolog) Error(ctx context.Context, msg string, args ...interface{}) {
	if l.logger.GetLevel() <= zerolog.ErrorLevel {
		l.logger.Error().Msgf(msg, args...)
	}
}
func (l *gormZerolog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sqlStr, iRowsAffected := fc()
	lvl := l.logger.GetLevel()
	if lvl > zerolog.TraceLevel {
		return
	}
	elapsedTime := time.Since(begin)
	if err != nil {
		l.logger.Err(err)
	}
	l.logger.Trace().
		Str("statement", sqlStr).
		Int64("rowsAffected", iRowsAffected).
		Str("el", fmt.Sprintf("%v ns", elapsedTime.Nanoseconds())).
		Send()
}
