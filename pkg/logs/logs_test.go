package logs

import (
	"context"
	"errors"
	"strings"
	"testing"

	berrs "github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
)

func TestBuildLineIncludesRequestIDAndCode(t *testing.T) {
	t.Parallel()

	ctx := WithRequestID(context.Background(), "req-1")
	line := buildLine(
		"auth_introspect",
		berrs.Wrap(errors.New("identity backend unavailable"), berrs.CodeGatewayAuthServiceUnavailable, "authentication service is unavailable"),
		fieldsFromContext(ctx),
		[]Field{String("route", "/api/v3/auth/me")},
	)

	if !strings.Contains(line, "action=auth_introspect") {
		t.Fatalf("expected action field, got %s", line)
	}
	if !strings.Contains(line, "request_id=req-1") {
		t.Fatalf("expected request_id field, got %s", line)
	}
	if !strings.Contains(line, "code=100601") {
		t.Fatalf("expected business code field, got %s", line)
	}
	if !strings.Contains(line, "cause=\"identity backend unavailable\"") {
		t.Fatalf("expected safe cause field, got %s", line)
	}
}

func TestBuildLineMasksSensitiveFields(t *testing.T) {
	t.Parallel()

	line := buildLine("login", nil, nil, []Field{
		String("password", "secret-value"),
		String("refresh_token", "abc"),
		String("route", "/api/v3/auth/login"),
	})

	if !strings.Contains(line, "password=[REDACTED]") {
		t.Fatalf("expected password redaction, got %s", line)
	}
	if !strings.Contains(line, "refresh_token=[REDACTED]") {
		t.Fatalf("expected refresh_token redaction, got %s", line)
	}
	if strings.Contains(line, "secret-value") || strings.Contains(line, "abc") {
		t.Fatalf("expected sensitive values to be removed, got %s", line)
	}
}

func TestErrorsContainingSensitiveKeywordsAreRedacted(t *testing.T) {
	t.Parallel()

	line := buildLine("token_parse", errors.New("authorization token leaked"), nil, nil)
	if !strings.Contains(line, "cause=[REDACTED]") {
		t.Fatalf("expected cause to be redacted, got %s", line)
	}
}
