package main

import (
	"github.com/khiemnd777/noah_api/modules/attribute/config"
	"github.com/khiemnd777/noah_api/shared/app"
	"github.com/khiemnd777/noah_api/shared/middleware"
	"github.com/khiemnd777/noah_api/shared/module"
	sharedruntime "github.com/khiemnd777/noah_api/shared/runtime"
	"github.com/khiemnd777/noah_api/shared/utils"
	frameworkattribute "github.com/khiemnd777/noah_framework/modules/attribute/module"
	frameworkapp "github.com/khiemnd777/noah_framework/pkg/app"
)

func main() {
	configPath, err := sharedruntime.ModuleConfigPath("attribute")
	if err != nil {
		configPath = utils.GetModuleConfigPath("attribute")
	}

	module.StartModule(module.ModuleOptions[config.ModuleConfig]{
		ConfigPath: configPath,
		ModuleName: "attribute",
		OnRegistry: func(moduleApp frameworkapp.Application, deps *module.ModuleDeps[config.ModuleConfig]) {
			if err := frameworkattribute.Register(
				app.Group(moduleApp, utils.GetModuleRoute(deps.Config.Server.Route), middleware.RequireAuth()),
				frameworkattribute.Options{
					Config:   deps.Config.FrameworkModuleConfig(),
					Database: deps.DBClient,
				},
			); err != nil {
				panic(err)
			}
		},
	})
}
