package context

type Resolver interface {
	Set(name string, value any)
	Resolve(name string) (any, bool)
}
