package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/khiemnd777/noah_api/shared/circuitbreaker"
	"github.com/khiemnd777/noah_api/shared/logger"
)

// RetryOptions configures retry behavior for a route
type RetryOptions struct {
	MaxAttempts int
	Delay       time.Duration
	ShouldRetry func(error) bool
}

func isWebSocketRequest(c *fiber.Ctx) bool {
	// RFC 6455
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

// WrapHandler applies Circuit Breaker + Retry logic to a single handler
func WrapHandler(name string, h fiber.Handler, opts ...RetryOptions) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if isWebSocketRequest(c) {
			logger.Info("🔌 WS bypass circuit: " + name)
			return h(c)
		}

		method := c.Method()
		retry := mergeRetryOptions(defaultRetryOptionsForMethod(method), opts)

		var err error
		for i := 0; i < retry.MaxAttempts; i++ {
			_, err = circuitbreaker.Run(name, func(ctx context.Context) (interface{}, error) {
				handleErr := h(c)

				if ferr, ok := handleErr.(*fiber.Error); ok && ferr.Code >= 400 && ferr.Code < 500 {
					return nil, circuitbreaker.ErrClientResponse
				}

				if handleErr == nil {
					statusCode := c.Response().StatusCode()
					if statusCode >= fiber.StatusBadRequest && statusCode < fiber.StatusInternalServerError {
						return nil, circuitbreaker.ErrClientResponse
					}
				}

				return nil, handleErr
			})

			if errors.Is(err, circuitbreaker.ErrClientResponse) {
				return nil
			}

			if isRecoveredPanic(err) {
				return err
			}

			if err == nil || !retry.ShouldRetry(err) {
				return err
			}

			if i == retry.MaxAttempts-1 {
				break
			}

			if !isSafeRetryMethod(method) && responseWasWritten(c.Response()) {
				return err
			}

			if isSafeRetryMethod(method) {
				c.Response().Reset()
			}

			logger.Warn(fmt.Sprintf("🔁 Retry [%s %s] #%d failed: %v", method, name, i+1, err))
			time.Sleep(retry.Delay)
		}

		log.Printf("❌ Handler failed after retries: %s %s", method, name)
		return err
	}
}

// WrapHandlers applies middleware(s) and wraps the final handler with CB + Retry
func WrapHandlers(name string, handlers []fiber.Handler, opts ...RetryOptions) []fiber.Handler {
	if len(handlers) == 0 {
		panic("no handlers provided")
	}
	last := handlers[len(handlers)-1]
	middlewares := handlers[:len(handlers)-1]
	wrapped := WrapHandler(name, last, opts...)
	return append(middlewares, wrapped)
}
