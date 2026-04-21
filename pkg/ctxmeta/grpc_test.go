package ctxmeta

import (
	"context"
	"net/http/httptest"
	"testing"

	"google.golang.org/grpc/metadata"
)

// TestGetClientIPFromIncomingContext verifies only internal trusted metadata is accepted.
// TestGetClientIPFromIncomingContext 验证仅接受内部受控 metadata 作为可信客户端 IP。
func TestGetClientIPFromIncomingContext(t *testing.T) {
	t.Parallel()

	t.Run("ignores public forwarded headers", func(t *testing.T) {
		t.Parallel()

		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
			"x-forwarded-for", "203.0.113.10",
			"x-real-ip", "203.0.113.11",
			"client-ip", "203.0.113.12",
		))
		if clientIP := GetClientIPFromIncomingContext(ctx); clientIP != "" {
			t.Fatalf("expected forwarded headers to be ignored, got %s", clientIP)
		}
	})

	t.Run("reads internal trusted client ip", func(t *testing.T) {
		t.Parallel()

		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
			"x-beehive-trusted-client-ip", "198.51.100.7",
		))
		if clientIP := GetClientIPFromIncomingContext(ctx); clientIP != "198.51.100.7" {
			t.Fatalf("expected trusted client ip, got %s", clientIP)
		}
	})
}

// TestOutgoingContextWithRequestMeta verifies gateway only forwards internal trusted client ip.
// TestOutgoingContextWithRequestMeta 验证 gateway 仅转发内部受控客户端 IP。
func TestOutgoingContextWithRequestMeta(t *testing.T) {
	t.Parallel()

	ctx := OutgoingContextWithRequestMeta(context.Background(), RequestMeta{
		ForwardedFor: "203.0.113.10",
		RealIP:       "203.0.113.11",
		ClientIP:     "198.51.100.7",
		UserAgent:    "go-test",
		RequestID:    "req-1",
	})
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		t.Fatalf("expected outgoing metadata")
	}
	if got := md.Get("x-beehive-trusted-client-ip"); len(got) != 1 || got[0] != "198.51.100.7" {
		t.Fatalf("expected internal trusted client ip metadata, got %v", got)
	}
	if got := md.Get("x-forwarded-for"); len(got) != 0 {
		t.Fatalf("expected public forwarded headers to be omitted, got %v", got)
	}
}

// TestBuildRequestMetaFromHTTP verifies trusted proxy extraction still builds client ip correctly.
// TestBuildRequestMetaFromHTTP 验证受信代理提取仍能正确构建客户端 IP。
func TestBuildRequestMetaFromHTTP(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "https://example.com", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	req.Header.Set("X-Forwarded-For", "198.51.100.7, 198.51.100.8")

	networks, err := ParseTrustedProxyCIDRs([]string{"127.0.0.0/8"})
	if err != nil {
		t.Fatalf("expected cidr parse to succeed, got %v", err)
	}

	meta := BuildRequestMetaFromHTTP(req, TrustedProxyConfig{
		Headers: []string{"X-Forwarded-For"},
		CIDRs:   networks,
	})
	if meta.ClientIP != "198.51.100.7" {
		t.Fatalf("expected trusted forwarded client ip, got %s", meta.ClientIP)
	}
}
