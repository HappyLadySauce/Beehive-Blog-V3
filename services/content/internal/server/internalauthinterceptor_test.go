package server

import (
	"context"
	"errors"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	errgrpcx "github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs/grpcx"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestInternalAuthInterceptor(t *testing.T) {
	t.Parallel()

	interceptor := NewInternalAuthInterceptor(config.Config{
		InternalAuthToken: "secret",
		AllowedCallers:    []string{"gateway"},
	})
	info := &grpc.UnaryServerInfo{FullMethod: "/content.Content/Ping"}
	handler := func(ctx context.Context, _ any) (any, error) {
		caller, ok := ctxmeta.TrustedInternalCallerFrom(ctx)
		if !ok || caller != "gateway" {
			t.Fatalf("expected trusted internal caller marker, got %q ok=%v", caller, ok)
		}
		return "ok", nil
	}

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

	badCtx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		ctxmeta.MetadataKeyInternalAuthToken, "wrong",
		ctxmeta.MetadataKeyInternalCaller, "gateway",
	))
	_, err = interceptor.Unary()(badCtx, nil, info, handler)
	assertUnauthenticatedInternalCaller(t, err)
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
	if !ok || parsed == nil || !errors.Is(parsed, errs.E(errs.CodeContentInternalCallerUnauthorized)) {
		t.Fatalf("expected content internal caller unauthorized code, got %v", parsed)
	}
}
