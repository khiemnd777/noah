package main

import (
	entsql "entgo.io/ent/dialect/sql"
	"github.com/khiemnd777/noah_framework/modules/metadata/config"
	"github.com/khiemnd777/noah_framework/modules/metadata/handler"
	"github.com/khiemnd777/noah_framework/modules/metadata/repository"
	"github.com/khiemnd777/noah_framework/modules/metadata/service"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/app"
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
		ConfigPath: utils.GetModuleConfigPath("metadata"),
		ModuleName: "metadata",
		InitEntClient: func(client frameworkdb.Client, cfg *config.ModuleConfig) (any, error) {
			return ent.EntBootstrapFromDatabase(client, func(drv *entsql.Driver) any {
				return generated.NewClient(generated.Driver(drv))
			}, cfg.Database.AutoMigrate)
		},
		OnRegistry: func(moduleApp frameworkapp.Application, deps *module.ModuleDeps[config.ModuleConfig]) {
			db := deps.DB
			baseRoute := utils.GetModuleRoute(deps.Config.Server.Route)

			cltRepo := repository.NewCollectionRepository(db)
			cltSvc := service.NewCollectionService(cltRepo)
			cltH := handler.NewCollectionHandler(cltSvc, deps)
			cltH.RegisterRoutes(app.Group(moduleApp, baseRoute, middleware.RequireAuth()))

			fRepo := repository.NewFieldRepository(db)
			fSvc := service.NewFieldService(fRepo, cltRepo)
			fH := handler.NewFieldHandler(fSvc, deps)
			fH.RegisterRoutes(app.Group(moduleApp, baseRoute, middleware.RequireAuth()))

			ipRepo := repository.NewImportFieldProfileRepository(db)
			ipSvc := service.NewImportFieldProfileService(ipRepo)
			ipH := handler.NewImportFieldProfileHandler(ipSvc, deps)
			ipH.RegisterRoutes(app.Group(moduleApp, baseRoute, middleware.RequireAuth()))

			imRepo := repository.NewImportFieldMappingRepository(db)
			imSvc := service.NewImportFieldMappingService(imRepo, ipRepo)
			imH := handler.NewImportFieldMappingHandler(imSvc, deps)
			imH.RegisterRoutes(app.Group(moduleApp, baseRoute, middleware.RequireAuth()))

			iEn := service.NewImportEngine(db, imSvc)
			iH := handler.NewImportHandler(iEn, deps)
			iH.RegisterRoutes(app.Group(moduleApp, baseRoute, middleware.RequireAuth()))

			eH := handler.NewExportHandler(iEn, deps)
			eH.RegisterRoutes(app.Group(moduleApp, baseRoute, middleware.RequireAuth()))
		},
	})
}
