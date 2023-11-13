package logger

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"time"

	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// See https://gorm.io/docs/logger.html
type gormLogger struct {
	slogger                   *slog.Logger
	IgnoreRecordNotFoundError bool
}

// _LogCommon This function creates slog messages with correct source code locations
func (l *gormLogger) _LogCommon(
	ctx context.Context,
	level slog.Level,
	msg string,
	args ...interface{},
) {
	if !l.slogger.Enabled(ctx, level) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(4, pcs[:]) // skip [Callers, Infof]

	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(args...)
	_ = l.slogger.Handler().Handle(ctx, r)
}

// GORM uses these log messages in the form:
// Info(ctx, "wurst `%s` from %s\n", brot, utils.FileWithLineNum())
func (l *gormLogger) _LogMsg(
	ctx context.Context,
	level slog.Level,
	msg string,
	args ...interface{},
) {
	l._LogCommon(ctx, level, fmt.Sprintf(msg, args...))
}

func (l *gormLogger) _Log(
	ctx context.Context,
	level slog.Level,
	msg string,
	args ...interface{},
) {
	l._LogCommon(ctx, level, msg, args...)
}

func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *gormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	l._LogMsg(ctx, slog.LevelInfo, msg, args...)
}

func (l *gormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	l._LogMsg(ctx, slog.LevelWarn, msg, args...)
}

func (l *gormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	l._LogMsg(ctx, slog.LevelError, msg, args...)
}

func (l *gormLogger) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64),
	err error,
) {
	elapsedTime := time.Since(begin)

	if err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) ||
		!l.IgnoreRecordNotFoundError) {
		sql, rows := fc()

		l._Log(
			ctx,
			LevelTrace,
			err.Error(),
			slog.Any("error", err),
			slog.String("query", sql),
			slog.Duration("elapsed", elapsedTime),
			slog.Int64("rows", rows),
		)
	} else {
		sql, rows := fc()

		l._Log(
			ctx,
			LevelTrace,
			"SQL query executed",
			slog.String("query", sql),
			slog.Duration("elapsed", elapsedTime),
			slog.Int64("rows", rows),
		)
	}
}

func NewGormLog(ignoreRecordNotFound bool) logger.Interface {
	return &gormLogger{
		slogger:                   slog.Default(),
		IgnoreRecordNotFoundError: true,
	}
}
