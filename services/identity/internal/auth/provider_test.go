package auth

import "testing"

// TestEnsureState verifies state generation keeps explicit input and fills blanks.
// TestEnsureState 验证 state 生成会保留显式输入并补齐空值。
func TestEnsureState(t *testing.T) {
	t.Parallel()

	if got := EnsureState("fixed-state"); got != "fixed-state" {
		t.Fatalf("expected explicit state to be preserved, got %q", got)
	}

	if got := EnsureState(""); got == "" {
		t.Fatalf("expected generated state to be non-empty")
	}
}
