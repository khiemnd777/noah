package fiber

import (
	"context"

	"github.com/gofiber/fiber/v2"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

type Router struct {
	router fiber.Router
}

type Context struct {
	ctx *fiber.Ctx
}

func NewRouter(router fiber.Router) *Router {
	return &Router{router: router}
}

func (r *Router) Get(path string, handlers ...frameworkhttp.Handler) {
	r.router.Get(path, wrapHandlers(handlers)...)
}

func (r *Router) Post(path string, handlers ...frameworkhttp.Handler) {
	r.router.Post(path, wrapHandlers(handlers)...)
}

func (r *Router) Put(path string, handlers ...frameworkhttp.Handler) {
	r.router.Put(path, wrapHandlers(handlers)...)
}

func (r *Router) Delete(path string, handlers ...frameworkhttp.Handler) {
	r.router.Delete(path, wrapHandlers(handlers)...)
}

func (r *Router) All(path string, handlers ...frameworkhttp.Handler) {
	r.router.All(path, wrapHandlers(handlers)...)
}

func (r *Router) Group(prefix string) frameworkhttp.Router {
	return &Router{router: r.router.Group(prefix)}
}

func (r *Router) Use(path string, handlers ...frameworkhttp.Handler) {
	args := make([]any, 0, len(handlers)+1)
	args = append(args, path)
	for _, handler := range wrapHandlers(handlers) {
		args = append(args, handler)
	}
	r.router.Use(args...)
}

func (c *Context) Get(key string) string {
	return c.ctx.Get(key)
}

func (c *Context) Locals(key string, value ...any) any {
	if len(value) > 0 {
		c.ctx.Locals(key, value[0])
	}
	return c.ctx.Locals(key)
}

func (c *Context) Method() string {
	return c.ctx.Method()
}

func (c *Context) Path() string {
	return c.ctx.Path()
}

func (c *Context) Param(name string) string {
	return c.ctx.Params(name)
}

func (c *Context) Header(name string) string {
	return c.ctx.Get(name)
}

func (c *Context) OriginalURL() string {
	return c.ctx.OriginalURL()
}

func (c *Context) QueryString() string {
	return string(c.ctx.Context().URI().QueryString())
}

func (c *Context) Local(key string) any {
	return c.ctx.Locals(key)
}

func (c *Context) SetLocal(key string, value any) {
	c.ctx.Locals(key, value)
}

func (c *Context) SendStatus(status int) error {
	return c.ctx.SendStatus(status)
}

func (c *Context) SendString(value string) error {
	return c.ctx.SendString(value)
}

func (c *Context) Status(status int) frameworkhttp.Context {
	c.ctx.Status(status)
	return c
}

func (c *Context) JSON(value any) error {
	return c.ctx.JSON(value)
}

func (c *Context) Set(key, value string) {
	c.ctx.Set(key, value)
}

func (c *Context) SetUserValue(key string, value any) {
	c.ctx.SetUserContext(context.WithValue(c.ctx.UserContext(), key, value))
}

func (c *Context) UserValue(key string) any {
	return c.ctx.UserContext().Value(key)
}

func (c *Context) Native() any {
	return c.ctx
}

func wrapHandlers(handlers []frameworkhttp.Handler) []fiber.Handler {
	result := make([]fiber.Handler, 0, len(handlers))
	for _, handler := range handlers {
		handler := handler
		result = append(result, func(ctx *fiber.Ctx) error {
			return handler(&Context{ctx: ctx})
		})
	}
	return result
}
