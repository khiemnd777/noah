package utils

import (
	"github.com/golang-jwt/jwt/v5"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"
)

type JWTTokenPayload = frameworkruntime.JWTTokenPayload

func GenerateJWTToken(secret string, payload JWTTokenPayload) (string, error) {
	return frameworkruntime.GenerateJWTToken(secret, payload)
}

func GetJWTClaims(c frameworkhttp.Context) (jwt.MapClaims, bool, error) {
	return frameworkruntime.GetJWTClaims(c)
}

func GetPermSetFromClaims(c frameworkhttp.Context) (map[string]struct{}, bool) {
	return frameworkruntime.GetPermSetFromClaims(c)
}
