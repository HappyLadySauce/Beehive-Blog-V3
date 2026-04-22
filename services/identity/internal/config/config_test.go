package config

import "testing"

// TestConfigValidate verifies identity configuration validation behavior.
// TestConfigValidate 验证 identity 配置校验行为。
func TestConfigValidate(t *testing.T) {
	t.Parallel()

	base := Config{
		InternalAuthToken: "secret",
		AllowedCallers:    []string{"gateway"},
		Postgres: PostgresConf{
			Host:   "127.0.0.1",
			User:   "user",
			DBName: "db",
		},
		StateRedis: RedisConf{
			Host: "127.0.0.1",
		},
		Security: SecurityConf{
			AccessTokenSecret:      "test-secret",
			AccessTokenTTLSeconds:  900,
			RefreshTokenTTLSeconds: 3600,
		},
	}

	t.Run("valid config", func(t *testing.T) {
		t.Parallel()

		if err := base.Validate(); err != nil {
			t.Fatalf("expected config validation to succeed, got %v", err)
		}
	})

	t.Run("missing internal auth token", func(t *testing.T) {
		t.Parallel()

		conf := base
		conf.InternalAuthToken = ""
		if err := conf.Validate(); err == nil {
			t.Fatalf("expected config validation to fail")
		}
	})

	t.Run("missing allowed callers", func(t *testing.T) {
		t.Parallel()

		conf := base
		conf.AllowedCallers = nil
		if err := conf.Validate(); err == nil {
			t.Fatalf("expected config validation to fail")
		}
	})

	t.Run("empty allowed caller value", func(t *testing.T) {
		t.Parallel()

		conf := base
		conf.AllowedCallers = []string{"gateway", " "}
		if err := conf.Validate(); err == nil {
			t.Fatalf("expected config validation to fail")
		}
	})

	t.Run("enabled qq requires complete configuration", func(t *testing.T) {
		t.Parallel()

		conf := base
		conf.SSO.QQ.Enabled = true
		if err := conf.Validate(); err == nil {
			t.Fatalf("expected config validation to fail")
		}
	})

	t.Run("enabled wechat requires complete configuration", func(t *testing.T) {
		t.Parallel()

		conf := base
		conf.SSO.WeChat.Enabled = true
		conf.SSO.WeChat.ClientID = "wechat-appid"
		conf.SSO.WeChat.ClientSecret = "wechat-secret"
		conf.SSO.WeChat.RedirectURL = "https://example.com/auth/wechat/callback"
		if err := conf.Validate(); err != nil {
			t.Fatalf("expected wechat config validation to succeed, got %v", err)
		}
	})

	t.Run("wechat defaults to snsapi_login when scopes empty", func(t *testing.T) {
		t.Parallel()

		conf := base
		conf.SSO.WeChat.Enabled = true
		conf.SSO.WeChat.ClientID = "wechat-appid"
		conf.SSO.WeChat.ClientSecret = "wechat-secret"
		conf.SSO.WeChat.RedirectURL = "https://example.com/auth/wechat/callback"
		conf.SSO.WeChat.Scopes = nil
		if err := conf.Validate(); err != nil {
			t.Fatalf("expected empty wechat scopes to be accepted, got %v", err)
		}
	})

	t.Run("wechat rejects unsupported scope", func(t *testing.T) {
		t.Parallel()

		conf := base
		conf.SSO.WeChat.Enabled = true
		conf.SSO.WeChat.ClientID = "wechat-appid"
		conf.SSO.WeChat.ClientSecret = "wechat-secret"
		conf.SSO.WeChat.RedirectURL = "https://example.com/auth/wechat/callback"
		conf.SSO.WeChat.Scopes = []string{"snsapi_base"}
		if err := conf.Validate(); err == nil {
			t.Fatalf("expected unsupported wechat scope to fail validation")
		}
	})

	t.Run("wechat rejects multiple scopes", func(t *testing.T) {
		t.Parallel()

		conf := base
		conf.SSO.WeChat.Enabled = true
		conf.SSO.WeChat.ClientID = "wechat-appid"
		conf.SSO.WeChat.ClientSecret = "wechat-secret"
		conf.SSO.WeChat.RedirectURL = "https://example.com/auth/wechat/callback"
		conf.SSO.WeChat.Scopes = []string{"snsapi_login", "snsapi_base"}
		if err := conf.Validate(); err == nil {
			t.Fatalf("expected multiple wechat scopes to fail validation")
		}
	})
}
