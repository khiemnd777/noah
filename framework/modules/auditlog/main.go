package main

import (
	entsql "entgo.io/ent/dialect/sql"
	"github.com/khiemnd777/noah_framework/modules/auditlog/config"
	"github.com/khiemnd777/noah_framework/modules/auditlog/ent/bootstrap"
	"github.com/khiemnd777/noah_framework/modules/auditlog/ent/generated"
	"github.com/khiemnd777/noah_framework/modules/auditlog/handler"
	"github.com/khiemnd777/noah_framework/modules/auditlog/repository"
	"github.com/khiemnd777/noah_framework/modules/auditlog/service"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/app"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/db/ent"
	sharedGenerated "github.com/khiemnd777/noah_framework/internal/legacy/shared/db/ent/generated"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/middleware"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/module"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/utils"
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
		OnRegistry: func(moduleApp frameworkapp.Application, deps *module.ModuleDeps[config.ModuleConfig]) {
			repo := repository.NewAuditLogRepository(deps.Ent.(*generated.Client), deps.SharedEnt.(*sharedGenerated.Client), deps)
			svc := service.NewAuditLogService(repo, deps)
			svc.InitPubSubEvents()
			h := handler.NewAuditLogHandler(svc)
			h.RegisterRoutes(app.Group(moduleApp, utils.GetModuleRoute(deps.Config.Server.Route), middleware.RequireAuth()))
		},
	})
}
