package server

import (
	"context"
	"errors"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	errgrpcx "github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs/grpcx"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// TestInternalAuthInterceptor verifies all business RPCs require authenticated internal callers.
// TestInternalAuthInterceptor 验证所有业务 RPC 都要求已认证的内部调用方。
func TestInternalAuthInterceptor(t *testing.T) {
	t.Parallel()

	interceptor := NewInternalAuthInterceptor(config.Config{
		InternalAuthToken: "secret",
		AllowedCallers:    []string{"gateway"},
	})

	info := &grpc.UnaryServerInfo{FullMethod: "/identity.Identity/Ping"}
	handler := func(ctx context.Context, _ any) (any, error) {
		caller, ok := ctxmeta.TrustedInternalCallerFrom(ctx)
		if !ok || caller != "gateway" {
			t.Fatalf("expected trusted internal caller marker, got %q ok=%v", caller, ok)
		}
		return "ok", nil
	}

	t.Run("valid token and caller pass", func(t *testing.T) {
		t.Parallel()

		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
			ctxmeta.MetadataKeyInternalAuthToken, "secret",
			ctxmeta.MetadataKeyInternalCaller, "gateway",
		))
		resp, err := interceptor.Unary()(ctx, nil, info, handler)
		if err != nil {
			t.Fatalf("expected interceptor to allow request, got %v", err)
		}
		if resp != "ok" {
			t.Fatalf("unexpected response: %v", resp)
		}
	})

	t.Run("missing token fails", func(t *testing.T) {
		t.Parallel()

		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
			ctxmeta.MetadataKeyInternalCaller, "gateway",
		))
		_, err := interceptor.Unary()(ctx, nil, info, handler)
		assertUnauthenticatedInternalCaller(t, err)
	})

	t.Run("wrong token fails", func(t *testing.T) {
		t.Parallel()

		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
			ctxmeta.MetadataKeyInternalAuthToken, "wrong",
			ctxmeta.MetadataKeyInternalCaller, "gateway",
		))
		_, err := interceptor.Unary()(ctx, nil, info, handler)
		assertUnauthenticatedInternalCaller(t, err)
	})

	t.Run("disallowed caller fails", func(t *testing.T) {
		t.Parallel()

		ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
			ctxmeta.MetadataKeyInternalAuthToken, "secret",
			ctxmeta.MetadataKeyInternalCaller, "search",
		))
		_, err := interceptor.Unary()(ctx, nil, info, handler)
		assertUnauthenticatedInternalCaller(t, err)
	})

	t.Run("non identity method bypasses", func(t *testing.T) {
		t.Parallel()

		otherInfo := &grpc.UnaryServerInfo{FullMethod: "/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo"}
		called := false
		resp, err := interceptor.Unary()(context.Background(), nil, otherInfo, func(ctx context.Context, req any) (any, error) {
			called = true
			if _, ok := ctxmeta.TrustedInternalCallerFrom(ctx); ok {
				t.Fatalf("did not expect trusted internal caller marker on bypass path")
			}
			return "reflection", nil
		})
		if err != nil {
			t.Fatalf("expected bypass to succeed, got %v", err)
		}
		if !called || resp != "reflection" {
			t.Fatalf("expected bypass handler to run, called=%v resp=%v", called, resp)
		}
	})
}

func assertUnauthenticatedInternalCaller(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Fatalf("expected error")
	}
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected grpc status error, got %v", err)
	}
	if st.Code() != codes.Unauthenticated {
		t.Fatalf("expected unauthenticated, got %s", st.Code())
	}
	parsed, ok := errgrpcx.ParseStatus(err)
	if !ok {
		t.Fatalf("expected grpc status to carry business error details")
	}
	if parsed == nil || !errors.Is(parsed, errs.E(errs.CodeIdentityInternalCallerUnauthorized)) {
		t.Fatalf("expected internal caller unauthorized code, got %v", parsed)
	}
}
