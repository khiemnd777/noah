package app

import (
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"
)

type RetryOptions = frameworkhttp.RetryOptions

func WrapHandler(name string, h frameworkhttp.Handler, opts ...RetryOptions) frameworkhttp.Handler {
	return frameworkruntime.WrapHandler(name, h, opts...)
}

func WrapHandlers(name string, handlers []frameworkhttp.Handler, opts ...RetryOptions) []frameworkhttp.Handler {
	return frameworkruntime.WrapHandlers(name, handlers, opts...)
}
