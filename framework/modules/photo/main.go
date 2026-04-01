package main

import (
	entsql "entgo.io/ent/dialect/sql"
	"github.com/khiemnd777/noah_framework/modules/photo/config"
	"github.com/khiemnd777/noah_framework/modules/photo/handler"
	"github.com/khiemnd777/noah_framework/modules/photo/jobs"
	"github.com/khiemnd777/noah_framework/modules/photo/repository"
	"github.com/khiemnd777/noah_framework/modules/photo/service"
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
		ConfigPath: utils.GetModuleConfigPath("photo"),
		ModuleName: "photo",
		InitEntClient: func(client frameworkdb.Client, cfg *config.ModuleConfig) (any, error) {
			return ent.EntBootstrapFromDatabase(client, func(drv *entsql.Driver) any {
				return generated.NewClient(generated.Driver(drv))
			}, cfg.Database.AutoMigrate)
		},
		OnRegistry: func(moduleApp frameworkapp.Application, deps *module.ModuleDeps[config.ModuleConfig]) {
			repo := repository.NewPhotoRepository(deps.Ent.(*generated.Client), deps)
			svc := service.NewPhotoService(repo, deps)
			h := handler.NewPhotoHandler(svc, deps)
			h.RegisterRoutes(app.Group(moduleApp, utils.GetModuleRoute(deps.Config.Server.Route), middleware.RequireAuth()))

			cron.RegisterJob(jobs.NewClearDeletedPhotoJob(svc))
		},
	})
}
