package department

import (
	"github.com/khiemnd777/noah_api/modules/main/config"
	"github.com/khiemnd777/noah_api/modules/main/department/handler"
	"github.com/khiemnd777/noah_api/modules/main/department/repository"
	"github.com/khiemnd777/noah_api/modules/main/department/service"
	"github.com/khiemnd777/noah_framework/shared/db/ent/generated"
	"github.com/khiemnd777/noah_framework/shared/module"
	frameworkapp "github.com/khiemnd777/noah_framework/pkg/app"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

func NewDepartmentModule(app frameworkapp.Application, router frameworkhttp.Router, deps *module.ModuleDeps[config.ModuleConfig]) {
	repo := repository.NewDepartmentRepository(deps.Ent.(*generated.Client), deps)
	svc := service.NewDepartmentService(repo, deps)
	h := handler.NewDepartmentHandler(svc, deps)
	h.RegisterRoutes(router)
}
