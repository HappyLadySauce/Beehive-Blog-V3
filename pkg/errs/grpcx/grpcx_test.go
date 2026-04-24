package grpcx_test

import (
	"errors"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	errgrpcx "github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs/grpcx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestToStatusAndParseStatus(t *testing.T) {
	t.Parallel()

	grpcErr := errgrpcx.ToStatus(
		errs.New(errs.CodeIdentityEmailAlreadyExists, "email already exists", errs.WithReference("ref-1")),
		"fallback",
	)

	st, ok := status.FromError(grpcErr)
	if !ok {
		t.Fatalf("expected grpc status")
	}
	if st.Code() != codes.AlreadyExists {
		t.Fatalf("unexpected grpc code: %s", st.Code())
	}

	parsed, ok := errgrpcx.ParseStatus(grpcErr)
	if !ok {
		t.Fatalf("expected parse success")
	}
	if parsed.Code != errs.CodeIdentityEmailAlreadyExists {
		t.Fatalf("unexpected business code: %d", parsed.Code)
	}
	if parsed.Reference != "ref-1" {
		t.Fatalf("unexpected reference: %s", parsed.Reference)
	}
	if !errors.Is(parsed, errs.E(errs.CodeIdentityEmailAlreadyExists)) {
		t.Fatalf("expected errors.Is to match business code")
	}
}

func TestContentCodeToGRPC(t *testing.T) {
	t.Parallel()

	cases := map[errs.Code]codes.Code{
		errs.CodeContentInternalCallerUnauthorized: codes.Unauthenticated,
		errs.CodeContentAccessForbidden:            codes.PermissionDenied,
		errs.CodeContentNotFound:                   codes.NotFound,
		errs.CodeContentTagNotFound:                codes.NotFound,
		errs.CodeContentRevisionNotFound:           codes.NotFound,
		errs.CodeContentRelationNotFound:           codes.NotFound,
		errs.CodeContentSlugAlreadyExists:          codes.AlreadyExists,
		errs.CodeContentTagAlreadyExists:           codes.AlreadyExists,
		errs.CodeContentRelationAlreadyExists:      codes.AlreadyExists,
		errs.CodeContentTagInUse:                   codes.FailedPrecondition,
	}
	for code, want := range cases {
		if got := errgrpcx.CodeToGRPC(code); got != want {
			t.Fatalf("CodeToGRPC(%d) = %s, want %s", code, got, want)
		}
	}
}
