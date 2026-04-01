package main

import (
	entsql "entgo.io/ent/dialect/sql"
	"github.com/khiemnd777/noah_framework/modules/notification/config"
	"github.com/khiemnd777/noah_framework/modules/notification/handler"
	"github.com/khiemnd777/noah_framework/modules/notification/repository"
	"github.com/khiemnd777/noah_framework/modules/notification/service"
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
		ConfigPath: utils.GetModuleConfigPath("notification"),
		ModuleName: "notification",
		InitEntClient: func(client frameworkdb.Client, cfg *config.ModuleConfig) (any, error) {
			return ent.EntBootstrapFromDatabase(client, func(drv *entsql.Driver) any {
				return generated.NewClient(generated.Driver(drv))
			}, cfg.Database.AutoMigrate)
		},
		OnRegistry: func(moduleApp frameworkapp.Application, deps *module.ModuleDeps[config.ModuleConfig]) {
			repo := repository.NewNotificationRepository(deps.Ent.(*generated.Client), deps)
			svc := service.NewNotificationService(repo, deps)
			svc.InitPubSubEvents()
			h := handler.NewNotificationHandler(svc, deps)
			baseRoute := utils.GetModuleRoute(deps.Config.Server.Route)
			h.RegisterRoutes(app.Group(moduleApp, baseRoute, middleware.RequireAuth()))
			h.RegisterRoutesInternal(app.Group(moduleApp, baseRoute, middleware.RequireInternal()))
		},
	})
}
