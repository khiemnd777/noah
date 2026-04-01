package app

import (
	"github.com/gofiber/fiber/v2"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"
)

func FiberContext(ctx frameworkhttp.Context) (*fiber.Ctx, error) {
	return frameworkruntime.FiberContext(ctx)
}

func MustFiberContext(ctx frameworkhttp.Context) *fiber.Ctx {
	return frameworkruntime.MustFiberContext(ctx)
}

func Group(app interface{ Router() frameworkhttp.Router }, prefix string, handlers ...frameworkhttp.Handler) frameworkhttp.Router {
	return frameworkruntime.Group(app, prefix, handlers...)
}
