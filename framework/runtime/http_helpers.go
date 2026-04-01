package runtime

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sony/gobreaker"

	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

func ParseBody[T any](c frameworkhttp.Context) (*T, error) {
	var body T
	if err := c.BodyParser(&body); err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}
	return &body, nil
}

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

func RouterGet(router frameworkhttp.Router, path string, handlers ...frameworkhttp.Handler) {
	router.Get(path, WrapHandlers(path, handlers)...)
}

func RouterPost(router frameworkhttp.Router, path string, handlers ...frameworkhttp.Handler) {
	router.Post(path, WrapHandlers(path, handlers)...)
}

func RouterPut(router frameworkhttp.Router, path string, handlers ...frameworkhttp.Handler) {
	router.Put(path, WrapHandlers(path, handlers)...)
}

func RouterDelete(router frameworkhttp.Router, path string, handlers ...frameworkhttp.Handler) {
	router.Delete(path, WrapHandlers(path, handlers)...)
}

func RouterRequest(method string, router frameworkhttp.Router, path string, handlers ...frameworkhttp.Handler) {
	router.Add(method, path, WrapHandlers(path, handlers)...)
}

func RouterGetWithOptions(router frameworkhttp.Router, path string, opts frameworkhttp.RetryOptions, handlers ...frameworkhttp.Handler) {
	router.Get(path, WrapHandlers(path, handlers, opts)...)
}

func RouterPostWithOptions(router frameworkhttp.Router, path string, opts frameworkhttp.RetryOptions, handlers ...frameworkhttp.Handler) {
	router.Post(path, WrapHandlers(path, handlers, opts)...)
}

func RouterRequestWithOptions(method string, router frameworkhttp.Router, path string, opts frameworkhttp.RetryOptions, handlers ...frameworkhttp.Handler) {
	router.Add(method, path, WrapHandlers(path, handlers, opts)...)
}

func WrapHandler(name string, h frameworkhttp.Handler, opts ...frameworkhttp.RetryOptions) frameworkhttp.Handler {
	retry := defaultRetryOptions()
	if len(opts) > 0 {
		retry = mergeRetryOptions(opts[0], retry)
	}

	return func(c frameworkhttp.Context) error {
		fiberCtx := MustFiberContext(c)

		if isWebSocketRequest(c) {
			return h(c)
		}

		var err error
		for i := 0; i < retry.MaxAttempts; i++ {
			_, err = RunWithCircuitBreaker(name, func(ctx context.Context) (any, error) {
				handleErr := h(c)

				if ferr, ok := handleErr.(*fiber.Error); ok && ferr.Code >= http.StatusBadRequest && ferr.Code < http.StatusInternalServerError {
					return nil, ErrClientResponse
				}

				if handleErr == nil {
					statusCode := fiberCtx.Response().StatusCode()
					if statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError {
						return nil, ErrClientResponse
					}
				}

				return nil, handleErr
			})

			if errors.Is(err, ErrClientResponse) {
				return nil
			}

			if err == nil || !retry.ShouldRetry(err) {
				return err
			}

			time.Sleep(retry.Delay)
		}

		return err
	}
}

func WrapHandlers(name string, handlers []frameworkhttp.Handler, opts ...frameworkhttp.RetryOptions) []frameworkhttp.Handler {
	if len(handlers) == 0 {
		panic("no handlers provided")
	}
	last := handlers[len(handlers)-1]
	middlewares := handlers[:len(handlers)-1]
	wrapped := WrapHandler(name, last, opts...)
	return append(middlewares, wrapped)
}

func isWebSocketRequest(c frameworkhttp.Context) bool {
	if c.Method() != fiber.MethodGet {
		return false
	}
	if strings.ToLower(c.Get("Upgrade")) != "websocket" {
		return false
	}
	if !strings.Contains(strings.ToLower(c.Get("Connection")), "upgrade") {
		return false
	}
	return true
}

func defaultRetryOptions() frameworkhttp.RetryOptions {
	appCfg, err := LoadYAML[AppConfig](APIPath("config.yaml"))
	if err != nil {
		return frameworkhttp.RetryOptions{
			MaxAttempts: 3,
			Delay:       200 * time.Millisecond,
			ShouldRetry: defaultShouldRetry,
		}
	}

	maxAttempts := appCfg.Retry.MaxAttempts
	if maxAttempts <= 0 {
		maxAttempts = 3
	}
	delay := appCfg.Retry.Delay
	if delay <= 0 {
		delay = 200 * time.Millisecond
	}

	return frameworkhttp.RetryOptions{
		MaxAttempts: maxAttempts,
		Delay:       delay,
		ShouldRetry: defaultShouldRetry,
	}
}

func mergeRetryOptions(custom frameworkhttp.RetryOptions, defaults frameworkhttp.RetryOptions) frameworkhttp.RetryOptions {
	if custom.MaxAttempts <= 0 {
		custom.MaxAttempts = defaults.MaxAttempts
	}
	if custom.Delay <= 0 {
		custom.Delay = defaults.Delay
	}
	if custom.ShouldRetry == nil {
		custom.ShouldRetry = defaults.ShouldRetry
	}
	return custom
}

func defaultShouldRetry(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrClientResponse) || errors.Is(err, gobreaker.ErrOpenState) {
		return false
	}
	if ferr, ok := err.(*fiber.Error); ok && ferr.Code >= http.StatusBadRequest && ferr.Code < http.StatusInternalServerError {
		return false
	}

	type statusError interface {
		HTTPStatusCode() int
	}
	var statusErr statusError
	if errors.As(err, &statusErr) && statusErr.HTTPStatusCode() < http.StatusInternalServerError {
		return false
	}

	return true
}
