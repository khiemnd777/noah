package main

import (
	entsql "entgo.io/ent/dialect/sql"
	"github.com/khiemnd777/noah_api/modules/photo/config"
	"github.com/khiemnd777/noah_api/modules/photo/jobs"
	sharedapp "github.com/khiemnd777/noah_api/shared/app"
	"github.com/khiemnd777/noah_api/shared/cron"
	"github.com/khiemnd777/noah_api/shared/db/ent"
	"github.com/khiemnd777/noah_api/shared/middleware"

	"github.com/khiemnd777/noah_api/shared/db/ent/generated"

	"github.com/khiemnd777/noah_api/modules/photo/handler"
	"github.com/khiemnd777/noah_api/modules/photo/repository"
	"github.com/khiemnd777/noah_api/modules/photo/service"
	"github.com/khiemnd777/noah_api/shared/module"
	"github.com/khiemnd777/noah_api/shared/utils"
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
		OnRegistry: func(app frameworkapp.Application, deps *module.ModuleDeps[config.ModuleConfig]) {
			repo := repository.NewPhotoRepository(deps.Ent.(*generated.Client), deps)
			svc := service.NewPhotoService(repo, deps)
			h := handler.NewPhotoHandler(svc, deps)
			h.RegisterRoutes(sharedapp.Group(app, utils.GetModuleRoute(deps.Config.Server.Route), middleware.RequireAuth()))

			cron.RegisterJob(jobs.NewClearDeletedPhotoJob(svc))
		},
	})
}
