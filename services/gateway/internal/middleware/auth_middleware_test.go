package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestAuthMiddlewareRejectsMissingToken(t *testing.T) {
	t.Parallel()

	m := NewAuthMiddleware(&fakeIdentityClient{
		introspectFn: func(_ context.Context, _ *pb.IntrospectAccessTokenRequest, _ ...grpc.CallOption) (*pb.IntrospectAccessTokenResponse, error) {
			return &pb.IntrospectAccessTokenResponse{Active: true}, nil
		},
	}, config.GatewaySecurityConf{TokenPrefix: "Bearer"}, config.IdentityRPCConf{InternalAuthToken: "secret", InternalCallerName: "gateway"})
	req := httptest.NewRequest(http.MethodGet, "/api/v3/auth/me", nil)
	rr := httptest.NewRecorder()

	m.Handle(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatalf("should not pass")
	})(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}

	var payload errorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unexpected json error: %v", err)
	}
	if payload.Message != "authorization is required" {
		t.Fatalf("unexpected message: %+v", payload)
	}
	if payload.Code != errs.CodeGatewayAuthorizationRequired {
		t.Fatalf("unexpected code: %+v", payload)
	}
}

func TestAuthMiddlewareInactiveToken(t *testing.T) {
	t.Parallel()

	m := NewAuthMiddleware(&fakeIdentityClient{
		introspectFn: func(_ context.Context, _ *pb.IntrospectAccessTokenRequest, _ ...grpc.CallOption) (*pb.IntrospectAccessTokenResponse, error) {
			return &pb.IntrospectAccessTokenResponse{Active: false}, nil
		},
	}, config.GatewaySecurityConf{TokenPrefix: "Bearer"}, config.IdentityRPCConf{InternalAuthToken: "secret", InternalCallerName: "gateway"})
	req := httptest.NewRequest(http.MethodGet, "/api/v3/auth/me", nil)
	req.Header.Set("Authorization", "Bearer token")
	rr := httptest.NewRecorder()

	m.Handle(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatalf("should not pass")
	})(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}

	var payload errorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unexpected json error: %v", err)
	}
	if payload.Message != "access token is inactive" {
		t.Fatalf("unexpected message: %+v", payload)
	}
	if payload.Code != errs.CodeGatewayAccessTokenInactive {
		t.Fatalf("unexpected code: %+v", payload)
	}
}

func TestAuthMiddlewareWritesTrustedContext(t *testing.T) {
	t.Parallel()

	m := NewAuthMiddleware(&fakeIdentityClient{
		introspectFn: func(ctx context.Context, _ *pb.IntrospectAccessTokenRequest, _ ...grpc.CallOption) (*pb.IntrospectAccessTokenResponse, error) {
			md, ok := metadata.FromOutgoingContext(ctx)
			if !ok {
				t.Fatalf("expected outgoing metadata")
			}
			if got := md.Get(ctxmeta.MetadataKeyInternalAuthToken); len(got) != 1 || got[0] != "secret" {
				t.Fatalf("expected internal auth token metadata, got %v", got)
			}
			if got := md.Get(ctxmeta.MetadataKeyInternalCaller); len(got) != 1 || got[0] != "gateway" {
				t.Fatalf("expected internal caller metadata, got %v", got)
			}
			return &pb.IntrospectAccessTokenResponse{
				Active:        true,
				UserId:        "u1",
				SessionId:     "s1",
				Role:          pb.Role_ROLE_ADMIN,
				AccountStatus: pb.AccountStatus_ACCOUNT_STATUS_ACTIVE,
				AuthSource:    pb.AuthSource_AUTH_SOURCE_LOCAL,
			}, nil
		},
	}, config.GatewaySecurityConf{TokenPrefix: "Bearer"}, config.IdentityRPCConf{InternalAuthToken: "secret", InternalCallerName: "gateway"})
	req := httptest.NewRequest(http.MethodGet, "/api/v3/auth/me", nil)
	req.Header.Set("Authorization", "Bearer token")
	rr := httptest.NewRecorder()

	called := false
	m.Handle(func(_ http.ResponseWriter, r *http.Request) {
		called = true
		authCtx, ok := AuthContextFrom(r.Context())
		if !ok || authCtx.UserID != "u1" || authCtx.SessionID != "s1" {
			t.Fatalf("unexpected auth context: %+v, ok=%v", authCtx, ok)
		}
	})(rr, req)

	if !called {
		t.Fatalf("expected handler to run")
	}
}

func TestAuthMiddlewareMasksUpstreamErrors(t *testing.T) {
	t.Parallel()

	m := NewAuthMiddleware(&fakeIdentityClient{
		introspectFn: func(_ context.Context, _ *pb.IntrospectAccessTokenRequest, _ ...grpc.CallOption) (*pb.IntrospectAccessTokenResponse, error) {
			return nil, status.Error(codes.Unavailable, "identity backend exploded")
		},
	}, config.GatewaySecurityConf{TokenPrefix: "Bearer"}, config.IdentityRPCConf{InternalAuthToken: "secret", InternalCallerName: "gateway"})
	req := httptest.NewRequest(http.MethodGet, "/api/v3/auth/me", nil)
	req.Header.Set("Authorization", "Bearer token")
	rr := httptest.NewRecorder()

	m.Handle(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatalf("should not pass")
	})(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rr.Code)
	}

	var payload errorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unexpected json error: %v", err)
	}
	if payload.Message != "authentication service is unavailable" {
		t.Fatalf("unexpected payload: %+v", payload)
	}
	if payload.Code != errs.CodeGatewayAuthServiceUnavailable {
		t.Fatalf("unexpected code: %+v", payload)
	}
	if strings.Contains(rr.Body.String(), "identity backend exploded") {
		t.Fatalf("response leaked upstream error: %s", rr.Body.String())
	}
}

func TestAuthMiddlewarePermissionDeniedMapsForbidden(t *testing.T) {
	t.Parallel()

	m := NewAuthMiddleware(&fakeIdentityClient{
		introspectFn: func(_ context.Context, _ *pb.IntrospectAccessTokenRequest, _ ...grpc.CallOption) (*pb.IntrospectAccessTokenResponse, error) {
			return nil, status.Error(codes.PermissionDenied, "forbidden upstream")
		},
	}, config.GatewaySecurityConf{TokenPrefix: "Bearer"}, config.IdentityRPCConf{InternalAuthToken: "secret", InternalCallerName: "gateway"})
	req := httptest.NewRequest(http.MethodGet, "/api/v3/auth/me", nil)
	req.Header.Set("Authorization", "Bearer token")
	rr := httptest.NewRecorder()

	m.Handle(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatalf("should not pass")
	})(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rr.Code)
	}

	var payload errorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unexpected json error: %v", err)
	}
	if payload.Code != errs.CodeGatewayAccessForbidden {
		t.Fatalf("unexpected code: %+v", payload)
	}
}
