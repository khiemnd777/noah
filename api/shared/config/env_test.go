package config

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestEnsureEnvLoadedPrefersRootSharedEnv(t *testing.T) {
	t.Helper()

	repoRoot := t.TempDir()
	apiDir := filepath.Join(repoRoot, "api")
	feDir := filepath.Join(repoRoot, "fe")

	mustMkdirAll(t, apiDir)
	mustMkdirAll(t, feDir)
	mustWriteFile(t, filepath.Join(repoRoot, ".env"), "APP_FE_ORIGIN=http://localhost:5173\n")
	mustWriteFile(t, filepath.Join(apiDir, ".env"), "APP_FE_ORIGIN=http://localhost:9999\nAPI_ONLY=true\nM_MAIN_CLIENT_BASE_URL=http://legacy.local\n")

	withWorkingDir(t, apiDir)
	withCleanEnv(t, "APP_ENV", "APP_FE_ORIGIN", "API_ONLY", "M_MAIN_CLIENT_BASE_URL")

	loadDotEnvOnce = sync.Once{}

	if err := EnsureEnvLoaded(); err != nil {
		t.Fatalf("EnsureEnvLoaded() error = %v", err)
	}

	if got := os.Getenv("APP_FE_ORIGIN"); got != "http://localhost:5173" {
		t.Fatalf("APP_FE_ORIGIN = %q, want %q", got, "http://localhost:5173")
	}

	if got := os.Getenv("API_ONLY"); got != "true" {
		t.Fatalf("API_ONLY = %q, want %q", got, "true")
	}

	if got := os.Getenv("M_MAIN_CLIENT_BASE_URL"); got != "http://localhost:5173" {
		t.Fatalf("M_MAIN_CLIENT_BASE_URL = %q, want %q", got, "http://localhost:5173")
	}
}

func TestEnsureEnvLoadedUsesRootProdSharedEnv(t *testing.T) {
	t.Helper()

	repoRoot := t.TempDir()
	apiDir := filepath.Join(repoRoot, "api")
	feDir := filepath.Join(repoRoot, "fe")

	mustMkdirAll(t, apiDir)
	mustMkdirAll(t, feDir)
	mustWriteFile(t, filepath.Join(repoRoot, ".env.prod"), "APP_FE_ORIGIN=https://admin.example.com\n")
	mustWriteFile(t, filepath.Join(apiDir, ".env.prod"), "API_ONLY=prod\nM_MAIN_CLIENT_BASE_URL=http://legacy.local\n")

	withWorkingDir(t, apiDir)
	withCleanEnv(t, "APP_FE_ORIGIN", "API_ONLY", "M_MAIN_CLIENT_BASE_URL")
	t.Setenv("APP_ENV", "production")

	loadDotEnvOnce = sync.Once{}

	if err := EnsureEnvLoaded(); err != nil {
		t.Fatalf("EnsureEnvLoaded() error = %v", err)
	}

	if got := os.Getenv("APP_FE_ORIGIN"); got != "https://admin.example.com" {
		t.Fatalf("APP_FE_ORIGIN = %q, want %q", got, "https://admin.example.com")
	}

	if got := os.Getenv("API_ONLY"); got != "prod" {
		t.Fatalf("API_ONLY = %q, want %q", got, "prod")
	}

	if got := os.Getenv("M_MAIN_CLIENT_BASE_URL"); got != "https://admin.example.com" {
		t.Fatalf("M_MAIN_CLIENT_BASE_URL = %q, want %q", got, "https://admin.example.com")
	}
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v", path, err)
	}
}

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("os.MkdirAll(%q) error = %v", path, err)
	}
}

func withWorkingDir(t *testing.T, dir string) {
	t.Helper()

	previous, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("os.Chdir(%q) error = %v", dir, err)
	}

	t.Cleanup(func() {
		_ = os.Chdir(previous)
	})
}

func withCleanEnv(t *testing.T, keys ...string) {
	t.Helper()

	snapshots := make(map[string]*string, len(keys))
	for _, key := range keys {
		if value, ok := os.LookupEnv(key); ok {
			copyValue := value
			snapshots[key] = &copyValue
		} else {
			snapshots[key] = nil
		}

		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("os.Unsetenv(%q) error = %v", key, err)
		}
	}

	t.Cleanup(func() {
		for _, key := range keys {
			if value := snapshots[key]; value != nil {
				_ = os.Setenv(key, *value)
				continue
			}

			_ = os.Unsetenv(key)
		}
	})
}
