package impl

import (
	"testing"

	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/metrics"
	"github.com/podengo-project/idmsvc-backend/internal/test"
	"github.com/podengo-project/idmsvc-backend/internal/test/mock/interface/client/pendo"
	"github.com/podengo-project/idmsvc-backend/internal/test/mock/interface/client/rbac"
	// client_rbac "github.com/podengo-project/idmsvc-backend/internal/test/mock/interface/client/rbac"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestGuardNewHandler(t *testing.T) {
	assert.PanicsWithValue(t, "'cfg' is nil", func() {
		guardNewHandler(nil, nil, nil, nil, nil)
	})

	cfg := test.GetTestConfig()
	assert.PanicsWithValue(t, "'db' is nil", func() {
		guardNewHandler(&config.Config{}, nil, nil, nil, nil)
	})

	sqlMock, gormDB, err := test.NewSqlMock(&gorm.Session{SkipHooks: true})
	require.NoError(t, err)
	require.NotNil(t, sqlMock)
	require.NotNil(t, gormDB)
	assert.PanicsWithValue(t, "'m' is nil", func() {
		guardNewHandler(cfg, gormDB, nil, nil, nil)
	})

	m := &metrics.Metrics{}
	assert.PanicsWithValue(t, "'rbac' is nil", func() {
		guardNewHandler(cfg, gormDB, m, nil, nil)
	})

	rbacClient := rbac.NewRbac(t)
	assert.PanicsWithValue(t, "'pendo' is nil", func() {
		guardNewHandler(cfg, gormDB, m, rbacClient, nil)
	})

	pendoClient := pendo.NewPendo(t)
	assert.NotPanics(t, func() {
		guardNewHandler(cfg, gormDB, m, rbacClient, pendoClient)
	})

	rbacClient.AssertExpectations(t)
	pendoClient.AssertExpectations(t)
	require.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestNewHandler(t *testing.T) {
	cfg := test.GetTestConfig()
	sqlMock, gormDB, err := test.NewSqlMock(&gorm.Session{SkipHooks: true})
	require.NoError(t, err)
	require.NotNil(t, sqlMock)
	require.NotNil(t, gormDB)
	m := &metrics.Metrics{}
	rbacClient := rbac.NewRbac(t)
	pendoClient := pendo.NewPendo(t)
	assert.NotPanics(t, func() {
		require.NotNil(t, NewHandler(cfg, gormDB, m, rbacClient, pendoClient))
	})

	rbacClient.AssertExpectations(t)
	pendoClient.AssertExpectations(t)
	require.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestAppSecrets(t *testing.T) {
	_, gormDB, err := test.NewSqlMock(&gorm.Session{SkipHooks: true})
	require.NoError(t, err)
	cfg := test.GetTestConfig()

	rbacClient := rbac.NewRbac(t)
	pendoClient := pendo.NewPendo(t)
	handler := NewHandler(cfg, gormDB, &metrics.Metrics{}, rbacClient, pendoClient)
	app := handler.(*application)

	assert.NotEmpty(t, app.config.Secrets.DomainRegKey)
	assert.Equal(t, len(app.config.Secrets.DomainRegKey), 32)
}
