package fiber

import (
	"context"
	"mime/multipart"

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

func (r *Router) Add(method, path string, handlers ...frameworkhttp.Handler) {
	r.router.Add(method, path, wrapHandlers(handlers)...)
}

func (r *Router) Group(prefix string) frameworkhttp.Router {
	return &Router{router: r.router.Group(prefix)}
}

func (r *Router) Mount(prefix string, handlers ...frameworkhttp.Handler) frameworkhttp.Router {
	return &Router{router: r.router.Group(prefix, wrapHandlers(handlers)...)}
}

func (r *Router) Use(path string, handlers ...frameworkhttp.Handler) {
	args := make([]any, 0, len(handlers)+1)
	args = append(args, path)
	for _, handler := range wrapHandlers(handlers) {
		args = append(args, handler)
	}
	r.router.Use(args...)
}

func (r *Router) Route(prefix string, fn func(frameworkhttp.Router)) {
	fn(&Router{router: r.router.Group(prefix)})
}

func (c *Context) Get(key string, defaultValue ...string) string {
	return c.ctx.Get(key, defaultValue...)
}

func (c *Context) BodyParser(out any) error {
	return c.ctx.BodyParser(out)
}

func (c *Context) Body() []byte {
	return c.ctx.Body()
}

func (c *Context) FormFile(key string) (*multipart.FileHeader, error) {
	return c.ctx.FormFile(key)
}

func (c *Context) FormValue(key string, defaultValue ...string) string {
	return c.ctx.FormValue(key, defaultValue...)
}

func (c *Context) Locals(key any, value ...any) any {
	if len(value) > 0 {
		c.ctx.Locals(key, value[0])
	}
	return c.ctx.Locals(key)
}

func (c *Context) Method(override ...string) string {
	return c.ctx.Method(override...)
}

func (c *Context) Next() error {
	return c.ctx.Next()
}

func (c *Context) Path() string {
	return c.ctx.Path()
}

func (c *Context) Param(name string) string {
	return c.ctx.Params(name)
}

func (c *Context) Params(name string, defaultValue ...string) string {
	return c.ctx.Params(name, defaultValue...)
}

func (c *Context) ParamsInt(name string) (int, error) {
	return c.ctx.ParamsInt(name)
}

func (c *Context) Query(name string, defaultValue ...string) string {
	return c.ctx.Query(name, defaultValue...)
}

func (c *Context) QueryBool(name string, defaultValue ...bool) bool {
	return c.ctx.QueryBool(name, defaultValue...)
}

func (c *Context) QueryInt(name string, defaultValue ...int) int {
	return c.ctx.QueryInt(name, defaultValue...)
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

func (c *Context) SendFile(path string, compress ...bool) error {
	return c.ctx.SendFile(path, compress...)
}

func (c *Context) SendString(value string) error {
	return c.ctx.SendString(value)
}

func (c *Context) Status(status int) frameworkhttp.Context {
	c.ctx.Status(status)
	return c
}

func (c *Context) JSON(value any, ctype ...string) error {
	return c.ctx.JSON(value, ctype...)
}

func (c *Context) Set(key, value string) {
	c.ctx.Set(key, value)
}

func (c *Context) SetUserValue(key string, value any) {
	c.ctx.SetUserContext(context.WithValue(c.ctx.UserContext(), key, value))
}

func (c *Context) SetUserContext(ctx context.Context) {
	c.ctx.SetUserContext(ctx)
}

func (c *Context) UserContext() context.Context {
	return c.ctx.UserContext()
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
