package utils

import (
	"context"

	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"
)

type ctxKey string

const CtxAccessTokenKey ctxKey = "accessToken"

// SetAccessTokenIntoContext returns a new context with the access token
func SetAccessTokenIntoContext(ctx context.Context, token string) context.Context {
	return frameworkruntime.ContextWithAccessToken(ctx, token)
}

func GetAccessTokenFromContext(ctx context.Context) string {
	return frameworkruntime.AccessTokenFromContext(ctx)
}
