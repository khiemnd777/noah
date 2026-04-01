package material

import (
	"github.com/khiemnd777/noah_api/modules/main/config"
	"github.com/khiemnd777/noah_api/modules/main/features/__relation/handler"
	_ "github.com/khiemnd777/noah_api/modules/main/features/__relation/registrar"
	"github.com/khiemnd777/noah_api/modules/main/features/__relation/repository"
	"github.com/khiemnd777/noah_api/modules/main/features/__relation/service"
	"github.com/khiemnd777/noah_api/modules/main/registry"
	"github.com/khiemnd777/noah_framework/shared/metadata/customfields"
	"github.com/khiemnd777/noah_framework/shared/module"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

type feature struct{}

func (feature) ID() string    { return "relation" }
func (feature) Priority() int { return 1 }

func (feature) Register(router frameworkhttp.Router, deps *module.ModuleDeps[config.ModuleConfig], cfMgr *customfields.Manager) error {
	repo := repository.NewRelationRepository()
	svc := service.NewRelationService(repo, deps)
	h := handler.NewRelationHandler(svc, deps)
	h.RegisterRoutes(router)
	return nil
}

func init() { registry.Register(feature{}) }
