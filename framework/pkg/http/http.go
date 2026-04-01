package http

import (
	"context"
	"mime/multipart"
	"time"
)

type Handler func(Context) error

type Context interface {
	BodyParser(out any) error
	Body() []byte
	FormFile(key string) (*multipart.FileHeader, error)
	FormValue(key string, defaultValue ...string) string
	Get(key string, defaultValue ...string) string
	Locals(key any, value ...any) any
	Method(override ...string) string
	Next() error
	OriginalURL() string
	Path() string
	Param(key string) string
	Params(key string, defaultValue ...string) string
	ParamsInt(key string) (int, error)
	Query(key string, defaultValue ...string) string
	QueryBool(key string, defaultValue ...bool) bool
	QueryInt(key string, defaultValue ...int) int
	QueryString() string
	SendFile(path string, compress ...bool) error
	SendStatus(status int) error
	SendString(body string) error
	Status(status int) Context
	JSON(value any, ctype ...string) error
	Set(key, value string)
	SetUserValue(key string, value any)
	SetUserContext(ctx context.Context)
	UserContext() context.Context
	UserValue(key string) any
	Native() any
}

type Router interface {
	Add(method, path string, handlers ...Handler)
	All(path string, handlers ...Handler)
	Delete(path string, handlers ...Handler)
	Get(path string, handlers ...Handler)
	Group(prefix string) Router
	Mount(prefix string, handlers ...Handler) Router
	Post(path string, handlers ...Handler)
	Put(path string, handlers ...Handler)
	Route(prefix string, fn func(Router))
	Use(path string, handlers ...Handler)
}

type RetryOptions struct {
	MaxAttempts int
	Delay       time.Duration
	ShouldRetry func(error) bool
}
