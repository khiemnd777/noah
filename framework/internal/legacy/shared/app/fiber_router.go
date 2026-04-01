package app

import (
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"
)

func RouterGet(router frameworkhttp.Router, path string, handlers ...frameworkhttp.Handler) {
	frameworkruntime.RouterGet(router, path, handlers...)
}

func RouterPost(router frameworkhttp.Router, path string, handlers ...frameworkhttp.Handler) {
	frameworkruntime.RouterPost(router, path, handlers...)
}

func RouterPut(router frameworkhttp.Router, path string, handlers ...frameworkhttp.Handler) {
	frameworkruntime.RouterPut(router, path, handlers...)
}

func RouterDelete(router frameworkhttp.Router, path string, handlers ...frameworkhttp.Handler) {
	frameworkruntime.RouterDelete(router, path, handlers...)
}

func RouterRequest(method string, router frameworkhttp.Router, path string, handlers ...frameworkhttp.Handler) {
	frameworkruntime.RouterRequest(method, router, path, handlers...)
}

func RouterGetWithOptions(router frameworkhttp.Router, path string, opts RetryOptions, handlers ...frameworkhttp.Handler) {
	frameworkruntime.RouterGetWithOptions(router, path, frameworkhttp.RetryOptions(opts), handlers...)
}

func RouterPostWithOptions(router frameworkhttp.Router, path string, opts RetryOptions, handlers ...frameworkhttp.Handler) {
	frameworkruntime.RouterPostWithOptions(router, path, frameworkhttp.RetryOptions(opts), handlers...)
}

func RouterRequestWithOptions(method string, router frameworkhttp.Router, path string, opts RetryOptions, handlers ...frameworkhttp.Handler) {
	frameworkruntime.RouterRequestWithOptions(method, router, path, frameworkhttp.RetryOptions(opts), handlers...)
}
