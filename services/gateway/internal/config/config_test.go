package config

import "testing"

// TestConfigValidate verifies gateway config validation behavior.
// TestConfigValidate 验证 gateway 配置校验行为。
func TestConfigValidate(t *testing.T) {
	t.Parallel()

	t.Run("valid trusted proxy cidrs", func(t *testing.T) {
		t.Parallel()

		conf := Config{
			IdentityRPC: IdentityRPCConf{
				InternalAuthToken:  "secret",
				InternalCallerName: "gateway",
			},
			Security: GatewaySecurityConf{
				TrustedProxyHeaders: []string{"X-Forwarded-For", "X-Real-IP"},
				TrustedProxyCIDRs:   []string{"127.0.0.0/8", "10.0.0.0/8"},
			},
		}
		if err := conf.Validate(); err != nil {
			t.Fatalf("expected config validation to succeed, got %v", err)
		}
	})

	t.Run("trusted proxy headers are case insensitive", func(t *testing.T) {
		t.Parallel()

		conf := Config{
			IdentityRPC: IdentityRPCConf{
				InternalAuthToken:  "secret",
				InternalCallerName: "gateway",
			},
			Security: GatewaySecurityConf{
				TrustedProxyHeaders: []string{"x-forwarded-for", "CLIENT-IP"},
				TrustedProxyCIDRs:   []string{"127.0.0.0/8"},
			},
		}
		if err := conf.Validate(); err != nil {
			t.Fatalf("expected case-insensitive trusted proxy headers to pass, got %v", err)
		}
	})

	t.Run("invalid trusted proxy cidrs", func(t *testing.T) {
		t.Parallel()

		conf := Config{
			IdentityRPC: IdentityRPCConf{
				InternalAuthToken:  "secret",
				InternalCallerName: "gateway",
			},
			Security: GatewaySecurityConf{
				TrustedProxyHeaders: []string{"X-Forwarded-For"},
				TrustedProxyCIDRs:   []string{"bad-cidr"},
			},
		}
		if err := conf.Validate(); err == nil {
			t.Fatalf("expected config validation to fail")
		}
	})

	t.Run("unsupported trusted proxy header should fail", func(t *testing.T) {
		t.Parallel()

		conf := Config{
			IdentityRPC: IdentityRPCConf{
				InternalAuthToken:  "secret",
				InternalCallerName: "gateway",
			},
			Security: GatewaySecurityConf{
				TrustedProxyHeaders: []string{"X-Forwaded-For"},
				TrustedProxyCIDRs:   []string{"127.0.0.0/8"},
			},
		}
		if err := conf.Validate(); err == nil {
			t.Fatalf("expected config validation to fail for unsupported trusted proxy header")
		}
	})

	t.Run("headers without cidrs should fail", func(t *testing.T) {
		t.Parallel()

		conf := Config{
			IdentityRPC: IdentityRPCConf{
				InternalAuthToken:  "secret",
				InternalCallerName: "gateway",
			},
			Security: GatewaySecurityConf{
				TrustedProxyHeaders: []string{"X-Forwarded-For"},
				TrustedProxyCIDRs:   []string{},
			},
		}
		if err := conf.Validate(); err == nil {
			t.Fatalf("expected config validation to fail when headers configured without cidrs")
		}
	})

	t.Run("cidrs without headers should fail", func(t *testing.T) {
		t.Parallel()

		conf := Config{
			IdentityRPC: IdentityRPCConf{
				InternalAuthToken:  "secret",
				InternalCallerName: "gateway",
			},
			Security: GatewaySecurityConf{
				TrustedProxyCIDRs: []string{"127.0.0.0/8"},
			},
		}
		if err := conf.Validate(); err == nil {
			t.Fatalf("expected config validation to fail when cidrs configured without headers")
		}
	})

	t.Run("empty headers and cidrs should pass", func(t *testing.T) {
		t.Parallel()

		conf := Config{
			IdentityRPC: IdentityRPCConf{
				InternalAuthToken:  "secret",
				InternalCallerName: "gateway",
			},
			Security: GatewaySecurityConf{
				TrustedProxyHeaders: []string{},
				TrustedProxyCIDRs:   []string{},
			},
		}
		if err := conf.Validate(); err != nil {
			t.Fatalf("expected config validation to succeed for socket-ip mode, got %v", err)
		}
	})

	t.Run("missing internal auth token should fail", func(t *testing.T) {
		t.Parallel()

		conf := Config{
			IdentityRPC: IdentityRPCConf{
				InternalCallerName: "gateway",
			},
		}
		if err := conf.Validate(); err == nil {
			t.Fatalf("expected config validation to fail when internal auth token is missing")
		}
	})

	t.Run("missing internal caller name should fail", func(t *testing.T) {
		t.Parallel()

		conf := Config{
			IdentityRPC: IdentityRPCConf{
				InternalAuthToken: "secret",
			},
		}
		if err := conf.Validate(); err == nil {
			t.Fatalf("expected config validation to fail when internal caller name is missing")
		}
	})
}
