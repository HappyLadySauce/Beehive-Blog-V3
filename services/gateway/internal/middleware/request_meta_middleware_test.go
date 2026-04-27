package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
)

func TestRequestMetaMiddlewareClientIPOrder(t *testing.T) {
	t.Parallel()

	m, err := NewRequestMetaMiddleware(config.GatewaySecurityConf{
		TrustedProxyCIDRs:   []string{"127.0.0.0/8"},
		TrustedProxyHeaders: []string{"X-Forwarded-For", "X-Real-IP", "Client-IP"},
	})
	if err != nil {
		t.Fatalf("expected middleware construction to succeed, got %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	req.Header.Set("X-Forwarded-For", "10.0.0.1, 10.0.0.2")
	req.Header.Set("X-Real-IP", "10.0.0.3")
	req.Header.Set("Client-IP", "10.0.0.4")
	rr := httptest.NewRecorder()

	var clientIP string
	m.Handle(func(_ http.ResponseWriter, r *http.Request) {
		meta, ok := RequestMetaFrom(r.Context())
		if !ok {
			t.Fatalf("request meta missing")
		}
		clientIP = meta.ClientIP
	})(rr, req)

	if clientIP != "10.0.0.1" {
		t.Fatalf("expected first forwarded IP, got %s", clientIP)
	}
}

func TestRequestMetaMiddlewareIgnoresForwardHeadersFromUntrustedSource(t *testing.T) {
	t.Parallel()

	m, err := NewRequestMetaMiddleware(config.GatewaySecurityConf{
		TrustedProxyCIDRs:   []string{"10.0.0.0/8"},
		TrustedProxyHeaders: []string{"X-Forwarded-For", "X-Real-IP", "Client-IP"},
	})
	if err != nil {
		t.Fatalf("expected middleware construction to succeed, got %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.RemoteAddr = "127.0.0.1:4567"
	req.Header.Set("X-Forwarded-For", "203.0.113.10")
	req.Header.Set("X-Real-IP", "203.0.113.11")
	req.Header.Set("Client-IP", "203.0.113.12")
	rr := httptest.NewRecorder()

	var clientIP string
	m.Handle(func(_ http.ResponseWriter, r *http.Request) {
		meta, ok := RequestMetaFrom(r.Context())
		if !ok {
			t.Fatalf("request meta missing")
		}
		clientIP = meta.ClientIP
	})(rr, req)

	if clientIP != "127.0.0.1" {
		t.Fatalf("expected socket ip for untrusted source, got %s", clientIP)
	}
}

func TestRequestMetaMiddlewareInvalidRemoteAddrFallsBackWithoutPanic(t *testing.T) {
	t.Parallel()

	m, err := NewRequestMetaMiddleware(config.GatewaySecurityConf{
		TrustedProxyCIDRs:   []string{"127.0.0.0/8"},
		TrustedProxyHeaders: []string{"X-Forwarded-For"},
	})
	if err != nil {
		t.Fatalf("expected middleware construction to succeed, got %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.RemoteAddr = "not-an-ip"
	req.Header.Set("X-Forwarded-For", "203.0.113.10")
	rr := httptest.NewRecorder()

	var clientIP string
	m.Handle(func(_ http.ResponseWriter, r *http.Request) {
		meta, ok := RequestMetaFrom(r.Context())
		if !ok {
			t.Fatalf("request meta missing")
		}
		clientIP = meta.ClientIP
	})(rr, req)

	if clientIP != "" {
		t.Fatalf("expected empty client ip for invalid remote addr, got %s", clientIP)
	}
}

func TestRequestMetaMiddlewareRejectsInvalidTrustedProxyCIDRs(t *testing.T) {
	t.Parallel()

	if _, err := NewRequestMetaMiddleware(config.GatewaySecurityConf{
		TrustedProxyCIDRs:   []string{"not-a-cidr"},
		TrustedProxyHeaders: []string{"X-Forwarded-For"},
	}); err == nil {
		t.Fatalf("expected invalid trusted proxy cidr error")
	}
}
