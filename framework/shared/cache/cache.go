package cache

import (
	"sync"
	"time"

	"github.com/khiemnd777/noah_framework/shared/config"
	frameworkcache "github.com/khiemnd777/noah_framework/pkg/cache"
	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"
)

var cacheSetup sync.Once

func ensureFrameworkCache() {
	cacheSetup.Do(func() {
		instances := make(map[string]frameworkcache.InstanceConfig, len(config.Get().Redis.Instances))
		for name, instance := range config.Get().Redis.Instances {
			instances[name] = frameworkcache.InstanceConfig{
				Host:      instance.Host,
				Port:      instance.Port,
				Username:  instance.Username,
				Password:  instance.Password,
				DB:        instance.DB,
				IsCluster: instance.IsCluster,
				UseTLS:    instance.UseTLS,
			}
		}

		_ = frameworkruntime.ConfigureDefaultCache(frameworkcache.Config{
			DefaultInstance: "cache",
			Instances:       instances,
		})
	})
}

func Get[T any](key string, ttl time.Duration, fallback func() (*T, error)) (*T, error) {
	ensureFrameworkCache()
	return frameworkcache.GetOrSet(key, ttl, fallback)
}

func GetList[T any](key string, ttl time.Duration, fallback func() ([]T, error)) ([]T, error) {
	ensureFrameworkCache()

	wrappedFallback := func() (*[]T, error) {
		list, err := fallback()
		if err != nil {
			return nil, err
		}
		return &list, nil
	}

	ptr, err := frameworkcache.GetOrSet(key, ttl, wrappedFallback)
	if err != nil {
		return nil, err
	}
	return *ptr, nil
}

type PaginatedCachedList[T any] struct {
	Items   []T  `json:"items"`
	HasMore bool `json:"hasMore"`
}

func GetListWithHasMore[T any](key string, ttl time.Duration, fallback func() ([]T, bool, error)) ([]T, bool, error) {
	ensureFrameworkCache()

	wrapped := func() (*PaginatedCachedList[T], error) {
		items, hasMore, err := fallback()
		if err != nil {
			return nil, err
		}
		return &PaginatedCachedList[T]{Items: items, HasMore: hasMore}, nil
	}

	result, err := frameworkcache.GetOrSet(key, ttl, wrapped)
	if err != nil {
		return nil, false, err
	}
	return result.Items, result.HasMore, nil
}
