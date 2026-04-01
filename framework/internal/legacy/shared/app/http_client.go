package app

import (
	"context"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"
)

// How to use:
// app.GetHttpClient().CallXxx

type HttpClient struct {
	client *frameworkruntime.HTTPClient
}

type StatusError = frameworkruntime.StatusError

func NewHttpClient() *HttpClient {
	return &HttpClient{
		client: frameworkruntime.NewHTTPClient(),
	}
}

func GetHttpClient() *HttpClient {
	return &HttpClient{client: frameworkruntime.DefaultHTTPClient()}
}

func (c *HttpClient) CallRequestWithPort(ctx context.Context, method, module string, port int, path, token, internalToken string, body, out any, opts ...RetryOptions) error {
	runtimeOpts := make([]frameworkhttp.RetryOptions, 0, len(opts))
	for _, opt := range opts {
		runtimeOpts = append(runtimeOpts, frameworkhttp.RetryOptions(opt))
	}
	return c.client.CallRequestWithPort(ctx, method, module, port, path, token, internalToken, body, out, runtimeOpts...)
}

func (c *HttpClient) CallRequest(ctx context.Context, method, module, path, token, internalToken string, body, out any, opts ...RetryOptions) error {
	runtimeOpts := make([]frameworkhttp.RetryOptions, 0, len(opts))
	for _, opt := range opts {
		runtimeOpts = append(runtimeOpts, frameworkhttp.RetryOptions(opt))
	}
	return c.client.CallRequest(ctx, method, module, path, token, internalToken, body, out, runtimeOpts...)
}

func (c *HttpClient) CallGet(ctx context.Context, module, path, token, internalToken string, out any, opts ...RetryOptions) error {
	return c.CallRequest(ctx, "GET", module, path, token, internalToken, nil, out, opts...)
}

func (c *HttpClient) CallPost(ctx context.Context, module, path, token, internalToken string, body, out any, opts ...RetryOptions) error {
	return c.CallRequest(ctx, "POST", module, path, token, internalToken, body, out, opts...)
}

func (c *HttpClient) CallPostAsync(ctx context.Context, module, path, token, internalToken string, body, out any, opts ...RetryOptions) {
	runtimeOpts := make([]frameworkhttp.RetryOptions, 0, len(opts))
	for _, opt := range opts {
		runtimeOpts = append(runtimeOpts, frameworkhttp.RetryOptions(opt))
	}
	c.client.CallPostAsync(ctx, module, path, token, internalToken, body, out, runtimeOpts...)
}
