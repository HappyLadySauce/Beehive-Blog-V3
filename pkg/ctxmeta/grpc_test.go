package ctxmeta

import (
	"context"
	"net/http/httptest"
	"testing"

	"google.golang.org/grpc/metadata"
)

// TestGetClientIPFromIncomingContext verifies trusted client IP requires authenticated internal caller context.
// TestGetClientIPFromIncomingContext 验证可信客户端 IP 必须依赖已认证的内部调用方上下文。
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

	t.Run("requires trusted caller marker", func(t *testing.T) {
		t.Parallel()

		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
			MetadataKeyTrustedClientIP, "198.51.100.7",
		))
		if clientIP := GetClientIPFromIncomingContext(ctx); clientIP != "" {
			t.Fatalf("expected empty client ip without trusted caller marker, got %s", clientIP)
		}

		ctx = WithTrustedInternalCaller(ctx, "gateway")
		if clientIP := GetClientIPFromIncomingContext(ctx); clientIP != "198.51.100.7" {
			t.Fatalf("expected trusted client ip, got %s", clientIP)
		}
	})
}

// TestBuildIdentityOutgoingContext verifies gateway forwards internal auth and trusted metadata.
// TestBuildIdentityOutgoingContext 验证 gateway 会转发内部认证与可信元数据。
func TestBuildIdentityOutgoingContext(t *testing.T) {
	t.Parallel()

	ctx := BuildIdentityOutgoingContext(context.Background(), RequestMeta{
		ForwardedFor: "203.0.113.10",
		RealIP:       "203.0.113.11",
		ClientIP:     "198.51.100.7",
		UserAgent:    "go-test",
		RequestID:    "req-1",
	}, InternalRPCAuth{
		Token:  "secret-token",
		Caller: "gateway",
	})
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		t.Fatalf("expected outgoing metadata")
	}
	if got := md.Get(MetadataKeyInternalAuthToken); len(got) != 1 || got[0] != "secret-token" {
		t.Fatalf("expected internal auth token metadata, got %v", got)
	}
	if got := md.Get(MetadataKeyInternalCaller); len(got) != 1 || got[0] != "gateway" {
		t.Fatalf("expected internal caller metadata, got %v", got)
	}
	if got := md.Get(MetadataKeyTrustedClientIP); len(got) != 1 || got[0] != "198.51.100.7" {
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
