package context

import (
	"context"
	"testing"

	"github.com/podengo-project/idmsvc-backend/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCtxWithDB(t *testing.T) {
	ctx := context.TODO()
	_, dbMock, err := test.NewSqlMock(nil)
	require.NoError(t, err)
	assert.NotPanics(t, func() {
		ctx = CtxWithDB(ctx, dbMock)
	})
}

func TestDBFromCtx(t *testing.T) {
	ctx := context.TODO()
	assert.PanicsWithValue(t, "'db' could not be read", func() {
		_ = DBFromCtx(ctx)
	})

	_, dbMock, err := test.NewSqlMock(nil)
	require.NoError(t, err)
	ctx = CtxWithDB(ctx, dbMock)

	db := DBFromCtx(ctx)
	require.NotNil(t, db)
}
