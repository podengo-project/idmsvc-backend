package test

import (
	"fmt"
	"os"

	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/secrets"
)

// Config for testing
func GetTestConfig() (cfg *config.Config) {
	if cfgpath, ok := os.LookupEnv("CONFIG_PATH"); !ok || cfgpath == "" {
		panic("CONFIG_PATH is not set or is empty")
	}
	cfg = &config.Config{}
	config.Load(cfg)

	if err := config.Validate(cfg); err != nil {
		panic(fmt.Errorf("Invalid configuration: %w", err))
	}

	sec, err := secrets.NewAppSecrets(cfg.Application.MainSecret)
	if err != nil {
		panic(err)
	}
	cfg.Secrets = *sec

	// override some default settings
	cfg.Application.MainSecret = secrets.GenerateRandomMainSecret()
	cfg.Application.TokenExpirationTimeSeconds = 3600
	cfg.Application.HostconfJwkValidity = config.DefaultHostconfJwkValidity
	cfg.Application.HostconfJwkRenewalThreshold = config.DefaultHostconfJwkRenewalThreshold
	cfg.Application.PaginationDefaultLimit = 10
	cfg.Application.PaginationMaxLimit = 100

	return cfg
}
