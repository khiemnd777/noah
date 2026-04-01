package utils

import (
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"
)

func GetAuthSecret() string {
	return frameworkruntime.AuthSecret()
}

func GetInternalToken() string {
	return frameworkruntime.InternalAuthToken()
}

func GetAccessToken(c frameworkhttp.Context) string {
	return frameworkruntime.GetAccessToken(c)
}
