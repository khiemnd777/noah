package app

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	frameworkapp "github.com/khiemnd777/noah_framework/pkg/app"
)

func FiberApplication(app frameworkapp.Application) (*fiber.App, error) {
	native := app.Native()
	fiberApp, ok := native.(*fiber.App)
	if !ok {
		return nil, fmt.Errorf("framework application is not fiber-backed")
	}
	return fiberApp, nil
}

func MustFiberApplication(app frameworkapp.Application) *fiber.App {
	fiberApp, err := FiberApplication(app)
	if err != nil {
		panic(err)
	}
	return fiberApp
}

func Group(app frameworkapp.Application, prefix string, handlers ...fiber.Handler) fiber.Router {
	fiberApp := MustFiberApplication(app)
	return fiberApp.Group(prefix, handlers...)
}
