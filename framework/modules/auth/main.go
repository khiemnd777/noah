package main

import (
	entsql "entgo.io/ent/dialect/sql"
	"github.com/khiemnd777/noah_framework/modules/auth/config"
	"github.com/khiemnd777/noah_framework/modules/auth/handler"
	"github.com/khiemnd777/noah_framework/modules/auth/repository"
	"github.com/khiemnd777/noah_framework/modules/auth/service"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/app"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/db/ent"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/db/ent/generated"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/module"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/utils"
	frameworkapp "github.com/khiemnd777/noah_framework/pkg/app"
	frameworkdb "github.com/khiemnd777/noah_framework/pkg/db"
)

func main() {
	module.StartModule(module.ModuleOptions[config.ModuleConfig]{
		ConfigPath: utils.GetModuleConfigPath("auth"),
		ModuleName: "auth",
		InitEntClient: func(client frameworkdb.Client, cfg *config.ModuleConfig) (any, error) {
			return ent.EntBootstrapFromDatabase(client, func(drv *entsql.Driver) any {
				return generated.NewClient(generated.Driver(drv))
			}, cfg.Database.AutoMigrate)
		},
		OnRegistry: func(moduleApp frameworkapp.Application, deps *module.ModuleDeps[config.ModuleConfig]) {
			authSecret := utils.GetAuthSecret()
			db := deps.Ent.(*generated.Client)
			repo := repository.NewAuthRepository(db)
			svc := service.NewAuthService(repo, authSecret)
			h := handler.NewAuthHandler(svc)
			h.RegisterRoutes(app.Group(moduleApp, utils.GetModuleRoute(deps.Config.Server.Route)))
		},
	})
}
