package router

import (
	"testing"

	"github.com/labstack/echo/v4"
	api_private "github.com/podengo-project/idmsvc-backend/internal/test/mock/api/private"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGuardNewGroupPrivate(t *testing.T) {
	assert.PanicsWithValue(t, "'e' is nil", func() {
		guardNewGroupPrivate(nil, nil)
	})

	e := echo.New()
	g := e.Group("/private")
	assert.PanicsWithValue(t, "'apiPrivate' is nil", func() {
		guardNewGroupPrivate(g, nil)
	})

	apiPrivate := api_private.NewServerInterface(t)
	assert.NotPanics(t, func() {
		guardNewGroupPrivate(g, apiPrivate)
	})
}

func TestNewGroupPrivate(t *testing.T) {
	e := echo.New()
	require.NotNil(t, e)
	handlers := api_private.NewServerInterface(t)
	require.NotNil(t, handlers)
	assert.NotPanics(t, func() {
		require.NotNil(t, newGroupPrivate(e.Group("/private"), handlers))
	})
}
