package middleware

import (
	"fmt"
	stdhttp "net/http"

	"github.com/khiemnd777/noah_framework/internal/legacy/shared/app/client_error"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/logger"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

func RecoverAndWrapError() frameworkhttp.Handler {
	return func(c frameworkhttp.Context) error {
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("panic: %v", r)
				logger.ErrorContext(c.UserContext(), "Recovered from panic", "error", err)
				_ = client_error.ResponseError(c, stdhttp.StatusInternalServerError, err)
			}
		}()
		return c.Next()
	}
}
