package httpx_test

import (
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	errhttpx "github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs/httpx"
)

func TestBuildResponse(t *testing.T) {
	t.Parallel()

	statusCode, payload := errhttpx.BuildResponse(
		errs.New(errs.CodeGatewayAccessTokenInvalid, "access token is invalid", errs.WithReference("ref-1")),
		"req-1",
	)

	if statusCode != 401 {
		t.Fatalf("expected 401, got %d", statusCode)
	}
	if payload.Code != errs.CodeGatewayAccessTokenInvalid {
		t.Fatalf("unexpected code: %d", payload.Code)
	}
	if payload.Reference != "ref-1" || payload.RequestID != "req-1" {
		t.Fatalf("unexpected payload: %+v", payload)
	}
}

func TestContentCodeToHTTP(t *testing.T) {
	t.Parallel()

	cases := map[errs.Code]int{
		errs.CodeContentInternalCallerUnauthorized: 401,
		errs.CodeContentAccessForbidden:            403,
		errs.CodeContentNotFound:                   404,
		errs.CodeContentTagNotFound:                404,
		errs.CodeContentRevisionNotFound:           404,
		errs.CodeContentSlugAlreadyExists:          409,
		errs.CodeContentTagAlreadyExists:           409,
		errs.CodeContentTagInUse:                   412,
	}
	for code, want := range cases {
		if got := errhttpx.CodeToHTTP(code); got != want {
			t.Fatalf("CodeToHTTP(%d) = %d, want %d", code, got, want)
		}
	}
}
