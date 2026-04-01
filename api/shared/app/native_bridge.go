package app

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

func FiberContext(ctx frameworkhttp.Context) (*fiber.Ctx, error) {
	native := ctx.Native()
	fiberCtx, ok := native.(*fiber.Ctx)
	if !ok {
		return nil, fmt.Errorf("framework context is not fiber-backed")
	}
	return fiberCtx, nil
}

func MustFiberContext(ctx frameworkhttp.Context) *fiber.Ctx {
	fiberCtx, err := FiberContext(ctx)
	if err != nil {
		panic(err)
	}
	return fiberCtx
}

func Group(app interface{ Router() frameworkhttp.Router }, prefix string, handlers ...frameworkhttp.Handler) frameworkhttp.Router {
	return app.Router().Mount(prefix, handlers...)
}
