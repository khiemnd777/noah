package main

import (
	entsql "entgo.io/ent/dialect/sql"
	"github.com/khiemnd777/noah_api/modules/attribute/config"
	"github.com/khiemnd777/noah_api/modules/attribute/handler"
	"github.com/khiemnd777/noah_api/modules/attribute/repository"
	"github.com/khiemnd777/noah_api/modules/attribute/service"
	"github.com/khiemnd777/noah_api/shared/app"
	"github.com/khiemnd777/noah_api/shared/db/ent"
	"github.com/khiemnd777/noah_api/shared/db/ent/generated"
	"github.com/khiemnd777/noah_api/shared/middleware"
	"github.com/khiemnd777/noah_api/shared/module"
	"github.com/khiemnd777/noah_api/shared/utils"
	frameworkapp "github.com/khiemnd777/noah_framework/pkg/app"
	frameworkdb "github.com/khiemnd777/noah_framework/pkg/db"
)

func main() {
	module.StartModule(module.ModuleOptions[config.ModuleConfig]{
		ConfigPath: utils.GetModuleConfigPath("attribute"),
		ModuleName: "attribute",
		InitEntClient: func(client frameworkdb.Client, cfg *config.ModuleConfig) (any, error) {
			return ent.EntBootstrapFromDatabase(client, func(drv *entsql.Driver) any {
				return generated.NewClient(generated.Driver(drv))
			}, cfg.Database.AutoMigrate)
		},
		OnRegistry: func(moduleApp frameworkapp.Application, deps *module.ModuleDeps[config.ModuleConfig]) {
			repo := repository.NewAttributeRepository(deps.Ent.(*generated.Client), deps)
			svc := service.NewAttributeService(repo, deps)
			h := handler.NewAttributeHandler(svc, deps)
			h.RegisterRoutes(app.Group(moduleApp, utils.GetModuleRoute(deps.Config.Server.Route), middleware.RequireAuth()))

			repoOption := repository.NewAttributeOptionRepository(deps.Ent.(*generated.Client), deps)
			svcOption := service.NewAttributeOptionService(repoOption, deps)
			hOption := handler.NewAttributeOptionHandler(svcOption, deps)
			hOption.RegisterRoutes(app.Group(moduleApp, utils.GetModuleRoute(deps.Config.Server.Route), middleware.RequireAuth()))
		},
	})
}
