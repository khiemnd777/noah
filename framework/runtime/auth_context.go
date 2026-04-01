package runtime

import (
	"context"
	"errors"
	stdhttp "net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

type contextKey string

const accessTokenContextKey contextKey = "accessToken"

type JWTTokenPayload struct {
	UserID       int
	DepartmentID int
	Email        string
	Roles        *map[string]struct{}
	Permissions  *map[string]struct{}
	Exp          time.Time
}

func AuthSecret() string {
	appCfg, err := LoadYAML[AppConfig](AppConfigPath())
	if err == nil {
		if envSecret := os.Getenv("JWT_TOKEN_SECRET"); envSecret != "" {
			return envSecret
		}
		return appCfg.Auth.Secret
	}
	return os.Getenv("JWT_TOKEN_SECRET")
}

func InternalAuthToken() string {
	appCfg, err := LoadYAML[AppConfig](AppConfigPath())
	if err == nil {
		if envToken := os.Getenv("INTERNAL_AUTH_TOKEN"); envToken != "" {
			return envToken
		}
		return appCfg.Auth.InternalAuthToken
	}
	return os.Getenv("INTERNAL_AUTH_TOKEN")
}

func GetAccessToken(c frameworkhttp.Context) string {
	authHeader := c.Get("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}

func GenerateJWTToken(secret string, payload JWTTokenPayload) (string, error) {
	claims := jwt.MapClaims{
		"user_id": payload.UserID,
		"dept_id": payload.DepartmentID,
		"email":   payload.Email,
		"roles":   payload.Roles,
		"perms":   payload.Permissions,
		"exp":     payload.Exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GetJWTClaims(c frameworkhttp.Context) (jwt.MapClaims, bool, error) {
	secret := AuthSecret()
	header := c.Get("Authorization")
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		return nil, false, respondRuntimeError(c, stdhttp.StatusUnauthorized, "Missing or invalid Authorization header")
	}

	tokenStr := strings.TrimPrefix(header, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, false, respondRuntimeError(c, stdhttp.StatusUnauthorized, "Invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil {
		return nil, false, respondRuntimeError(c, stdhttp.StatusUnauthorized, "Invalid token claims")
	}

	return claims, true, nil
}

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

func RequireInternal(token string) frameworkhttp.Handler {
	return func(c frameworkhttp.Context) error {
		if c.Get("X-Internal-Token") != token {
			return c.Status(stdhttp.StatusUnauthorized).SendString("Unauthorized internal call")
		}
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

func ContextWithAccessToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, accessTokenContextKey, token)
}

func GetPermSetFromClaims(c frameworkhttp.Context) (map[string]struct{}, bool) {
	if v := c.Locals("permSet"); v != nil {
		switch vv := v.(type) {
		case map[string]struct{}:
			if len(vv) > 0 {
				return vv, true
			}
		case *map[string]struct{}:
			if vv != nil && len(*vv) > 0 {
				return *vv, true
			}
		}
	}
	if v := c.Locals("permissions"); v != nil {
		set := make(map[string]struct{})
		switch vv := v.(type) {
		case []string:
			for _, perm := range vv {
				if perm = normalizePermission(perm); perm != "" {
					set[perm] = struct{}{}
				}
			}
		case []any:
			for _, perm := range vv {
				if s, ok := perm.(string); ok {
					if s = normalizePermission(s); s != "" {
						set[s] = struct{}{}
					}
				}
			}
		}
		if len(set) > 0 {
			c.Locals("permSet", set)
			return set, true
		}
	}

	claims, ok, _ := GetJWTClaims(c)
	if !ok || claims == nil {
		return nil, false
	}
	raw, exists := claims["perms"]
	if !exists || raw == nil {
		return nil, false
	}

	set := make(map[string]struct{})
	switch v := raw.(type) {
	case *map[string]struct{}:
		if v != nil {
			for k := range *v {
				if k = normalizePermission(k); k != "" {
					set[k] = struct{}{}
				}
			}
		}
	case map[string]struct{}:
		for k := range v {
			if k = normalizePermission(k); k != "" {
				set[k] = struct{}{}
			}
		}
	case map[string]bool:
		for k, val := range v {
			if val {
				if k = normalizePermission(k); k != "" {
					set[k] = struct{}{}
				}
			}
		}
	case map[string]any:
		for k := range v {
			if k = normalizePermission(k); k != "" {
				set[k] = struct{}{}
			}
		}
	case []string:
		for _, k := range v {
			if k = normalizePermission(k); k != "" {
				set[k] = struct{}{}
			}
		}
	case []any:
		for _, it := range v {
			if s, ok := it.(string); ok {
				if s = normalizePermission(s); s != "" {
					set[s] = struct{}{}
				}
			}
		}
	case string:
		if strings.Contains(v, ",") {
			for _, part := range strings.Split(v, ",") {
				if part = normalizePermission(part); part != "" {
					set[part] = struct{}{}
				}
			}
		} else {
			for _, part := range strings.Fields(v) {
				if part = normalizePermission(part); part != "" {
					set[part] = struct{}{}
				}
			}
		}
	default:
		return nil, false
	}

	if len(set) == 0 {
		return nil, false
	}
	c.Locals("permSet", set)
	return set, true
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

func normalizePermission(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
