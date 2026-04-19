package auth

import "testing"

// TestIsProviderLoginReady verifies the phase-3 provider readiness contract.
// TestIsProviderLoginReady 验证第三阶段的 provider 就绪约束。
func TestIsProviderLoginReady(t *testing.T) {
	t.Parallel()

	if !IsProviderLoginReady(ProviderGitHub) {
		t.Fatalf("expected github provider to be ready")
	}
	if IsProviderLoginReady(ProviderQQ) {
		t.Fatalf("expected qq provider to be not ready")
	}
	if IsProviderLoginReady(ProviderWeChat) {
		t.Fatalf("expected wechat provider to be not ready")
	}
}
