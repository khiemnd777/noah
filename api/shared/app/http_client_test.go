package app

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestHttpClientPostDoesNotRetryByDefault(t *testing.T) {
	initAppTestConfig(t, "127.0.0.1")

	var attempts int32
	client := &HttpClient{client: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		atomic.AddInt32(&attempts, 1)
		return jsonResponse(http.StatusInternalServerError, "boom"), nil
	})}}

	err := client.CallRequestWithPort(context.Background(), http.MethodPost, "test", 1234, "/retry", "", "", map[string]string{"value": "x"}, nil)
	if err == nil {
		t.Fatal("CallRequestWithPort() error = nil, want non-nil")
	}
	if got := atomic.LoadInt32(&attempts); got != 1 {
		t.Fatalf("attempts = %d, want 1", got)
	}
}

func TestHttpClientGetRetriesByDefault(t *testing.T) {
	initAppTestConfig(t, "127.0.0.1")

	var attempts int32
	client := &HttpClient{client: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		attempt := atomic.AddInt32(&attempts, 1)
		if attempt < 3 {
			return jsonResponse(http.StatusInternalServerError, "boom"), nil
		}
		return jsonResponse(http.StatusOK, `{"status":"ok"}`), nil
	})}}

	var out struct {
		Status string `json:"status"`
	}
	err := client.CallRequestWithPort(context.Background(), http.MethodGet, "test", 1234, "/retry", "", "", nil, &out)
	if err != nil {
		t.Fatalf("CallRequestWithPort() error = %v", err)
	}
	if got := atomic.LoadInt32(&attempts); got != 3 {
		t.Fatalf("attempts = %d, want 3", got)
	}
	if out.Status != "ok" {
		t.Fatalf("out.Status = %q, want %q", out.Status, "ok")
	}
}

func TestHttpClientPostCanExplicitlyRetryAndRebuildBody(t *testing.T) {
	initAppTestConfig(t, "127.0.0.1")

	var attempts int32
	bodies := make([]string, 0, 3)
	client := &HttpClient{client: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		bodyBytes, _ := io.ReadAll(r.Body)
		bodies = append(bodies, string(bodyBytes))

		attempt := atomic.AddInt32(&attempts, 1)
		if attempt < 3 {
			return jsonResponse(http.StatusInternalServerError, "boom"), nil
		}
		return jsonResponse(http.StatusOK, ""), nil
	})}}

	payload := struct {
		Value string `json:"value"`
	}{Value: "retry-body"}
	err := client.CallRequestWithPort(
		context.Background(),
		http.MethodPost,
		"test",
		1234,
		"/retry",
		"",
		"",
		payload,
		nil,
		RetryOptions{
			MaxAttempts: 3,
			Delay:       time.Millisecond,
		},
	)
	if err != nil {
		t.Fatalf("CallRequestWithPort() error = %v", err)
	}
	if got := atomic.LoadInt32(&attempts); got != 3 {
		t.Fatalf("attempts = %d, want 3", got)
	}
	if len(bodies) != 3 {
		t.Fatalf("len(bodies) = %d, want 3", len(bodies))
	}
	for i, body := range bodies {
		if body != `{"value":"retry-body"}` {
			t.Fatalf("body[%d] = %q, want %q", i, body, `{"value":"retry-body"}`)
		}
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func jsonResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
