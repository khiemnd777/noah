package cache

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
)

type InstanceConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	DB        int
	IsCluster bool
	UseTLS    bool
}

type Config struct {
	DefaultInstance string
	Instances       map[string]InstanceConfig
}

type Store interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, ttl time.Duration) error
	Delete(key string) error
	DeleteByPattern(pattern string) error
}

type PaginatedList[T any] struct {
	Items   []T  `json:"items"`
	HasMore bool `json:"hasMore"`
}

var (
	defaultStore Store
	storeMu      sync.RWMutex
)

func SetDefaultStore(store Store) {
	storeMu.Lock()
	defer storeMu.Unlock()
	defaultStore = store
}

func DefaultStore() (Store, error) {
	storeMu.RLock()
	defer storeMu.RUnlock()

	if defaultStore == nil {
		return nil, errors.New("framework cache store is not configured")
	}
	return defaultStore, nil
}

func GetOrSet[T any](key string, ttl time.Duration, fallback func() (*T, error)) (*T, error) {
	store, err := DefaultStore()
	if err != nil {
		return nil, err
	}

	payload, err := store.Get(key)
	if err != nil {
		return nil, err
	}
	if len(payload) > 0 {
		var out T
		if err := json.Unmarshal(payload, &out); err != nil {
			return nil, err
		}
		return &out, nil
	}

	value, err := fallback()
	if err != nil {
		return nil, err
	}

	payload, err = json.Marshal(value)
	if err != nil {
		return nil, err
	}
	if err := store.Set(key, payload, ttl); err != nil {
		return nil, err
	}

	return value, nil
}
