package main

import (
	"log"

	frameworkfolder "github.com/khiemnd777/noah_framework/modules/folder/module"
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
		Name:       "folder",
		ConfigPath: frameworkruntime.FrameworkModulePath("folder", "config.yaml"),
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
			return frameworkfolder.Register(router, frameworkfolder.Options{
				Database: deps.DB,
				Config: frameworkfolder.Config{
					Cache: frameworkfolder.CacheConfig{
						ShortTTL: deps.AppConfig.Cache.TTL.Short,
						LongTTL:  deps.AppConfig.Cache.TTL.Long,
					},
				},
			})
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
