package logic

import (
	"errors"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	errgrpcx "github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs/grpcx"
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

	err := toStatusError(errs.New(errs.CodeIdentityAccountDisabled, "account disabled"), "fallback")
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected grpc status error, got %v", err)
	}
	if st.Code() != codes.FailedPrecondition {
		t.Fatalf("expected failed precondition, got %s", st.Code())
	}
	if st.Message() != "account disabled" {
		t.Fatalf("expected account disabled, got %s", st.Message())
	}
	parsed, ok := errgrpcx.ParseStatus(err)
	if !ok {
		t.Fatalf("expected parsed grpc status detail")
	}
	if !errors.Is(parsed, errs.E(errs.CodeIdentityAccountDisabled)) {
		t.Fatalf("expected business code to survive grpc adaptation")
	}
}
