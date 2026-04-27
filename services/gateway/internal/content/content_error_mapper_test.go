package content

import (
	"context"
	"errors"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	errgrpcx "github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs/grpcx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMapUpstreamErrorPreservesContentDomainCode(t *testing.T) {
	t.Parallel()

	err := errgrpcx.ToStatus(errs.New(errs.CodeContentSlugAlreadyExists, "slug already exists"), "fallback")
	mapped := MapUpstreamError(context.Background(), "content_create", "/route", err)
	if !errors.Is(mapped, errs.E(errs.CodeContentSlugAlreadyExists)) {
		t.Fatalf("expected slug conflict to be preserved, got %v", mapped)
	}
}

func TestMapUpstreamErrorMasksContentInternalCallerUnauthorized(t *testing.T) {
	t.Parallel()

	err := errgrpcx.ToStatus(errs.New(errs.CodeContentInternalCallerUnauthorized, "internal caller authentication failed"), "fallback")
	mapped := MapUpstreamError(context.Background(), "content_create", "/route", err)
	if !errors.Is(mapped, errs.E(errs.CodeGatewayNotReady)) {
		t.Fatalf("expected gateway not ready error, got %v", mapped)
	}
	if errors.Is(mapped, errs.E(errs.CodeContentInternalCallerUnauthorized)) {
		t.Fatalf("expected internal auth code to be masked, got %v", mapped)
	}
}

func TestMapUpstreamErrorMapsContentUnavailable(t *testing.T) {
	t.Parallel()

	mapped := MapUpstreamError(context.Background(), "content_list", "/route", status.Error(codes.Unavailable, "unavailable"))
	if !errors.Is(mapped, errs.E(errs.CodeGatewayContentServiceUnavailable)) {
		t.Fatalf("expected content service unavailable, got %v", mapped)
	}
	if errors.Is(mapped, errs.E(errs.CodeGatewayAuthServiceUnavailable)) {
		t.Fatalf("expected content unavailable not auth unavailable, got %v", mapped)
	}
}

func TestMapUpstreamErrorPreservesTagInUse(t *testing.T) {
	t.Parallel()

	err := errgrpcx.ToStatus(errs.New(errs.CodeContentTagInUse, "tag is in use"), "fallback")
	mapped := MapUpstreamError(context.Background(), "content_tag_delete", "/route", err)
	if !errors.Is(mapped, errs.E(errs.CodeContentTagInUse)) {
		t.Fatalf("expected tag in use error to be preserved, got %v", mapped)
	}
}

func TestMapUpstreamErrorPreservesRelationConflict(t *testing.T) {
	t.Parallel()

	err := errgrpcx.ToStatus(errs.New(errs.CodeContentRelationAlreadyExists, "content relation already exists"), "fallback")
	mapped := MapUpstreamError(context.Background(), "content_relation_create", "/route", err)
	if !errors.Is(mapped, errs.E(errs.CodeContentRelationAlreadyExists)) {
		t.Fatalf("expected relation conflict to be preserved, got %v", mapped)
	}
}

func TestMapUpstreamErrorMapsAlreadyExistsFallbackByRoute(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		action string
		route  string
		code   errs.Code
	}{
		{
			name:   "relation conflict",
			action: "content_relation_create",
			route:  "/api/v3/studio/content/items/:content_id/relations",
			code:   errs.CodeContentRelationAlreadyExists,
		},
		{
			name:   "tag conflict",
			action: "content_tag_create",
			route:  "/api/v3/studio/content/tags",
			code:   errs.CodeContentTagAlreadyExists,
		},
		{
			name:   "slug conflict",
			action: "content_create",
			route:  "/api/v3/studio/content/items",
			code:   errs.CodeContentSlugAlreadyExists,
		},
		{
			name:   "unknown stage route does not match tag substring",
			action: "content_stage_create",
			route:  "/api/v3/studio/content/stages",
			code:   errs.CodeContentSlugAlreadyExists,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mapped := MapUpstreamError(context.Background(), tc.action, tc.route, status.Error(codes.AlreadyExists, "already exists"))
			if !errors.Is(mapped, errs.E(tc.code)) {
				t.Fatalf("expected %d, got %v", tc.code, mapped)
			}
		})
	}
}
