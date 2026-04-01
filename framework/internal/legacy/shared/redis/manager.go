package redis

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/khiemnd777/noah_framework/internal/legacy/shared/config"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/logger"
	frameworkcache "github.com/khiemnd777/noah_framework/pkg/cache"
	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"
)

var (
	initOnce sync.Once
)

func Init() {
	initOnce.Do(func() {
		if err := frameworkruntime.ConfigureCache(toFrameworkCacheConfig(config.Get().Redis)); err != nil {
			logger.Warn(fmt.Sprintf("⚠️ Redis infra not connected: %v", err))
		}
	})
}

func backend(name string) (frameworkcache.Backend, error) {
	Init()
	return frameworkruntime.CacheBackend(name)
}

func Get(name, key string) (string, error) {
	rdb, err := backend(name)
	if err != nil {
		return "", err
	}

	val, err := rdb.Get(key)
	if err != nil {
		logger.Warn("❌ Redis GET error: "+key, err)
		return "", err
	}
	if len(val) == 0 {
		logger.Info("🔍 Redis GET: key not found: " + key)
		return "", nil
	}
	return string(val), nil
}

func GetInt(name, key string) (int, error) {
	val, err := Get(name, key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(val)
}

func Set(name, key string, value interface{}, ttl time.Duration) error {
	rdb, err := backend(name)
	if err != nil {
		return err
	}
	err = rdb.Set(key, []byte(fmt.Sprint(value)), ttl)
	if err != nil {
		logger.Warn("❌ Redis SET error: "+key, err)
	}
	return err
}

func Incr(name, key string) (int64, error) {
	rdb, err := backend(name)
	if err != nil {
		return 0, err
	}
	val, err := rdb.Increment(key)
	if err != nil {
		logger.Warn("❌ Redis SET error: "+key, err)
	}
	return val, err
}

func Del(name, key string) error {
	rdb, err := backend(name)
	if err != nil {
		return err
	}
	err = rdb.Delete(key)
	if err != nil {
		logger.Warn("❌ Redis DEL error: "+key, err)
	}
	return err
}

func DelByPattern(name string, pattern string) error {
	rdb, err := backend(name)
	if err != nil {
		return err
	}
	return rdb.DeleteByPattern(pattern)
}

func Exists(name, key string) (bool, error) {
	rdb, err := backend(name)
	if err != nil {
		return false, err
	}
	ok, err := rdb.Exists(key)
	if err != nil {
		logger.Warn("❌ Redis EXISTS error: "+key, err)
		return false, err
	}
	return ok, nil
}

func ScanKeys(name, pattern string) ([]string, error) {
	rdb, err := backend(name)
	if err != nil {
		return nil, err
	}
	return rdb.Scan(pattern)
}

func TTL(name, key string) (time.Duration, error) {
	rdb, err := backend(name)
	if err != nil {
		return 0, err
	}
	ttl, err := rdb.TTL(key)
	if err != nil {
		logger.Warn("❌ Redis TTL error: "+key, err)
		return 0, err
	}
	return ttl, nil
}

func InitFromConfig(cfg config.RedisConfig) {
	if err := frameworkruntime.ConfigureCache(toFrameworkCacheConfig(cfg)); err != nil {
		logger.Warn(fmt.Sprintf("⚠️ Redis infra not connected: %v", err))
		return
	}
	for name, rc := range cfg.Instances {
		logger.Info(fmt.Sprintf("✅ Redis [%s] configured: %s:%d", name, rc.Host, rc.Port))
	}
}

func toFrameworkCacheConfig(cfg config.RedisConfig) frameworkcache.Config {
	instances := make(map[string]frameworkcache.InstanceConfig, len(cfg.Instances))
	for name, instance := range cfg.Instances {
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

	defaultInstance := "cache"
	if _, ok := instances["cache"]; !ok {
		for name := range instances {
			defaultInstance = name
			break
		}
	}

	return frameworkcache.Config{
		DefaultInstance: defaultInstance,
		Instances:       instances,
	}
}
