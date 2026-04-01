package circuitbreaker

import (
	"context"
	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"
)

var ErrClientResponse = frameworkruntime.ErrClientResponse

func Init() {
	_ = frameworkruntime.InitCircuitBreaker()
}

func Run(name string, fn func(context.Context) (interface{}, error)) (interface{}, error) {
	return frameworkruntime.RunWithCircuitBreaker(name, func(ctx context.Context) (any, error) {
		return fn(ctx)
	})
}
