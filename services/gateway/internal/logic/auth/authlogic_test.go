package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/middleware"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeIdentityClient struct {
	registerFn     func(ctx context.Context, in *pb.RegisterLocalUserRequest, opts ...grpc.CallOption) (*pb.RegisterLocalUserResponse, error)
	ssoStartFn     func(ctx context.Context, in *pb.StartSsoLoginRequest, opts ...grpc.CallOption) (*pb.StartSsoLoginResponse, error)
	logoutFn       func(ctx context.Context, in *pb.LogoutSessionRequest, opts ...grpc.CallOption) (*pb.LogoutSessionResponse, error)
	getCurrentFn   func(ctx context.Context, in *pb.GetCurrentUserRequest, opts ...grpc.CallOption) (*pb.GetCurrentUserResponse, error)
	introspectFn   func(ctx context.Context, in *pb.IntrospectAccessTokenRequest, opts ...grpc.CallOption) (*pb.IntrospectAccessTokenResponse, error)
	loginFn        func(ctx context.Context, in *pb.LoginLocalUserRequest, opts ...grpc.CallOption) (*pb.LoginLocalUserResponse, error)
	finishSsoFn    func(ctx context.Context, in *pb.FinishSsoLoginRequest, opts ...grpc.CallOption) (*pb.FinishSsoLoginResponse, error)
	refreshTokenFn func(ctx context.Context, in *pb.RefreshSessionTokenRequest, opts ...grpc.CallOption) (*pb.RefreshSessionTokenResponse, error)
}

func (f *fakeIdentityClient) RegisterLocalUser(ctx context.Context, in *pb.RegisterLocalUserRequest, opts ...grpc.CallOption) (*pb.RegisterLocalUserResponse, error) {
	return f.registerFn(ctx, in, opts...)
}

func (f *fakeIdentityClient) LoginLocalUser(ctx context.Context, in *pb.LoginLocalUserRequest, opts ...grpc.CallOption) (*pb.LoginLocalUserResponse, error) {
	if f.loginFn == nil {
		return nil, status.Error(codes.Unimplemented, "not implemented")
	}
	return f.loginFn(ctx, in, opts...)
}

func (f *fakeIdentityClient) StartSsoLogin(ctx context.Context, in *pb.StartSsoLoginRequest, opts ...grpc.CallOption) (*pb.StartSsoLoginResponse, error) {
	return f.ssoStartFn(ctx, in, opts...)
}

func (f *fakeIdentityClient) FinishSsoLogin(ctx context.Context, in *pb.FinishSsoLoginRequest, opts ...grpc.CallOption) (*pb.FinishSsoLoginResponse, error) {
	if f.finishSsoFn == nil {
		return nil, status.Error(codes.Unimplemented, "not implemented")
	}
	return f.finishSsoFn(ctx, in, opts...)
}

func (f *fakeIdentityClient) RefreshSessionToken(ctx context.Context, in *pb.RefreshSessionTokenRequest, opts ...grpc.CallOption) (*pb.RefreshSessionTokenResponse, error) {
	if f.refreshTokenFn == nil {
		return nil, status.Error(codes.Unimplemented, "not implemented")
	}
	return f.refreshTokenFn(ctx, in, opts...)
}

func (f *fakeIdentityClient) LogoutSession(ctx context.Context, in *pb.LogoutSessionRequest, opts ...grpc.CallOption) (*pb.LogoutSessionResponse, error) {
	return f.logoutFn(ctx, in, opts...)
}

func (f *fakeIdentityClient) GetCurrentUser(ctx context.Context, in *pb.GetCurrentUserRequest, opts ...grpc.CallOption) (*pb.GetCurrentUserResponse, error) {
	return f.getCurrentFn(ctx, in, opts...)
}

func (f *fakeIdentityClient) IntrospectAccessToken(ctx context.Context, in *pb.IntrospectAccessTokenRequest, opts ...grpc.CallOption) (*pb.IntrospectAccessTokenResponse, error) {
	if f.introspectFn == nil {
		return nil, status.Error(codes.Unimplemented, "not implemented")
	}
	return f.introspectFn(ctx, in, opts...)
}

