package runtime

import (
	"github.com/gofiber/fiber/v2"
	frameworkfiber "github.com/khiemnd777/noah_framework/internal/http/fiber"
	frameworkapp "github.com/khiemnd777/noah_framework/pkg/app"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

func NewApplication(cfg frameworkapp.Config) frameworkapp.Application {
	return frameworkfiber.NewApplication(cfg)
}

func WrapFiberApplication(app *fiber.App) frameworkapp.Application {
	return frameworkfiber.WrapApplication(app)
}

func WrapFiberRouter(router fiber.Router) frameworkhttp.Router {
	return frameworkfiber.NewRouter(router)
}

func AsFiberApp(app frameworkapp.Application) (*fiber.App, bool) {
	native := app.Native()
	fiberApp, ok := native.(*fiber.App)
	return fiberApp, ok
}
