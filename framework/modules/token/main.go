package main

import (
	entsql "entgo.io/ent/dialect/sql"
	"github.com/khiemnd777/noah_framework/modules/token/config"
	"github.com/khiemnd777/noah_framework/modules/token/handler"
	"github.com/khiemnd777/noah_framework/modules/token/jobs"
	"github.com/khiemnd777/noah_framework/modules/token/repository"
	"github.com/khiemnd777/noah_framework/modules/token/service"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/app"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/cron"
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
		ConfigPath: utils.GetModuleConfigPath("token"),
		ModuleName: "token",
		InitEntClient: func(client frameworkdb.Client, cfg *config.ModuleConfig) (any, error) {
			return ent.EntBootstrapFromDatabase(client, func(drv *entsql.Driver) any {
				return generated.NewClient(generated.Driver(drv))
			}, cfg.Database.AutoMigrate)
		},
		OnRegistry: func(moduleApp frameworkapp.Application, deps *module.ModuleDeps[config.ModuleConfig]) {
			authSecret := utils.GetAuthSecret()
			db := deps.Ent.(*generated.Client)
			repo := repository.NewTokenRepository(db)
			svc := service.NewTokenService(repo, authSecret)
			h := handler.NewTokenHandler(svc)
			h.RegisterRoutes(app.Group(moduleApp, utils.GetModuleRoute(deps.Config.Server.Route), middleware.RequireInternal()))

			cron.RegisterJob(jobs.NewClearResfreshTokenJob(svc))
		},
	})
}
