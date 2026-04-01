package main

import (
	entsql "entgo.io/ent/dialect/sql"
	"github.com/khiemnd777/noah_framework/modules/realtime/config"
	"github.com/khiemnd777/noah_framework/modules/realtime/handler"
	"github.com/khiemnd777/noah_framework/modules/realtime/service"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/app"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/db/ent"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/db/ent/generated"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/middleware"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/module"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/utils"
	frameworkapp "github.com/khiemnd777/noah_framework/pkg/app"
	frameworkdb "github.com/khiemnd777/noah_framework/pkg/db"
)

func main() {
	module.StartModule(module.ModuleOptions[config.ModuleConfig]{
		ConfigPath: utils.GetModuleConfigPath("realtime"),
		ModuleName: "realtime",
		InitEntClient: func(client frameworkdb.Client, cfg *config.ModuleConfig) (any, error) {
			return ent.EntBootstrapFromDatabase(client, func(drv *entsql.Driver) any {
				return generated.NewClient(generated.Driver(drv))
			}, cfg.Database.AutoMigrate)
		},
		OnRegistry: func(moduleApp frameworkapp.Application, deps *module.ModuleDeps[config.ModuleConfig]) {
			hub := service.NewHub()
			hub.InitPubSubEvents()
			go hub.Run()

			jwtSecret := utils.GetAuthSecret()
			h := handler.NewHandler(hub, jwtSecret)
			baseRoute := utils.GetModuleRoute(deps.Config.Server.Route)
			h.RegisterRoutes(app.Group(moduleApp, baseRoute))
			h.RegisterInternalRoutes(app.Group(moduleApp, baseRoute, middleware.RequireInternal()))
		},
	})
}
