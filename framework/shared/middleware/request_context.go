package middleware

import (
	"strings"

	"github.com/google/uuid"
	"github.com/khiemnd777/noah_framework/shared/logger"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

func AttachRequestContext() frameworkhttp.Handler {
	return func(c frameworkhttp.Context) error {
		requestID := strings.TrimSpace(c.Get("X-Request-ID"))
		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Set("X-Request-ID", requestID)
		c.Locals("requestID", requestID)

		ctx := logger.ContextWithFields(c.UserContext(), map[string]any{
			"request_id": requestID,
			"method":     c.Method(),
			"path":       c.Path(),
		})
		c.SetUserContext(ctx)

		return c.Next()
	}
}
