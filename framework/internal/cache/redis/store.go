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

	return NewStoreFromInstance(instanceCfg)
}

func NewStoreFromInstance(instanceCfg frameworkcache.InstanceConfig) (*Store, error) {
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

func (s *Store) Increment(key string) (int64, error) {
	return s.client.Incr(context.Background(), key).Result()
}

func (s *Store) Exists(key string) (bool, error) {
	count, err := s.client.Exists(context.Background(), key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *Store) Scan(pattern string) ([]string, error) {
	var (
		cursor uint64
		keys   []string
	)

	for {
		batch, next, err := s.client.Scan(context.Background(), cursor, pattern, 100).Result()
		if err != nil {
			return nil, err
		}
		keys = append(keys, batch...)
		cursor = next
		if cursor == 0 {
			return keys, nil
		}
	}
}

func (s *Store) TTL(key string) (time.Duration, error) {
	return s.client.TTL(context.Background(), key).Result()
}

func (s *Store) Publish(channel string, payload []byte) error {
	return s.client.Publish(context.Background(), channel, payload).Err()
}

func (s *Store) Subscribe(ctx context.Context, channel string) (frameworkcache.Subscription, error) {
	pubsub := s.client.Subscribe(ctx, channel)
	if _, err := pubsub.Receive(ctx); err != nil {
		_ = pubsub.Close()
		return nil, err
	}

	out := make(chan frameworkcache.Message)
	go func() {
		defer close(out)
		ch := pubsub.Channel()
		for {
			select {
			case <-ctx.Done():
				_ = pubsub.Close()
				return
			case msg, ok := <-ch:
				if !ok || msg == nil {
					return
				}
				out <- frameworkcache.Message{
					Channel: msg.Channel,
					Payload: []byte(msg.Payload),
				}
			}
		}
	}()

	return &subscription{pubsub: pubsub, ch: out}, nil
}

type subscription struct {
	pubsub *redis.PubSub
	ch     <-chan frameworkcache.Message
}

func (s *subscription) Channel() <-chan frameworkcache.Message {
	return s.ch
}

func (s *subscription) Close() error {
	return s.pubsub.Close()
}
