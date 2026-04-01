package module

import (
	"github.com/khiemnd777/noah_framework/modules/observability/handler"
	"github.com/khiemnd777/noah_framework/modules/observability/repository"
	"github.com/khiemnd777/noah_framework/modules/observability/service"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

type Config struct {
	Loki repository.LokiConfig
}

type Options struct {
	Config            Config
	RequirePermission handler.PermissionGuard
}

func Register(router frameworkhttp.Router, opts Options) {
	repo := repository.New(opts.Config.Loki)
	svc := service.New(repo)
	h := handler.New(svc, opts.RequirePermission)
	h.RegisterRoutes(router)
}
