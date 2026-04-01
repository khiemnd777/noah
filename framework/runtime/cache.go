package runtime

import (
	"fmt"
	"sync"

	frameworkredis "github.com/khiemnd777/noah_framework/internal/cache/redis"
	frameworkcache "github.com/khiemnd777/noah_framework/pkg/cache"
)

var (
	cacheMu         sync.RWMutex
	configuredCache = map[string]frameworkcache.Backend{}
)

func NewCacheStore(cfg frameworkcache.Config) (frameworkcache.Store, error) {
	return frameworkredis.NewStore(cfg)
}

func ConfigureCache(cfg frameworkcache.Config) error {
	defaultInstance := cfg.DefaultInstance
	if defaultInstance == "" {
		defaultInstance = "cache"
	}

	backends := make(map[string]frameworkcache.Backend, len(cfg.Instances))
	for name, instanceCfg := range cfg.Instances {
		store, err := frameworkredis.NewStoreFromInstance(instanceCfg)
		if err != nil {
			return err
		}
		backends[name] = store
	}

	defaultBackend, ok := backends[defaultInstance]
	if !ok {
		return fmt.Errorf("cache instance %q is not configured", defaultInstance)
	}

	cacheMu.Lock()
	configuredCache = backends
	cacheMu.Unlock()

	frameworkcache.SetDefaultStore(defaultBackend)
	return nil
}

func ConfigureDefaultCache(cfg frameworkcache.Config) error {
	return ConfigureCache(cfg)
}

func CacheBackend(name string) (frameworkcache.Backend, error) {
	cacheMu.RLock()
	defer cacheMu.RUnlock()

	backend, ok := configuredCache[name]
	if !ok {
		return nil, fmt.Errorf("cache instance %q is not configured", name)
	}
	return backend, nil
}
