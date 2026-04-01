package main

import (
	entsql "entgo.io/ent/dialect/sql"

	sharedapp "github.com/khiemnd777/noah_api/shared/app"
	"github.com/khiemnd777/noah_api/shared/db/ent"
	"github.com/khiemnd777/noah_api/shared/middleware"
	"github.com/khiemnd777/noah_api/shared/utils"

	"github.com/khiemnd777/noah_api/modules/metadata/config"
	"github.com/khiemnd777/noah_api/modules/metadata/handler"
	"github.com/khiemnd777/noah_api/modules/metadata/repository"
	"github.com/khiemnd777/noah_api/modules/metadata/service"
	"github.com/khiemnd777/noah_api/shared/db/ent/generated"
	"github.com/khiemnd777/noah_api/shared/module"
	frameworkapp "github.com/khiemnd777/noah_framework/pkg/app"
	frameworkdb "github.com/khiemnd777/noah_framework/pkg/db"
)

func main() {
	module.StartModule(module.ModuleOptions[config.ModuleConfig]{
		ConfigPath: utils.GetModuleConfigPath("metadata"),
		ModuleName: "metadata",
		InitEntClient: func(client frameworkdb.Client, cfg *config.ModuleConfig) (any, error) {
			return ent.EntBootstrapFromDatabase(client, func(drv *entsql.Driver) any {
				return generated.NewClient(generated.Driver(drv))
			}, cfg.Database.AutoMigrate)
		},
		OnRegistry: func(app frameworkapp.Application, deps *module.ModuleDeps[config.ModuleConfig]) {
			db := deps.DB

			// Collection
			cltRepo := repository.NewCollectionRepository(db)
			cltSvc := service.NewCollectionService(cltRepo)
			cltH := handler.NewCollectionHandler(cltSvc, deps)
			cltH.RegisterRoutes(sharedapp.Group(app, utils.GetModuleRoute(deps.Config.Server.Route), middleware.RequireAuth()))

			// Field
			fRepo := repository.NewFieldRepository(db)
			fSvc := service.NewFieldService(fRepo, cltRepo)
			fH := handler.NewFieldHandler(fSvc, deps)
			fH.RegisterRoutes(sharedapp.Group(app, utils.GetModuleRoute(deps.Config.Server.Route), middleware.RequireAuth()))

			// Import Field Profiles
			ipRepo := repository.NewImportFieldProfileRepository(db)
			ipSvc := service.NewImportFieldProfileService(ipRepo)
			ipH := handler.NewImportFieldProfileHandler(ipSvc, deps)
			ipH.RegisterRoutes(sharedapp.Group(app, utils.GetModuleRoute(deps.Config.Server.Route), middleware.RequireAuth()))

			// Import Field Mappings
			imRepo := repository.NewImportFieldMappingRepository(db)
			imSvc := service.NewImportFieldMappingService(imRepo, ipRepo)
			imH := handler.NewImportFieldMappingHandler(imSvc, deps)
			imH.RegisterRoutes(sharedapp.Group(app, utils.GetModuleRoute(deps.Config.Server.Route), middleware.RequireAuth()))

			// Import
			iEn := service.NewImportEngine(db, imSvc)
			iH := handler.NewImportHandler(iEn, deps)
			iH.RegisterRoutes(sharedapp.Group(app, utils.GetModuleRoute(deps.Config.Server.Route), middleware.RequireAuth()))

			//Export
			eH := handler.NewExportHandler(iEn, deps)
			eH.RegisterRoutes(sharedapp.Group(app, utils.GetModuleRoute(deps.Config.Server.Route), middleware.RequireAuth()))
		},
	})
}
