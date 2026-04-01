package cache

import (
	"strings"

	"github.com/khiemnd777/noah_framework/shared/logger"
	frameworkcache "github.com/khiemnd777/noah_framework/pkg/cache"
)

// UpdateAndInvalidate runs DB update and clears cache
func UpdateAndInvalidate(key string, updateFn func() error) error {
	err := updateFn()
	if err != nil {
		return err
	}

	ensureFrameworkCache()
	store, err := frameworkcache.DefaultStore()
	if err != nil {
		return err
	}
	return store.Delete(key)
}

// UpdateManyAndInvalidate runs updateFn and deletes multiple cache keys if update succeeds.
func UpdateManyAndInvalidate(keys []string, updateFn func() error) error {
	err := updateFn()
	if err != nil {
		return err
	}

	InvalidateKeys(keys...)
	return nil
}

func isPatternKey(key string) bool {
	return strings.Contains(key, "*")
}

// InvalidateKeys deletes multiple cache keys or patterns.
func InvalidateKeys(keys ...string) {
	ensureFrameworkCache()
	store, err := frameworkcache.DefaultStore()
	if err != nil {
		logger.Warn("❌ Framework cache store unavailable: " + err.Error())
		return
	}

	for _, key := range keys {
		if isPatternKey(key) {
			if err := store.DeleteByPattern(key); err != nil {
				logger.Warn("❌ Failed to invalidate pattern cache: " + key)
			} else {
				logger.Info("🧹 Pattern cache invalidated: " + key)
			}
			continue
		}

		if err := store.Delete(key); err != nil {
			logger.Warn("❌ Failed to invalidate cache: " + key)
		} else {
			logger.Info("🧹 Cache invalidated: " + key)
		}
	}
}
