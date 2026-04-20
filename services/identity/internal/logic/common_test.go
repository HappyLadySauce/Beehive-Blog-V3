package logic

import (
	"testing"

	identityservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestParseID verifies trusted identifier parsing behavior.
// TestParseID 验证可信标识解析行为。
func TestParseID(t *testing.T) {
	t.Parallel()

	if _, err := parseID("user_id", "123"); err != nil {
		t.Fatalf("expected parseID to succeed, got %v", err)
	}
	if _, err := parseID("user_id", "abc"); err == nil {
		t.Fatalf("expected parseID to fail on invalid input")
	}
}

// TestToStatusError verifies service-layer error mapping.
// TestToStatusError 验证 service 层错误映射。
func TestToStatusError(t *testing.T) {
	t.Parallel()

	err := toStatusError(identityservice.NewError(identityservice.ErrorKindFailedPrecondition, "account_disabled", nil), "fallback")
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected grpc status error, got %v", err)
	}
	if st.Code() != codes.FailedPrecondition {
		t.Fatalf("expected failed precondition, got %s", st.Code())
	}
	if st.Message() != "account_disabled" {
		t.Fatalf("expected account_disabled, got %s", st.Message())
	}
}
