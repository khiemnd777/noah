package main

import (
	entsql "entgo.io/ent/dialect/sql"

	sharedapp "github.com/khiemnd777/noah_api/shared/app"
	"github.com/khiemnd777/noah_api/shared/db/ent"
	"github.com/khiemnd777/noah_api/shared/db/ent/generated"
	"github.com/khiemnd777/noah_api/shared/middleware"
	"github.com/khiemnd777/noah_api/shared/utils"

	"github.com/khiemnd777/noah_api/modules/user/config"
	"github.com/khiemnd777/noah_api/modules/user/handler"
	"github.com/khiemnd777/noah_api/modules/user/repository"
	"github.com/khiemnd777/noah_api/modules/user/service"
	"github.com/khiemnd777/noah_api/shared/module"
	frameworkapp "github.com/khiemnd777/noah_framework/pkg/app"
	frameworkdb "github.com/khiemnd777/noah_framework/pkg/db"
)

func main() {
	module.StartModule(module.ModuleOptions[config.ModuleConfig]{
		ConfigPath: utils.GetModuleConfigPath("user"),
		ModuleName: "user",
		InitEntClient: func(client frameworkdb.Client, cfg *config.ModuleConfig) (any, error) {
			return ent.EntBootstrapFromDatabase(client, func(drv *entsql.Driver) any {
				return generated.NewClient(generated.Driver(drv))
			}, cfg.Database.AutoMigrate)
		},
		OnRegistry: func(app frameworkapp.Application, deps *module.ModuleDeps[config.ModuleConfig]) {
			repo := repository.NewUserRepository(deps.Ent.(*generated.Client))
			svc := service.NewUserService(repo, deps)
			h := handler.NewUserHandler(svc, deps)
			h.RegisterRoutes(sharedapp.Group(app, utils.GetModuleRoute(deps.Config.Server.Route), middleware.RequireAuth()))
		},
	})
}
