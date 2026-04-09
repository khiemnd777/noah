package main

import (
	"database/sql"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/khiemnd777/noah_api/modules/i18n/config"
	"github.com/khiemnd777/noah_api/modules/i18n/handler"
	"github.com/khiemnd777/noah_api/modules/i18n/repository"
	"github.com/khiemnd777/noah_api/modules/i18n/service"
	"github.com/khiemnd777/noah_api/shared/db/ent"
	"github.com/khiemnd777/noah_api/shared/db/ent/generated"
	"github.com/khiemnd777/noah_api/shared/middleware"
	"github.com/khiemnd777/noah_api/shared/module"
	"github.com/khiemnd777/noah_api/shared/utils"
)

func main() {
	module.StartModule(module.ModuleOptions[config.ModuleConfig]{
		ConfigPath: utils.GetModuleConfigPath("i18n"),
		ModuleName: "i18n",
		InitEntClient: func(provider string, db *sql.DB, cfg *config.ModuleConfig) (any, error) {
			return ent.EntBootstrap(provider, db, func(drv *entsql.Driver) any {
				return generated.NewClient(generated.Driver(drv))
			}, cfg.Database.AutoMigrate)
		},
		OnRegistry: func(app *fiber.App, deps *module.ModuleDeps[config.ModuleConfig]) {
			db := deps.Ent.(*generated.Client)
			repo := repository.NewLanguageRepository(db)
			svc := service.NewLanguageService(repo)
			h := handler.NewLanguageHandler(svc, deps)
			h.RegisterRoutes(app.Group(utils.GetModuleRoute(deps.Config.Server.Route), middleware.RequireAuth()))
		},
	})
}
