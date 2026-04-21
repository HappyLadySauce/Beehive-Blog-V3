package config

import "testing"

// TestConfigValidate verifies gateway config validation behavior.
// TestConfigValidate 验证 gateway 配置校验行为。
func TestConfigValidate(t *testing.T) {
	t.Parallel()

	t.Run("valid trusted proxy cidrs", func(t *testing.T) {
		t.Parallel()

		conf := Config{
			Security: GatewaySecurityConf{
				TrustedProxyCIDRs: []string{"127.0.0.0/8", "10.0.0.0/8"},
			},
		}
		if err := conf.Validate(); err != nil {
			t.Fatalf("expected config validation to succeed, got %v", err)
		}
	})

	t.Run("invalid trusted proxy cidrs", func(t *testing.T) {
		t.Parallel()

		conf := Config{
			Security: GatewaySecurityConf{
				TrustedProxyCIDRs: []string{"bad-cidr"},
			},
		}
		if err := conf.Validate(); err == nil {
			t.Fatalf("expected config validation to fail")
		}
	})
}
