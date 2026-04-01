package app

import (
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

// Get registers GET route with Circuit Breaker + Retry
func RouterGet(router frameworkhttp.Router, path string, handlers ...frameworkhttp.Handler) {
	router.Get(path, WrapHandlers(path, handlers)...)
}

// Post registers POST route with Circuit Breaker + Retry
func RouterPost(router frameworkhttp.Router, path string, handlers ...frameworkhttp.Handler) {
	router.Post(path, WrapHandlers(path, handlers)...)
}

// Post registers POST route with Circuit Breaker + Retry
func RouterPut(router frameworkhttp.Router, path string, handlers ...frameworkhttp.Handler) {
	router.Put(path, WrapHandlers(path, handlers)...)
}

func RouterDelete(router frameworkhttp.Router, path string, handlers ...frameworkhttp.Handler) {
	router.Delete(path, WrapHandlers(path, handlers)...)
}

// Request supports any method (PUT, DELETE, etc) with Circuit Breaker + Retry
func RouterRequest(method string, router frameworkhttp.Router, path string, handlers ...frameworkhttp.Handler) {
	router.Add(method, path, WrapHandlers(path, handlers)...)
}

// GetWithOptions allows passing custom retry options
func RouterGetWithOptions(router frameworkhttp.Router, path string, opts RetryOptions, handlers ...frameworkhttp.Handler) {
	router.Get(path, WrapHandlers(path, handlers, opts)...)
}

// PostWithOptions allows passing custom retry options
func RouterPostWithOptions(router frameworkhttp.Router, path string, opts RetryOptions, handlers ...frameworkhttp.Handler) {
	router.Post(path, WrapHandlers(path, handlers, opts)...)
}

// RequestWithOptions allows passing custom retry options
func RouterRequestWithOptions(method string, router frameworkhttp.Router, path string, opts RetryOptions, handlers ...frameworkhttp.Handler) {
	router.Add(method, path, WrapHandlers(path, handlers, opts)...)
}
