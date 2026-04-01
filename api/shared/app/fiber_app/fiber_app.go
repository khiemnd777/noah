package fiber_app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/khiemnd777/noah_api/shared/app/client_error"
	"github.com/khiemnd777/noah_api/shared/config"
	sharedlogger "github.com/khiemnd777/noah_api/shared/logger"
	"github.com/khiemnd777/noah_api/shared/middleware"
)

type Server struct {
	app *fiber.App
}

func callerLocation(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown:0"
	}

	fn := runtime.FuncForPC(pc)
	funcName := "unknown"
	if fn != nil {
		funcName = fn.Name()
	}

	return fmt.Sprintf("%s:%d (%s)", filepath.Base(file), line, funcName)
}

func responseErrorNative(c *fiber.Ctx, statusCode int, err error, extraMessage ...string) error {
	message := "Server error"
	if statusCode >= fiber.StatusBadRequest && statusCode < fiber.StatusInternalServerError {
		message = "Client error"
	}

	if len(extraMessage) > 0 && extraMessage[0] != "" {
		message = fmt.Sprintf("%s: %s", message, extraMessage[0])
	}

	if os.Getenv("APP_ENV") == "development" && err != nil && (len(extraMessage) == 0 || extraMessage[0] != err.Error()) {
		message = fmt.Sprintf("%s\n%s", message, err.Error())
	}

	location := callerLocation(1)
	logMessage := fmt.Sprintf("%s | at %s", message, location)
	if statusCode >= fiber.StatusInternalServerError {
		sharedlogger.ErrorContext(c.UserContext(), logMessage, "status_code", statusCode, "error", err)
	} else {
		sharedlogger.WarnContext(c.UserContext(), logMessage, "status_code", statusCode, "error", err)
	}

	return c.Status(statusCode).JSON(client_error.ErrorResponse{
		Code:    statusCode,
		Message: message,
	})
}

func Init() (*Server, *fiber.App) {
	srvCfg := config.Get().Server

	app := fiber.New(fiber.Config{
		BodyLimit: srvCfg.BodyLimitMB * 1024 * 1024,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if fe, ok := err.(*fiber.Error); ok {
				return responseErrorNative(c, fe.Code, fe, fe.Error())
			}

			return responseErrorNative(c, fiber.StatusInternalServerError, err, err.Error())
		},
	})

	// Default health check route
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	app.Use(middleware.AttachRequestContext())
	app.Use(logger.New())                     // Simple request logger
	app.Use(middleware.RecoverAndWrapError()) // Recover from panics and return error response

	return &Server{app: app}, app
}

func NewFiberApp() *fiber.App {
	srvCfg := config.Get().Server

	app := fiber.New(fiber.Config{
		BodyLimit: srvCfg.BodyLimitMB * 1024 * 1024,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if fe, ok := err.(*fiber.Error); ok {
				return responseErrorNative(c, fe.Code, fe, fe.Message)
			}

			return responseErrorNative(c, fiber.StatusInternalServerError, err)
		},
	})
	app.Use(middleware.AttachRequestContext())
	app.Use(logger.New())                     // Simple request logger
	app.Use(middleware.RecoverAndWrapError()) // Recover from panics and return error response
	return app
}

func (s *Server) Start() {
	addr := fmt.Sprintf("%s:%d", config.Get().Server.Host, config.Get().Server.Port)
	info := fmt.Sprintf("✅ Fiber listening on http://%s", addr)
	log.Println(info)

	if err := s.app.Listen(addr); err != nil {
		log.Fatal(info)
	}
}
