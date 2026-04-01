package lifecycle

import (
	"errors"
	"os"
	"sync"
	"time"
)

type ModuleInfo struct {
	PID      int       `json:"pid"`
	Host     string    `json:"host"`
	Port     int       `json:"port"`
	RunAt    time.Time `json:"run_at"`
	External bool      `json:"external"`
	Route    string    `json:"router"`
}

type Registry map[string]ModuleInfo

type Store interface {
	Load() (Registry, error)
	Save(Registry) error
	Update(func(Registry)) error
}

type Hook func() error

type Manager interface {
	AppendOnStart(Hook)
	AppendOnStop(Hook)
	RunStart() error
	RunStop() error
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
		return nil, errors.New("framework lifecycle store is not configured")
	}
	return defaultStore, nil
}

func Register(name, route, host string, port int, external bool) error {
	store, err := DefaultStore()
	if err != nil {
		return err
	}

	return store.Update(func(reg Registry) {
		reg[name] = ModuleInfo{
			PID:      os.Getpid(),
			Host:     host,
			Port:     port,
			Route:    route,
			RunAt:    time.Now(),
			External: external,
		}
	})
}
