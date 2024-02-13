package impl

import (
	"testing"

	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/metrics"
	"github.com/podengo-project/idmsvc-backend/internal/test"
	"github.com/podengo-project/idmsvc-backend/internal/test/mock/interface/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNewHandler(t *testing.T) {
	sqlMock, gormDB, err := test.NewSqlMock(&gorm.Session{SkipHooks: true})
	inventoryMock := client.NewHostInventory(t)
	require.NoError(t, err)
	require.NotNil(t, sqlMock)
	require.NotNil(t, gormDB)
	assert.Panics(t, func() {
		NewHandler(nil, nil, nil, nil)
	})
	assert.Panics(t, func() {
		NewHandler(&config.Config{}, nil, nil, nil)
	})
	cfg := test.GetTestConfig()
	assert.NotPanics(t, func() {
		NewHandler(cfg, gormDB, &metrics.Metrics{}, inventoryMock)
	})
}

func TestAppSecrets(t *testing.T) {
	_, gormDB, err := test.NewSqlMock(&gorm.Session{SkipHooks: true})
	inventoryMock := client.NewHostInventory(t)
	require.NoError(t, err)
	cfg := test.GetTestConfig()

	handler := NewHandler(cfg, gormDB, &metrics.Metrics{}, inventoryMock)
	app := handler.(*application)

	assert.NotEmpty(t, app.config.Secrets.DomainRegKey)
	assert.Equal(t, len(app.config.Secrets.DomainRegKey), 32)
}
