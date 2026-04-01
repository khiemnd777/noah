package runtime

import (
	frameworkjsonstore "github.com/khiemnd777/noah_framework/internal/lifecycle/jsonstore"
	frameworklifecycle "github.com/khiemnd777/noah_framework/pkg/lifecycle"
)

func NewLifecycleStore(path string) frameworklifecycle.Store {
	return frameworkjsonstore.New(path)
}

func ConfigureDefaultLifecycleStore(path string) frameworklifecycle.Store {
	store := NewLifecycleStore(path)
	frameworklifecycle.SetDefaultStore(store)
	return store
}
