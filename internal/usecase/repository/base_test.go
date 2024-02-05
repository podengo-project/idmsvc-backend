package repository

// https://pkg.go.dev/github.com/stretchr/testify/suite

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/podengo-project/idmsvc-backend/internal/test"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type SuiteBase struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
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
}
