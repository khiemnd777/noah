package runtime

import (
	frameworkredis "github.com/khiemnd777/noah_framework/internal/cache/redis"
	frameworkcache "github.com/khiemnd777/noah_framework/pkg/cache"
)

func NewCacheStore(cfg frameworkcache.Config) (frameworkcache.Store, error) {
	return frameworkredis.NewStore(cfg)
}

func ConfigureDefaultCache(cfg frameworkcache.Config) error {
	store, err := NewCacheStore(cfg)
	if err != nil {
		return err
	}

	frameworkcache.SetDefaultStore(store)
	return nil
}
