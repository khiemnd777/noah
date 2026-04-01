// scripts/create_module/templates/config_config.go.tmpl
package config

import (
	sharedconfig "github.com/khiemnd777/noah_api/shared/config"
	frameworkattribute "github.com/khiemnd777/noah_framework/modules/attribute/module"
	frameworkservice "github.com/khiemnd777/noah_framework/modules/attribute/service"
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

func (c *ModuleConfig) FrameworkModuleConfig() frameworkattribute.Config {
	cacheTTL := sharedconfig.Get().Cache.TTL
	return frameworkattribute.Config{
		AutoMigrate: c.Database.AutoMigrate,
		Cache: frameworkservice.CacheConfig{
			ShortTTL: cacheTTL.Short,
			LongTTL:  cacheTTL.Long,
		},
	}
}
