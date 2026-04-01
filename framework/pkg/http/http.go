package http

import "time"

type Handler func(Context) error

type Context interface {
	Get(key string) string
	Locals(key string, value ...any) any
	Method() string
	OriginalURL() string
	Path() string
	Param(key string) string
	QueryString() string
	SendStatus(status int) error
	SendString(body string) error
	Status(status int) Context
	JSON(value any) error
	Set(key, value string)
	SetUserValue(key string, value any)
	UserValue(key string) any
	Native() any
}

type Router interface {
	All(path string, handlers ...Handler)
	Delete(path string, handlers ...Handler)
	Get(path string, handlers ...Handler)
	Group(prefix string) Router
	Post(path string, handlers ...Handler)
	Put(path string, handlers ...Handler)
	Use(path string, handlers ...Handler)
}

type RetryOptions struct {
	MaxAttempts int
	Delay       time.Duration
	ShouldRetry func(error) bool
}
