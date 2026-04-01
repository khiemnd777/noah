package runtime

import (
	"context"
	"errors"
	"fmt"
	"log"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/sony/gobreaker"
)

var (
	breakers       sync.Map
	breakerFactory func(name string) *gobreaker.CircuitBreaker
)

var ErrClientResponse = errors.New("client response")

type statusCoder interface {
	HTTPStatusCode() int
}

func InitCircuitBreaker() error {
	appCfg, err := LoadYAML[AppConfig](APIPath("config.yaml"))
	if err != nil {
		return err
	}

	cbCfg := appCfg.CircuitBreaker
	breakerFactory = func(name string) *gobreaker.CircuitBreaker {
		settings := gobreaker.Settings{
			Name:     name,
			Interval: cbCfg.Interval,
			Timeout:  cbCfg.Timeout,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				return counts.ConsecutiveFailures >= uint32(cbCfg.ConsecutiveFailures)
			},
			IsSuccessful: func(err error) bool {
				if err == nil || errors.Is(err, ErrClientResponse) {
					return true
				}
				var sc statusCoder
				if errors.As(err, &sc) && sc.HTTPStatusCode() < 500 {
					return true
				}
				return false
			},
			OnStateChange: func(cbName string, from gobreaker.State, to gobreaker.State) {
				log.Printf("circuit state changed [%s]: %s -> %s", cbName, from.String(), to.String())
			},
		}
		return gobreaker.NewCircuitBreaker(settings)
	}

	return nil
}

func RunWithCircuitBreaker(name string, fn func(context.Context) (any, error)) (any, error) {
	var panicErr error
	breaker := getBreaker(name)

	result, err := breaker.Execute(func() (any, error) {
		defer func() {
			if r := recover(); r != nil {
				panicErr = fmt.Errorf("panic: %v", r)
				log.Printf("panic in circuit breaker [%s]: %v\n%s", name, r, string(debug.Stack()))
			}
		}()
		return fn(context.Background())
	})

	if panicErr != nil {
		return nil, panicErr
	}

	switch {
	case errors.Is(err, gobreaker.ErrOpenState):
		log.Printf("circuit open: blocked call [%s]", name)
		return nil, err
	case errors.Is(err, ErrClientResponse):
		return nil, err
	case err != nil:
		if strings.Contains(err.Error(), "client error") {
			return nil, err
		}
		return nil, err
	default:
		return result, nil
	}
}

func getBreaker(name string) *gobreaker.CircuitBreaker {
	if breakerFactory == nil {
		if err := InitCircuitBreaker(); err != nil {
			panic(err)
		}
	}
	if cb, ok := breakers.Load(name); ok {
		return cb.(*gobreaker.CircuitBreaker)
	}
	created := breakerFactory(name)
	actual, _ := breakers.LoadOrStore(name, created)
	return actual.(*gobreaker.CircuitBreaker)
}
