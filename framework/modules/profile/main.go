package main

import (
	"log"

	frameworkprofile "github.com/khiemnd777/noah_framework/modules/profile/module"
	frameworkapp "github.com/khiemnd777/noah_framework/pkg/app"
	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"
)

type moduleConfig struct {
	Server   frameworkruntime.ServerConfig   `yaml:"server"`
	External bool                            `yaml:"external"`
	Database frameworkruntime.DatabaseConfig `yaml:"database"`
}

func main() {
	err := frameworkruntime.StartModule[moduleConfig](frameworkruntime.ModuleOptions[moduleConfig]{
		Name:       "profile",
		ConfigPath: frameworkruntime.FrameworkModulePath("profile", "config.yaml"),
		GetServer: func(cfg *moduleConfig) frameworkruntime.ServerConfig {
			return cfg.Server
		},
		GetDatabase: func(cfg *moduleConfig) frameworkruntime.DatabaseConfig {
			return cfg.Database
		},
		OnRegister: func(_ frameworkapp.Application, deps *frameworkruntime.ModuleDeps[moduleConfig]) error {
			router := deps.App.Router().Mount(
				deps.Config.Server.Route,
				frameworkruntime.RequireAuth(deps.AppConfig.Auth.Secret),
			)
			return frameworkprofile.Register(router, frameworkprofile.Options{
				Database: deps.DB,
				Config: frameworkprofile.Config{
					Cache: frameworkprofile.CacheConfig{
						MediumTTL: deps.AppConfig.Cache.TTL.Medium,
					},
				},
			})
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
