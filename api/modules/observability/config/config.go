package config

import (
	sharedconfig "github.com/khiemnd777/noah_api/shared/config"
	frameworkobservability "github.com/khiemnd777/noah_framework/modules/observability/module"
	frameworkrepository "github.com/khiemnd777/noah_framework/modules/observability/repository"
)

type ModuleConfig struct {
	Server   sharedconfig.ServerConfig   `yaml:"server"`
	Database sharedconfig.DatabaseConfig `yaml:"database"`
}

func (c *ModuleConfig) GetServer() sharedconfig.ServerConfig {
	return c.Server
}

func (c *ModuleConfig) GetDatabase() sharedconfig.DatabaseConfig {
	return c.Database
}

func FrameworkModuleConfig() frameworkobservability.Config {
	cfg := sharedconfig.Get().Observability
	return frameworkobservability.Config{
		Loki: frameworkrepository.LokiConfig{
			BaseURL:        cfg.Loki.BaseURL,
			TenantID:       cfg.Loki.TenantID,
			BearerToken:    cfg.Loki.BearerToken,
			Timeout:        cfg.Loki.Timeout,
			StreamSelector: cfg.Loki.StreamSelector,
			MaxQueryLimit:  cfg.Loki.MaxQueryLimit,
		},
	}
}
