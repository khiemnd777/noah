package main

import (
	entsql "entgo.io/ent/dialect/sql"
	"github.com/khiemnd777/noah_api/modules/main/config"
	"github.com/khiemnd777/noah_api/modules/main/department/handler"
	"github.com/khiemnd777/noah_api/modules/main/department/repository"
	"github.com/khiemnd777/noah_api/modules/main/department/service"
	_ "github.com/khiemnd777/noah_api/modules/main/features"
	"github.com/khiemnd777/noah_api/modules/main/registry"
	sharedapp "github.com/khiemnd777/noah_api/shared/app"
	"github.com/khiemnd777/noah_api/shared/db/ent"
	"github.com/khiemnd777/noah_api/shared/db/ent/generated"
	"github.com/khiemnd777/noah_api/shared/metadata/customfields"
	"github.com/khiemnd777/noah_api/shared/middleware"
	"github.com/khiemnd777/noah_api/shared/module"
	"github.com/khiemnd777/noah_api/shared/utils"
	frameworkapp "github.com/khiemnd777/noah_framework/pkg/app"
	frameworkdb "github.com/khiemnd777/noah_framework/pkg/db"
)

func main() {
	module.StartModule(module.ModuleOptions[config.ModuleConfig]{
		ConfigPath: utils.GetModuleConfigPath("main"),
		ModuleName: "main",
		InitEntClient: func(client frameworkdb.Client, cfg *config.ModuleConfig) (any, error) {
			return ent.EntBootstrapFromDatabase(client, func(drv *entsql.Driver) any {
				return generated.NewClient(generated.Driver(drv))
			}, cfg.Database.AutoMigrate)
		},
		OnRegistry: func(app frameworkapp.Application, deps *module.ModuleDeps[config.ModuleConfig]) {
			repo := repository.NewDepartmentRepository(deps.Ent.(*generated.Client), deps)
			router := sharedapp.Group(app, utils.GetModuleRoute(deps.Config.Server.Route), middleware.RequireAuth())

			router.Use("/:dept_id<int>/*",
				middleware.RequireDepartmentMember("dept_id"),
			)

			// Department
			svc := service.NewDepartmentService(repo, deps)
			h := handler.NewDepartmentHandler(svc, deps)
			h.RegisterRoutes(router)

			// Features
			cfStore := &customfields.PGStore{DB: deps.DB}
			cfMgr := customfields.NewManager(cfStore)
			registry.Init(router, deps, cfMgr, registry.InitOptions{
				EnabledIDs: deps.Config.Features.Enabled,
			})
		},
	})
}
