package context

import (
	"context"
	"log/slog"
)

type keySlog string

// CtxWithLog create a context that contain the specified
// *slog.Logger value.
func CtxWithLog(ctx context.Context, log *slog.Logger) context.Context {
	key := keySlog("log")
	if ctx == nil {
		panic("'ctx' is nil")
	}
	if log == nil {
		panic("'log' is nil")
	}
	return context.WithValue(ctx, key, log)
}

// LogFromCtx get a log from a specified context.
func LogFromCtx(ctx context.Context) *slog.Logger {
	key := keySlog("log")
	if ctx == nil {
		panic("'ctx' is nil")
	}
	l, ok := ctx.Value(key).(*slog.Logger)
	if !ok {
		panic("'log' could not be read")
	}
	return l
}
