package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type fakeIdentityClient struct {
	introspectFn func(ctx context.Context, in *pb.IntrospectAccessTokenRequest, opts ...grpc.CallOption) (*pb.IntrospectAccessTokenResponse, error)
}

func (f *fakeIdentityClient) RegisterLocalUser(context.Context, *pb.RegisterLocalUserRequest, ...grpc.CallOption) (*pb.RegisterLocalUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
func (f *fakeIdentityClient) LoginLocalUser(context.Context, *pb.LoginLocalUserRequest, ...grpc.CallOption) (*pb.LoginLocalUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
func (f *fakeIdentityClient) StartSsoLogin(context.Context, *pb.StartSsoLoginRequest, ...grpc.CallOption) (*pb.StartSsoLoginResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
func (f *fakeIdentityClient) FinishSsoLogin(context.Context, *pb.FinishSsoLoginRequest, ...grpc.CallOption) (*pb.FinishSsoLoginResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
func (f *fakeIdentityClient) RefreshSessionToken(context.Context, *pb.RefreshSessionTokenRequest, ...grpc.CallOption) (*pb.RefreshSessionTokenResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
func (f *fakeIdentityClient) LogoutSession(context.Context, *pb.LogoutSessionRequest, ...grpc.CallOption) (*pb.LogoutSessionResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
func (f *fakeIdentityClient) GetCurrentUser(context.Context, *pb.GetCurrentUserRequest, ...grpc.CallOption) (*pb.GetCurrentUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
func (f *fakeIdentityClient) IntrospectAccessToken(ctx context.Context, in *pb.IntrospectAccessTokenRequest, opts ...grpc.CallOption) (*pb.IntrospectAccessTokenResponse, error) {
	return f.introspectFn(ctx, in, opts...)
}
func (f *fakeIdentityClient) Ping(context.Context, *pb.PingRequest, ...grpc.CallOption) (*pb.PingResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

type errorResponse struct {
	Code      errs.Code `json:"code"`
	Message   string    `json:"message"`
	Reference string    `json:"reference"`
	RequestID string    `json:"request_id"`
}

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
			if got := md.Get("x-beehive-internal-auth-token"); len(got) != 1 || got[0] != "secret" {
				t.Fatalf("expected internal auth token metadata, got %v", got)
			}
			if got := md.Get("x-beehive-internal-caller"); len(got) != 1 || got[0] != "gateway" {
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

func TestRequestMetaMiddlewareClientIPOrder(t *testing.T) {
	t.Parallel()

	m, err := NewRequestMetaMiddleware(config.GatewaySecurityConf{
		TrustedProxyCIDRs:   []string{"127.0.0.0/8"},
		TrustedProxyHeaders: []string{"X-Forwarded-For", "X-Real-IP", "Client-IP"},
	})
	if err != nil {
		t.Fatalf("expected middleware construction to succeed, got %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	req.Header.Set("X-Forwarded-For", "10.0.0.1, 10.0.0.2")
	req.Header.Set("X-Real-IP", "10.0.0.3")
	req.Header.Set("Client-IP", "10.0.0.4")
	rr := httptest.NewRecorder()

	var clientIP string
	m.Handle(func(_ http.ResponseWriter, r *http.Request) {
		meta, ok := RequestMetaFrom(r.Context())
		if !ok {
			t.Fatalf("request meta missing")
		}
		clientIP = meta.ClientIP
	})(rr, req)

	if clientIP != "10.0.0.1" {
		t.Fatalf("expected first forwarded IP, got %s", clientIP)
	}
}

func TestRequestMetaMiddlewareIgnoresForwardHeadersFromUntrustedSource(t *testing.T) {
	t.Parallel()

	m, err := NewRequestMetaMiddleware(config.GatewaySecurityConf{
		TrustedProxyCIDRs:   []string{"10.0.0.0/8"},
		TrustedProxyHeaders: []string{"X-Forwarded-For", "X-Real-IP", "Client-IP"},
	})
	if err != nil {
		t.Fatalf("expected middleware construction to succeed, got %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.RemoteAddr = "127.0.0.1:4567"
	req.Header.Set("X-Forwarded-For", "203.0.113.10")
	req.Header.Set("X-Real-IP", "203.0.113.11")
	req.Header.Set("Client-IP", "203.0.113.12")
	rr := httptest.NewRecorder()

	var clientIP string
	m.Handle(func(_ http.ResponseWriter, r *http.Request) {
		meta, ok := RequestMetaFrom(r.Context())
		if !ok {
			t.Fatalf("request meta missing")
		}
		clientIP = meta.ClientIP
	})(rr, req)

	if clientIP != "127.0.0.1" {
		t.Fatalf("expected socket ip for untrusted source, got %s", clientIP)
	}
}

func TestRequestMetaMiddlewareInvalidRemoteAddrFallsBackWithoutPanic(t *testing.T) {
	t.Parallel()

	m, err := NewRequestMetaMiddleware(config.GatewaySecurityConf{
		TrustedProxyCIDRs:   []string{"127.0.0.0/8"},
		TrustedProxyHeaders: []string{"X-Forwarded-For"},
	})
	if err != nil {
		t.Fatalf("expected middleware construction to succeed, got %v", err)
	}
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.RemoteAddr = "not-an-ip"
	req.Header.Set("X-Forwarded-For", "203.0.113.10")
	rr := httptest.NewRecorder()

	var clientIP string
	m.Handle(func(_ http.ResponseWriter, r *http.Request) {
		meta, ok := RequestMetaFrom(r.Context())
		if !ok {
			t.Fatalf("request meta missing")
		}
		clientIP = meta.ClientIP
	})(rr, req)

	if clientIP != "" {
		t.Fatalf("expected empty client ip for invalid remote addr, got %s", clientIP)
	}
}

func TestRequestMetaMiddlewareRejectsInvalidTrustedProxyCIDRs(t *testing.T) {
	t.Parallel()

	if _, err := NewRequestMetaMiddleware(config.GatewaySecurityConf{
		TrustedProxyCIDRs:   []string{"not-a-cidr"},
		TrustedProxyHeaders: []string{"X-Forwarded-For"},
	}); err == nil {
		t.Fatalf("expected invalid trusted proxy cidr error")
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