func (f *fakeIdentityClient) Ping(context.Context, *pb.PingRequest, ...grpc.CallOption) (*pb.PingResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func TestAuthRegisterMapsResponse(t *testing.T) {
	t.Parallel()

	client := &fakeIdentityClient{
		registerFn: func(_ context.Context, in *pb.RegisterLocalUserRequest, _ ...grpc.CallOption) (*pb.RegisterLocalUserResponse, error) {
			if in.GetUsername() != "alice" {
				t.Fatalf("unexpected username %s", in.GetUsername())
			}
			return &pb.RegisterLocalUserResponse{
				TokenPair: &pb.TokenPair{AccessToken: "a", RefreshToken: "r", ExpiresIn: 100, TokenType: "Bearer"},
				CurrentUser: &pb.CurrentUser{
					UserId: "1", Username: "alice", Email: "alice@example.com", Role: pb.Role_ROLE_MEMBER, Status: pb.AccountStatus_ACCOUNT_STATUS_ACTIVE,
				},
				SessionInfo: &pb.SessionInfo{SessionId: "s1", UserId: "1", AuthSource: pb.AuthSource_AUTH_SOURCE_LOCAL, Status: pb.SessionStatus_SESSION_STATUS_ACTIVE},
			}, nil
		},
	}
	logic := NewAuthRegisterLogic(context.Background(), &svc.ServiceContext{IdentityClient: client})
	resp, err := logic.AuthRegister(&types.AuthRegisterReq{Username: "alice", Email: "alice@example.com", Password: "12345678"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.AccessToken != "a" || resp.SessionId != "s1" || resp.User.UserId != "1" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestAuthSsoStartErrorMapsTo412(t *testing.T) {
	t.Parallel()

	client := &fakeIdentityClient{
		ssoStartFn: func(_ context.Context, _ *pb.StartSsoLoginRequest, _ ...grpc.CallOption) (*pb.StartSsoLoginResponse, error) {
			return nil, status.Error(codes.FailedPrecondition, "provider not ready")
		},
	}
	logic := NewAuthSsoStartLogic(context.Background(), &svc.ServiceContext{IdentityClient: client})
	_, err := logic.AuthSsoStart(&types.AuthSsoStartReq{Provider: "qq", RedirectUri: "https://example.com/cb"})
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, errs.E(errs.CodeGatewayBadRequest)) && !errors.Is(err, errs.E(errs.CodeIdentitySSOProviderNotReady)) {
		t.Fatalf("expected mapped domain error, got %v", err)
	}
}

func TestAuthLogoutUsesTrustedSessionID(t *testing.T) {
	t.Parallel()

	var gotSessionID string
	client := &fakeIdentityClient{
		logoutFn: func(_ context.Context, in *pb.LogoutSessionRequest, _ ...grpc.CallOption) (*pb.LogoutSessionResponse, error) {
			gotSessionID = in.GetSessionId()
			return &pb.LogoutSessionResponse{Ok: true}, nil
		},
	}
	ctx := middleware.WithAuthContext(context.Background(), middleware.AuthContext{SessionID: "trusted-session"})
	logic := NewAuthLogoutLogic(ctx, &svc.ServiceContext{IdentityClient: client})
	resp, err := logic.AuthLogout(&types.AuthLogoutReq{RefreshToken: "refresh"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Ok || gotSessionID != "trusted-session" {
		t.Fatalf("unexpected logout result: resp=%+v session=%s", resp, gotSessionID)
	}
}

func TestAuthMeRequiresTrustedContext(t *testing.T) {
	t.Parallel()

	logic := NewAuthMeLogic(context.Background(), &svc.ServiceContext{IdentityClient: &fakeIdentityClient{
		getCurrentFn: func(_ context.Context, _ *pb.GetCurrentUserRequest, _ ...grpc.CallOption) (*pb.GetCurrentUserResponse, error) {
			return &pb.GetCurrentUserResponse{}, nil
		},
	}})
	_, err := logic.AuthMe(&types.AuthMeReq{})
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, errs.E(errs.CodeGatewayAuthorizationRequired)) {
		t.Fatalf("expected gateway authorization required error, got %v", err)
	}
}
