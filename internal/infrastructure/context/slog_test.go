package context

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCtxWithLog(t *testing.T) {
	ctx := context.TODO()
	assert.PanicsWithValue(t, "'log' is nil", func() {
		ctx = CtxWithLog(ctx, nil)
	})
	require.NotNil(t, ctx)

	assert.NotPanics(t, func() {
		ctx = CtxWithLog(ctx, slog.Default())
	})
	require.NotNil(t, ctx)
}

func TestLogFromCtx(t *testing.T) {
	var log *slog.Logger
	assert.PanicsWithValue(t, "'log' could not be read", func() {
		log = LogFromCtx(context.TODO())
	})
	require.Nil(t, log)

	assert.NotPanics(t, func() {
		log = LogFromCtx(CtxWithLog(context.TODO(), slog.Default()))
	})
	require.NotNil(t, log)
}
