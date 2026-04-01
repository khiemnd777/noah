package runtime

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
)

type HTTPClient struct {
	client *http.Client
}

type StatusError struct {
	StatusCode int
	Message    string
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("status %d: %s", e.StatusCode, e.Message)
}

func (e *StatusError) HTTPStatusCode() int {
	return e.StatusCode
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{client: http.DefaultClient}
}

var (
	httpClientOnce sync.Once
	httpClient     *HTTPClient
)

func DefaultHTTPClient() *HTTPClient {
	httpClientOnce.Do(func() {
		httpClient = NewHTTPClient()
	})
	return httpClient
}

func (c *HTTPClient) CallRequestWithPort(ctx context.Context, method, module string, port int, path, token, internalToken string, body, out any, opts ...frameworkhttp.RetryOptions) error {
	appCfg, err := LoadYAML[AppConfig](APIPath("config.yaml"))
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s:%d%s", appCfg.Server.Host, port, path)
	reqBody, err := marshalBody(body)
	if err != nil {
		return err
	}

	retry := defaultRetryOptions()
	if len(opts) > 0 {
		retry = mergeRetryOptions(opts[0], retry)
	}

	var lastErr error
	for i := 0; i < retry.MaxAttempts; i++ {
		_, err := RunWithCircuitBreaker(fmt.Sprintf("%s:%s", module, path), func(_ context.Context) (any, error) {
			req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
			if err != nil {
				return nil, err
			}

			setHeaders(req, token, internalToken)

			resp, err := c.client.Do(req)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()

			if resp.StatusCode >= http.StatusBadRequest {
				return nil, parseErrorResponse(resp)
			}

			if out != nil {
				return nil, json.NewDecoder(resp.Body).Decode(out)
			}

			return nil, nil
		})
		if err == nil || !retry.ShouldRetry(err) {
			return err
		}
		lastErr = err
		time.Sleep(retry.Delay)
	}

	return lastErr
}

func (c *HTTPClient) CallRequest(ctx context.Context, method, module, path, token, internalToken string, body, out any, opts ...frameworkhttp.RetryOptions) error {
	appCfg, err := LoadYAML[AppConfig](APIPath("config.yaml"))
	if err != nil {
		return err
	}
	return c.CallRequestWithPort(ctx, method, module, appCfg.Server.Port, path, token, internalToken, body, out, opts...)
}

func (c *HTTPClient) CallGet(ctx context.Context, module, path, token, internalToken string, out any, opts ...frameworkhttp.RetryOptions) error {
	return c.CallRequest(ctx, http.MethodGet, module, path, token, internalToken, nil, out, opts...)
}

func (c *HTTPClient) CallPost(ctx context.Context, module, path, token, internalToken string, body, out any, opts ...frameworkhttp.RetryOptions) error {
	return c.CallRequest(ctx, http.MethodPost, module, path, token, internalToken, body, out, opts...)
}

func (c *HTTPClient) CallPostAsync(ctx context.Context, module, path, token, internalToken string, body, out any, opts ...frameworkhttp.RetryOptions) {
	go func() {
		_ = c.CallPost(ctx, module, path, token, internalToken, body, out, opts...)
	}()
}

func marshalBody(body any) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func setHeaders(req *http.Request, token, internalToken string) {
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if internalToken != "" {
		req.Header.Set("X-Internal-Token", internalToken)
	}
}

func parseErrorResponse(resp *http.Response) error {
	if resp == nil || resp.Body == nil {
		return errors.New("empty response body")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body failed: %w", err)
	}
	return &StatusError{
		StatusCode: resp.StatusCode,
		Message:    string(bytes.TrimSpace(body)),
	}
}
