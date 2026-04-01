package main

import (
	entsql "entgo.io/ent/dialect/sql"

	sharedapp "github.com/khiemnd777/noah_api/shared/app"
	"github.com/khiemnd777/noah_api/shared/db/ent"
	"github.com/khiemnd777/noah_api/shared/middleware"
	"github.com/khiemnd777/noah_api/shared/module"
	"github.com/khiemnd777/noah_api/shared/utils"

	"github.com/khiemnd777/noah_api/modules/auditlog/config"
	"github.com/khiemnd777/noah_api/modules/auditlog/ent/bootstrap"
	"github.com/khiemnd777/noah_api/modules/auditlog/ent/generated"
	"github.com/khiemnd777/noah_api/modules/auditlog/handler"
	"github.com/khiemnd777/noah_api/modules/auditlog/repository"
	"github.com/khiemnd777/noah_api/modules/auditlog/service"
	sharedGenerated "github.com/khiemnd777/noah_api/shared/db/ent/generated"
	frameworkapp "github.com/khiemnd777/noah_framework/pkg/app"
	frameworkdb "github.com/khiemnd777/noah_framework/pkg/db"
)

func main() {
	module.StartModule(module.ModuleOptions[config.ModuleConfig]{
		ConfigPath: utils.GetModuleConfigPath("auditlog"),
		ModuleName: "auditlog",
		InitEntClient: func(client frameworkdb.Client, cfg *config.ModuleConfig) (any, error) {
			return bootstrap.EntBootstrapFromDatabase(client, func(drv *entsql.Driver) any {
				return generated.NewClient(generated.Driver(drv))
			}, cfg.Database.AutoMigrate)
		},
		InitSharedEntClient: func(client frameworkdb.Client, cfg *config.ModuleConfig) (any, error) {
			return ent.EntBootstrapFromDatabase(client, func(drv *entsql.Driver) any {
				return sharedGenerated.NewClient(sharedGenerated.Driver(drv))
			}, cfg.SharedDatabase.AutoMigrate)
		},
		OnRegistry: func(app frameworkapp.Application, deps *module.ModuleDeps[config.ModuleConfig]) {
			repo := repository.NewAuditLogRepository(deps.Ent.(*generated.Client), deps.SharedEnt.(*sharedGenerated.Client), deps)
			svc := service.NewAuditLogService(repo, deps)
			svc.InitPubSubEvents()
			h := handler.NewAuditLogHandler(svc)
			h.RegisterRoutes(sharedapp.Group(app, utils.GetModuleRoute(deps.Config.Server.Route), middleware.RequireAuth()))
		},
	})
}
