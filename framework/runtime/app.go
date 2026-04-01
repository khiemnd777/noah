package runtime

import (
	"github.com/gofiber/fiber/v2"
	frameworkfiber "github.com/khiemnd777/noah_framework/internal/http/fiber"
	frameworkapp "github.com/khiemnd777/noah_framework/pkg/app"
)

func NewApplication(cfg frameworkapp.Config) frameworkapp.Application {
	return frameworkfiber.NewApplication(cfg)
}

func AsFiberApp(app frameworkapp.Application) (*fiber.App, bool) {
	fiberApp, ok := app.(interface{ Native() *fiber.App })
	if !ok {
		return nil, false
	}
	return fiberApp.Native(), true
}
