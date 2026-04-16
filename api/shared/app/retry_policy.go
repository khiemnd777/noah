package app

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/khiemnd777/noah_api/shared/circuitbreaker"
	"github.com/khiemnd777/noah_api/shared/config"
	"github.com/sony/gobreaker"
	"github.com/valyala/fasthttp"
)

func defaultRetryOptionsForMethod(method string) RetryOptions {
	cfgRetry := config.Get().Retry
	maxAttempts := 1
	if isSafeRetryMethod(method) {
		maxAttempts = cfgRetry.MaxAttempts
	}

	return RetryOptions{
		MaxAttempts: maxAttempts,
		Delay:       cfgRetry.Delay,
		ShouldRetry: defaultShouldRetry,
	}
}

func mergeRetryOptions(base RetryOptions, opts []RetryOptions) RetryOptions {
	if len(opts) == 0 {
		return ensureRetryOptions(base, base)
	}

	return ensureRetryOptions(opts[0], base)
}

func ensureRetryOptions(current, fallback RetryOptions) RetryOptions {
	if current.ShouldRetry == nil {
		current.ShouldRetry = fallback.ShouldRetry
	}
	if current.MaxAttempts <= 0 {
		current.MaxAttempts = fallback.MaxAttempts
	}
	if current.Delay == 0 {
		current.Delay = fallback.Delay
	}
	if current.MaxAttempts <= 0 {
		current.MaxAttempts = 1
	}
	return current
}

func defaultShouldRetry(err error) bool {
	if err == nil {
		return false
	}
	if isRecoveredPanic(err) {
		return false
	}
	if errors.Is(err, circuitbreaker.ErrClientResponse) || errors.Is(err, gobreaker.ErrOpenState) {
		return false
	}
	if ferr, ok := err.(*fiber.Error); ok && ferr.Code >= http.StatusBadRequest && ferr.Code < http.StatusInternalServerError {
		return false
	}
	var statusErr *StatusError
	if errors.As(err, &statusErr) && statusErr.StatusCode < http.StatusInternalServerError {
		return false
	}
	return true
}

func isSafeRetryMethod(method string) bool {
	switch strings.ToUpper(method) {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	default:
		return false
	}
}

func isRecoveredPanic(err error) bool {
	if err == nil {
		return false
	}
	return strings.HasPrefix(err.Error(), "panic:")
}

func responseWasWritten(resp *fasthttp.Response) bool {
	if resp == nil {
		return false
	}
	if resp.StatusCode() != fiber.StatusOK {
		return true
	}
	return len(resp.Body()) > 0
}
