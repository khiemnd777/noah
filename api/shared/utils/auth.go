package utils

import (
	"os"

	"github.com/khiemnd777/noah_api/shared/config"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

func GetAuthSecret() string {
	cfg, _ := LoadConfig[config.AppConfig](GetAppConfigPath())

	if envSecret := os.Getenv("JWT_TOKEN_SECRET"); envSecret != "" {
		return envSecret
	}
	return cfg.Auth.Secret
}

func GetInternalToken() string {
	cfg, _ := LoadConfig[config.AppConfig](GetAppConfigPath())

	if envSecret := os.Getenv("INTERNAL_AUTH_TOKEN"); envSecret != "" {
		return envSecret
	}
	return cfg.Auth.InternalAuthToken
}

// GetAccessToken extracts the Bearer token from the Authorization header
func GetAccessToken(c frameworkhttp.Context) string {
	authHeader := c.Get("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}
