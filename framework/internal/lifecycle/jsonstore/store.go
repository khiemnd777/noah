package jsonstore

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"syscall"

	frameworklifecycle "github.com/khiemnd777/noah_framework/pkg/lifecycle"
)

type Store struct {
	path string
	mu   sync.Mutex
}

func New(path string) *Store {
	return &Store{path: path}
}

func (s *Store) Load() (frameworklifecycle.Registry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return withLock(s.path, func() (frameworklifecycle.Registry, error) {
		return s.loadUnlocked()
	})
}

func (s *Store) Save(reg frameworklifecycle.Registry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return withLockErr(s.path, func() error {
		return s.saveUnlocked(reg)
	})
}

func (s *Store) Update(update func(frameworklifecycle.Registry)) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return withLockErr(s.path, func() error {
		reg, err := s.loadUnlocked()
		if err != nil {
			return err
		}
		update(reg)
		return s.saveUnlocked(reg)
	})
}

func (s *Store) loadUnlocked() (frameworklifecycle.Registry, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return frameworklifecycle.Registry{}, nil
		}
		return nil, err
	}

	reg := frameworklifecycle.Registry{}
	if err := json.Unmarshal(data, &reg); err != nil {
		return nil, err
	}
	return reg, nil
}

func (s *Store) saveUnlocked(reg frameworklifecycle.Registry) error {
	data, err := json.MarshalIndent(reg, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}

	tmpPath := s.path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmpPath, s.path)
}

func withLock[T any](path string, fn func() (T, error)) (T, error) {
	lockPath := path + ".lock"
	if err := os.MkdirAll(filepath.Dir(lockPath), 0o755); err != nil {
		var zero T
		return zero, err
	}

	lockFile, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		var zero T
		return zero, err
	}
	defer lockFile.Close()

	if err := syscall.Flock(int(lockFile.Fd()), syscall.LOCK_EX); err != nil {
		var zero T
		return zero, err
	}
	defer func() {
		_ = syscall.Flock(int(lockFile.Fd()), syscall.LOCK_UN)
	}()

	return fn()
}

func withLockErr(path string, fn func() error) error {
	_, err := withLock(path, func() (struct{}, error) {
		return struct{}{}, fn()
	})
	return err
}
