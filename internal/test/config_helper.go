package test

import "github.com/podengo-project/idmsvc-backend/internal/config"

// Config for testing
func GetTestConfig() (cfg *config.Config) {
	cfg = &config.Config{}
	cfg.Application = config.Application{
		DomainRegTokenKey:      "random",
		PaginationDefaultLimit: 10,
		PaginationMaxLimit:     100,
		PathPrefix:             config.DefaultPathPrefix,
	}
	return config.Load(cfg)
}
