package main

import (
	entsql "entgo.io/ent/dialect/sql"
	"github.com/khiemnd777/noah_api/modules/observability/config"
	"github.com/khiemnd777/noah_api/shared/app"
	"github.com/khiemnd777/noah_api/shared/db/ent"
	"github.com/khiemnd777/noah_api/shared/db/ent/generated"
	"github.com/khiemnd777/noah_api/shared/middleware"
	"github.com/khiemnd777/noah_api/shared/middleware/rbac"
	"github.com/khiemnd777/noah_api/shared/module"
	sharedruntime "github.com/khiemnd777/noah_api/shared/runtime"
	"github.com/khiemnd777/noah_api/shared/utils"
	frameworkobservability "github.com/khiemnd777/noah_framework/modules/observability/module"
	frameworkapp "github.com/khiemnd777/noah_framework/pkg/app"
	frameworkdb "github.com/khiemnd777/noah_framework/pkg/db"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

func main() {
	configPath, err := sharedruntime.ModuleConfigPath("observability")
	if err != nil {
		configPath = utils.GetModuleConfigPath("observability")
	}

	module.StartModule(module.ModuleOptions[config.ModuleConfig]{
		ConfigPath: configPath,
		ModuleName: "observability",
		InitEntClient: func(client frameworkdb.Client, cfg *config.ModuleConfig) (any, error) {
			return ent.EntBootstrapFromDatabase(client, func(drv *entsql.Driver) any {
				return generated.NewClient(generated.Driver(drv))
			}, cfg.Database.AutoMigrate)
		},
		OnRegistry: func(moduleApp frameworkapp.Application, deps *module.ModuleDeps[config.ModuleConfig]) {
			frameworkobservability.Register(
				app.Group(moduleApp, utils.GetModuleRoute(deps.Config.Server.Route), middleware.RequireAuth()),
				frameworkobservability.Options{
					Config: config.FrameworkModuleConfig(),
					RequirePermission: func(c frameworkhttp.Context, permission string) error {
						return rbac.GuardAnyPermission(c, deps.Ent.(*generated.Client), permission)
					},
				},
			)
		},
	})
}
