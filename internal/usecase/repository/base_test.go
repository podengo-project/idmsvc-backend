package repository

// https://pkg.go.dev/github.com/stretchr/testify/suite

import (
	"bytes"
	"context"
	"log/slog"

	"github.com/DATA-DOG/go-sqlmock"
	app_context "github.com/podengo-project/idmsvc-backend/internal/infrastructure/context"
	"github.com/podengo-project/idmsvc-backend/internal/test"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type SuiteBase struct {
	suite.Suite
	DB        *gorm.DB
	LogBuffer bytes.Buffer
	Log       *slog.Logger
	Ctx       context.Context
	mock      sqlmock.Sqlmock
}

// https://pkg.go.dev/github.com/stretchr/testify/suite#SetupTestSuite
func (s *SuiteBase) SetupTest() {
	var err error
	s.mock, s.DB, err = test.NewSqlMock(&gorm.Session{
		SkipHooks: true,
	})
	if err != nil {
		s.Suite.FailNow("Error calling gorm.Open: %s", err.Error())
		return
	}
	s.LogBuffer.Reset()
	s.Log = slog.New(slog.NewTextHandler(&s.LogBuffer, nil))
	s.Ctx = app_context.CtxWithLog(
		app_context.CtxWithDB(context.Background(), s.DB),
		s.Log,
	)
}
