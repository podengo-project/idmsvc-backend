package test

import "github.com/podengo-project/idmsvc-backend/internal/config"

// Config for testing
func GetTestConfig() (cfg *config.Config) {
	cfg = &config.Config{}
	cfg = config.Load(cfg)
	// override some default settings
	cfg.Application.DomainRegTokenKey = "random"
	cfg.Application.PaginationDefaultLimit = 10
	cfg.Application.PaginationMaxLimit = 100
	return cfg
}
