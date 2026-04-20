package errs_test

import (
	"errors"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
)

func TestNewAndParse(t *testing.T) {
	t.Parallel()

	err := errs.New(
		errs.CodeIdentityInvalidArgument,
		"invalid argument",
		errs.WithReference("ref-1"),
		errs.WithMeta(map[string]any{"field": "username"}),
	)

	parsed := errs.Parse(err)
	if parsed == nil {
		t.Fatalf("expected parsed error")
	}
	if parsed.Code != errs.CodeIdentityInvalidArgument {
		t.Fatalf("unexpected code: %d", parsed.Code)
	}
	if parsed.Reference != "ref-1" {
		t.Fatalf("unexpected reference: %s", parsed.Reference)
	}
	if parsed.Meta["field"] != "username" {
		t.Fatalf("unexpected meta: %+v", parsed.Meta)
	}
}

func TestWrapAndIsCode(t *testing.T) {
	t.Parallel()

	cause := errors.New("db failed")
	err := errs.Wrap(cause, errs.CodeIdentityDependencyUnavailable, "dependency unavailable")
	if !errors.Is(err, errs.E(errs.CodeIdentityDependencyUnavailable)) {
		t.Fatalf("expected code match")
	}

	parsed := errs.Parse(err)
	if parsed == nil || !errors.Is(parsed, cause) {
		t.Fatalf("expected wrapped cause")
	}
	if parsed.Error() != "dependency unavailable" {
		t.Fatalf("unexpected error string: %s", parsed.Error())
	}
}

func TestErrorsIsByBusinessCode(t *testing.T) {
	t.Parallel()

	err := errs.Wrap(errors.New("repo failed"), errs.CodeIdentityInvalidCredentials, "invalid credentials")
	if !errors.Is(err, errs.E(errs.CodeIdentityInvalidCredentials)) {
		t.Fatalf("expected business code match")
	}
	if code, ok := errs.CodeOf(err); !ok || code != errs.CodeIdentityInvalidCredentials {
		t.Fatalf("unexpected code extraction: code=%v ok=%v", code, ok)
	}
}

func TestErrorsJoinPreservesPrimaryBusinessCode(t *testing.T) {
	t.Parallel()

	joined := errors.Join(errors.New("redis unavailable"), errors.New("postgres unavailable"))
	err := errs.Wrap(joined, errs.CodeIdentityDependencyUnavailable, "identity dependencies are unavailable")
	if !errors.Is(err, errs.E(errs.CodeIdentityDependencyUnavailable)) {
		t.Fatalf("expected primary business code match")
	}
}
