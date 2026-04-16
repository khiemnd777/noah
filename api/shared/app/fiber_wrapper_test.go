package app

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/khiemnd777/noah_api/shared/config"
)

func TestRouterPostDoesNotRetryMutationByDefault(t *testing.T) {
	initAppTestConfig(t, "127.0.0.1")

	var attempts int32
	fApp := fiber.New()
	RouterPost(fApp, "/post-no-retry", func(c *fiber.Ctx) error {
		atomic.AddInt32(&attempts, 1)
		return fmt.Errorf("post-processing failed")
	})

	resp := performFiberRequest(t, fApp, http.MethodPost, "/post-no-retry")
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
	}
	if got := atomic.LoadInt32(&attempts); got != 1 {
		t.Fatalf("attempts = %d, want 1", got)
	}
}

func TestRouterGetRetriesTransientFailuresByDefault(t *testing.T) {
	initAppTestConfig(t, "127.0.0.1")

	var attempts int32
	fApp := fiber.New()
	RouterGet(fApp, "/get-retry", func(c *fiber.Ctx) error {
		attempt := atomic.AddInt32(&attempts, 1)
		if attempt < 3 {
			return fmt.Errorf("transient failure")
		}
		return c.SendString("ok")
	})

	resp := performFiberRequest(t, fApp, http.MethodGet, "/get-retry")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if got := atomic.LoadInt32(&attempts); got != 3 {
		t.Fatalf("attempts = %d, want 3", got)
	}
}

func TestRouterPostWithOptionsAllowsExplicitRetry(t *testing.T) {
	initAppTestConfig(t, "127.0.0.1")

	var attempts int32
	fApp := fiber.New()
	RouterPostWithOptions(fApp, "/post-opt-in", RetryOptions{
		MaxAttempts: 3,
		Delay:       time.Millisecond,
	}, func(c *fiber.Ctx) error {
		attempt := atomic.AddInt32(&attempts, 1)
		if attempt < 3 {
			return fmt.Errorf("transient failure")
		}
		return c.SendString("ok")
	})

	resp := performFiberRequest(t, fApp, http.MethodPost, "/post-opt-in")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if got := atomic.LoadInt32(&attempts); got != 3 {
		t.Fatalf("attempts = %d, want 3", got)
	}
}

func TestRouterGetDoesNotRetryRecoveredPanic(t *testing.T) {
	initAppTestConfig(t, "127.0.0.1")

	var attempts int32
	fApp := fiber.New()
	RouterGet(fApp, "/panic-no-retry", func(c *fiber.Ctx) error {
		atomic.AddInt32(&attempts, 1)
		panic("boom")
	})

	resp := performFiberRequest(t, fApp, http.MethodGet, "/panic-no-retry")
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
	}
	if got := atomic.LoadInt32(&attempts); got != 1 {
		t.Fatalf("attempts = %d, want 1", got)
	}
}

func performFiberRequest(t *testing.T, app *fiber.App, method, path string) *http.Response {
	t.Helper()

	req := httptest.NewRequest(method, path, nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	return resp
}

func initAppTestConfig(t *testing.T, host string) {
	t.Helper()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")
	configBody := fmt.Sprintf(`project:
  name: test
server:
  host: %s
  port: 0
retry:
  maxattempts: 3
  delay: 1ms
circuitbreaker:
  interval: 1s
  timeout: 1s
  consecutivefailures: 10
`, host)

	if err := os.WriteFile(configPath, []byte(configBody), 0o644); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	if err := config.Init(configPath); err != nil {
		t.Fatalf("config.Init() error = %v", err)
	}
}
