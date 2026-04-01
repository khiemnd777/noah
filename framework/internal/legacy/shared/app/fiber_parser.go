package app

import (
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"
)

func ParseBody[T any](c frameworkhttp.Context) (*T, error) {
	return frameworkruntime.ParseBody[T](c)
}
