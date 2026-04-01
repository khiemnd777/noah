package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	fiberws "github.com/gofiber/websocket/v2"
	"github.com/khiemnd777/noah_framework/shared/logger"
	frameworkapp "github.com/khiemnd777/noah_framework/pkg/app"
	frameworkhttp "github.com/khiemnd777/noah_framework/pkg/http"
	frameworkruntime "github.com/khiemnd777/noah_framework/runtime"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

type LoadBalancer struct {
	targets []*url.URL
	alive   []bool
	counter uint64
}

func NewLoadBalancer(targets []string) (*LoadBalancer, error) {
	var urls []*url.URL
	for _, t := range targets {
		u, err := url.Parse(t)
		if err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}
	alive := make([]bool, len(urls))
	for i := range alive {
		alive[i] = true
	}
	return &LoadBalancer{targets: urls, alive: alive}, nil
}

func (lb *LoadBalancer) NextTarget() *url.URL {
	for i := 0; i < len(lb.targets); i++ {
		index := (lb.counter + uint64(i)) % uint64(len(lb.targets))
		if lb.alive[index] {
			lb.counter = index + 1
			return lb.targets[index]
		}
	}
	return lb.targets[0]
}

func singleJoiningSlash(a, b string) string {
	aSlash := strings.HasSuffix(a, "/")
	bSlash := strings.HasPrefix(b, "/")
	switch {
	case aSlash && bSlash:
		return a + b[1:]
	case !aSlash && !bSlash:
		return a + "/" + b
	default:
		return a + b
	}
}

func isWebSocketRequest(c frameworkhttp.Context) bool {
	// RFC 6455
	if c.Method() != fiber.MethodGet {
		return false
	}
	if strings.ToLower(c.Get("Upgrade")) != "websocket" {
		return false
	}
	if !strings.Contains(strings.ToLower(c.Get("Connection")), "upgrade") {
		return false
	}
	return true
}

// RegisterReverseProxy mounts a reverse proxy at given route with load balancing
func RegisterReverseProxy(app frameworkapp.Application, route string, targets []string) error {
	lb, err := NewLoadBalancer(targets)
	if err != nil {
		return err
	}

	app.Router().All(route+"/*", func(c frameworkhttp.Context) error {
		fiberCtx := frameworkruntime.MustFiberContext(c)
		target := lb.NextTarget()

		// ✅ WS: bypass circuit breaker + use WS bridge proxy
		if isWebSocketRequest(c) {
			// lưu context cần thiết cho ws handler
			fiberCtx.Locals("__proxy_target", target.String())
			fiberCtx.Locals("__proxy_path", fiberCtx.Params("*"))
			fiberCtx.Locals("__proxy_query", string(fiberCtx.Context().URI().QueryString()))
			fiberCtx.Locals("__proxy_auth", fiberCtx.Get("Authorization"))

			logger.Debug(fmt.Sprintf("[Gateway] WS Proxy %s → %s", fiberCtx.OriginalURL(), target.String()))

			// IMPORTANT: MUST upgrade using fiberws.New
			return fiberws.New(proxyWebSocket)(fiberCtx)
		}

		// ✅ HTTP: keep current reverse proxy (can be wrapped by circuit breaker outside if you want)
		proxy := httputil.NewSingleHostReverseProxy(target)

		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.URL.Path = singleJoiningSlash(target.Path, fiberCtx.Params("*"))
			req.URL.RawQuery = string(fiberCtx.Context().URI().QueryString())
			req.Host = target.Host

			// Clean hop-by-hop headers (HTTP only)
			hopHeaders := []string{
				"Connection", "Keep-Alive", "Proxy-Authenticate",
				"Proxy-Authorization", "Te", "Trailer",
				"Transfer-Encoding", "Upgrade",
			}
			for _, h := range hopHeaders {
				delete(req.Header, h)
			}

			// Forward headers
			if auth := fiberCtx.Get("Authorization"); auth != "" {
				req.Header.Set("Authorization", auth)
			}
			req.Header.Set("X-Internal-Token", frameworkruntime.InternalAuthToken())

			logger.Debug(fmt.Sprintf("[Gateway] Proxy %s → %s", fiberCtx.OriginalURL(), target.String()))
		}

		fasthttpadaptor.NewFastHTTPHandler(proxy)(fiberCtx.Context())
		return nil
	})

	return nil
}
