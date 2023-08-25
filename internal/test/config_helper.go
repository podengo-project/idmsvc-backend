package test

import "github.com/podengo-project/idmsvc-backend/internal/config"

// Config for testing
func GetTestConfig() (cfg *config.Config) {
	cfg = &config.Config{}
	return config.Load(cfg)
}
