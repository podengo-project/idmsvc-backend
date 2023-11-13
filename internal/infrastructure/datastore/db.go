package datastore

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/logger"
	"golang.org/x/exp/slog"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const DbMigrationPath = "./scripts/db/migrations"

func getUrl(config *config.Config) string {
	var sslStr string
	if config.Database.CACertPath == "" {
		sslStr = "sslmode=disable"
	} else {
		sslStr = fmt.Sprintf("sslmode=verify-full sslrootcert=%s", config.Database.CACertPath)
	}
	return fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s %s",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
		sslStr,
	)
}

func NewDB(cfg *config.Config) (db *gorm.DB) {
	var err error
	dbUrl := getUrl(cfg)

	if db, err = gorm.Open(pg.Open(dbUrl),
		&gorm.Config{
			Logger:                 logger.NewGormLog(true),
			SkipDefaultTransaction: true,
			// CreateBatchSize:        50,
			TranslateError: true,
		}); err != nil {
		slog.Error("Error creating database connector", slog.Any("error", err))
		return nil
	}
	return db
}

func NewDbMigration(config *config.Config) (db *gorm.DB, m *migrate.Migrate, err error) {
	var sqlDb *sql.DB
	dbURL := getUrl(config)
	sqlDb, err = sql.Open("postgres", dbURL)
	if err != nil {
		return nil, nil, fmt.Errorf("could not connect to database: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDb, &postgres.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("could not get database driver: %w", err)
	}

	if m, err = migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", DbMigrationPath),
		"postgres",
		driver); err != nil {
		return nil, nil, fmt.Errorf("could not create migration instance: %w", err)
	}

	return db, m, err
}

func Close(db *gorm.DB) {
	var (
		sqlDB *sql.DB
		err   error
	)
	if db == nil {
		slog.Warn("Close called with db=nil connector")
		return
	}
	if sqlDB, err = db.DB(); err != nil {
		slog.Error("Error retrieving the sql driver", slog.Any("error", err))
		return
	}
	if err = sqlDB.Close(); err != nil {
		slog.Error("Error closing database connector", slog.Any("error", err))
		return
	}
}
