package middleware

import (
	stdhttp "net/http"
	"strings"

	"github.com/khiemnd777/noah_api/shared/app/client_error"
	"github.com/khiemnd777/noah_api/shared/logger"
	"github.com/khiemnd777/noah_api/shared/utils"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

func RequireAuth() frameworkhttp.Handler {
	return func(c frameworkhttp.Context) error {
		header := c.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			return client_error.ResponseError(c, stdhttp.StatusUnauthorized, nil, "Missing or invalid Authorization header")
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")

		claims, ok, err := utils.GetJWTClaims(c)

		if !ok || err != nil {
			return client_error.ResponseError(c, stdhttp.StatusUnauthorized, err, "Invalid token claims")
		}

		c.Locals("userID", int(claims["user_id"].(float64)))
		c.Locals("deptID", int(claims["dept_id"].(float64)))

		// Inject token into context for downstream access
		ctxWithToken := utils.SetAccessTokenIntoContext(c.UserContext(), tokenStr)
		ctxWithToken = logger.ContextWithFields(ctxWithToken, map[string]any{
			"user_id":       int(claims["user_id"].(float64)),
			"department_id": int(claims["dept_id"].(float64)),
		})
		c.SetUserContext(ctxWithToken)

		return c.Next()
	}
}

func RequireInternal() frameworkhttp.Handler {
	// Only use for internal audit, trace, impersonate, và trust-based routing.
	return func(c frameworkhttp.Context) error {
		token := c.Get("X-Internal-Token")
		baseIntrTkn := utils.GetInternalToken()
		if token != baseIntrTkn {
			return c.Status(401).SendString("Unauthorized internal call")
		}
		return c.Next()
	}
}
