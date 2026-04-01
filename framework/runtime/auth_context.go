package runtime

import (
	"context"
	"errors"
	stdhttp "net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

type contextKey string

const accessTokenContextKey contextKey = "accessToken"

func RequireAuth(secret string) frameworkhttp.Handler {
	return func(c frameworkhttp.Context) error {
		header := c.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			return respondRuntimeError(c, stdhttp.StatusUnauthorized, "Missing or invalid Authorization header")
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return respondRuntimeError(c, stdhttp.StatusUnauthorized, "Invalid or expired token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return respondRuntimeError(c, stdhttp.StatusUnauthorized, "Invalid token claims")
		}

		if userID, ok := claimInt(claims, "user_id"); ok {
			c.Locals("userID", userID)
		}
		if deptID, ok := claimInt(claims, "dept_id"); ok {
			c.Locals("deptID", deptID)
		}
		c.SetUserContext(context.WithValue(c.UserContext(), accessTokenContextKey, tokenStr))
		return c.Next()
	}
}

func UserIDFromContext(c frameworkhttp.Context) (int, bool) {
	return localInt(c, "userID")
}

func DeptIDFromContext(c frameworkhttp.Context) (int, bool) {
	return localInt(c, "deptID")
}

func AccessTokenFromContext(ctx context.Context) string {
	if value, ok := ctx.Value(accessTokenContextKey).(string); ok {
		return value
	}
	return ""
}

func claimInt(claims jwt.MapClaims, key string) (int, bool) {
	raw, ok := claims[key]
	if !ok {
		return 0, false
	}

	switch value := raw.(type) {
	case float64:
		return int(value), true
	case int:
		return value, true
	case int64:
		return int(value), true
	case string:
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}

func localInt(c frameworkhttp.Context, key string) (int, bool) {
	switch value := c.Locals(key).(type) {
	case int:
		return value, true
	case int64:
		return int(value), true
	case float64:
		return int(value), true
	case string:
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}

func respondRuntimeError(c frameworkhttp.Context, statusCode int, message string) error {
	return c.Status(statusCode).JSON(map[string]any{
		"statusCode":    statusCode,
		"statusMessage": message,
	})
}
