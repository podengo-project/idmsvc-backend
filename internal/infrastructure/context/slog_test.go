package context

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCtxWithLog(t *testing.T) {
	assert.PanicsWithValue(t, "'ctx' is nil", func() {
		_ = CtxWithLog(nil, nil)
	})

	ctx := context.TODO()
	assert.PanicsWithValue(t, "'log' is nil", func() {
		_ = CtxWithLog(ctx, nil)
	})
	require.NotNil(t, ctx)

	assert.NotPanics(t, func() {
		ctx = CtxWithLog(ctx, slog.Default())
	})
	require.NotNil(t, ctx)
}

func TestLogFromCtx(t *testing.T) {
	var log *slog.Logger

	assert.PanicsWithValue(t, "'ctx' is nil", func() {
		log = LogFromCtx(nil)
	})

	ctx := context.TODO()
	assert.PanicsWithValue(t, "'log' could not be read", func() {
		log = LogFromCtx(ctx)
	})
	require.Nil(t, log)

	assert.NotPanics(t, func() {
		ctx = CtxWithLog(ctx, slog.Default())
		log = LogFromCtx(ctx)
	})
	require.NotNil(t, log)
}
