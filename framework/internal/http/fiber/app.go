package fiber

import (
	"fmt"
	"net"
	"strconv"

	"github.com/gofiber/fiber/v2"
	frameworkapp "github.com/khiemnd777/noah_framework/pkg/app"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

type Application struct {
	app    *fiber.App
	router frameworkhttp.Router
	addr   string
}

func NewApplication(cfg frameworkapp.Config) *Application {
	app := fiber.New(fiber.Config{
		BodyLimit: cfg.BodyLimitMB * 1024 * 1024,
	})

	wrapped := WrapApplication(app)
	wrapped.addr = cfg.Host + ":" + strconv.Itoa(cfg.Port)
	return wrapped
}

func WrapApplication(app *fiber.App) *Application {
	return &Application{
		app:    app,
		router: NewRouter(app),
	}
}

func (a *Application) Router() frameworkhttp.Router {
	return a.router
}

func (a *Application) Native() any {
	return a.app
}

func (a *Application) Listen(addr string) error {
	return a.app.Listen(addr)
}

func (a *Application) Serve(listener any) error {
	native, ok := listener.(net.Listener)
	if !ok {
		return fmt.Errorf("unsupported listener type %T", listener)
	}
	return a.app.Listener(native)
}

func (a *Application) Run() error {
	return a.app.Listen(a.addr)
}
