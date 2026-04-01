package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	frameworkcache "github.com/khiemnd777/noah_framework/pkg/cache"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
}

func NewStore(cfg frameworkcache.Config) (*Store, error) {
	instanceName := cfg.DefaultInstance
	if instanceName == "" {
		instanceName = "cache"
	}

	instanceCfg, ok := cfg.Instances[instanceName]
	if !ok {
		return nil, fmt.Errorf("cache instance %q is not configured", instanceName)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     instanceCfg.Host + ":" + strconv.Itoa(instanceCfg.Port),
		Username: instanceCfg.Username,
		Password: instanceCfg.Password,
		DB:       instanceCfg.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &Store{client: client}, nil
}

func (s *Store) Get(key string) ([]byte, error) {
	value, err := s.client.Get(context.Background(), key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (s *Store) Set(key string, value []byte, ttl time.Duration) error {
	return s.client.Set(context.Background(), key, value, ttl).Err()
}

func (s *Store) Delete(key string) error {
	return s.client.Del(context.Background(), key).Err()
}

func (s *Store) DeleteByPattern(pattern string) error {
	var cursor uint64

	for {
		keys, next, err := s.client.Scan(context.Background(), cursor, pattern, 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := s.client.Del(context.Background(), keys...).Err(); err != nil {
				return err
			}
		}
		cursor = next
		if cursor == 0 {
			return nil
		}
	}
}
