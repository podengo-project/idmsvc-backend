package test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewSqlMock(session *gorm.Session) (sqlmock.Sqlmock, *gorm.DB, error) {
	sqlDB, sqlMock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		TranslateError:         true,
	})
	if err != nil {
		return nil, nil, err
	}
	if session != nil {
		gormDB = gormDB.Session(session)
	}

	return sqlMock, gormDB, nil
}
