package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/middleware"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type fakeIdentityClient struct {
	registerFn       func(ctx context.Context, in *pb.RegisterLocalUserRequest, opts ...grpc.CallOption) (*pb.RegisterLocalUserResponse, error)
	ssoStartFn       func(ctx context.Context, in *pb.StartSsoLoginRequest, opts ...grpc.CallOption) (*pb.StartSsoLoginResponse, error)
	logoutFn         func(ctx context.Context, in *pb.LogoutSessionRequest, opts ...grpc.CallOption) (*pb.LogoutSessionResponse, error)
	getCurrentFn     func(ctx context.Context, in *pb.GetCurrentUserRequest, opts ...grpc.CallOption) (*pb.GetCurrentUserResponse, error)
	introspectFn     func(ctx context.Context, in *pb.IntrospectAccessTokenRequest, opts ...grpc.CallOption) (*pb.IntrospectAccessTokenResponse, error)
	loginFn          func(ctx context.Context, in *pb.LoginLocalUserRequest, opts ...grpc.CallOption) (*pb.LoginLocalUserResponse, error)
	finishSsoFn      func(ctx context.Context, in *pb.FinishSsoLoginRequest, opts ...grpc.CallOption) (*pb.FinishSsoLoginResponse, error)
	refreshTokenFn   func(ctx context.Context, in *pb.RefreshSessionTokenRequest, opts ...grpc.CallOption) (*pb.RefreshSessionTokenResponse, error)
	updateProfileFn  func(ctx context.Context, in *pb.UpdateOwnProfileRequest, opts ...grpc.CallOption) (*pb.UpdateOwnProfileResponse, error)
	changePasswordFn func(ctx context.Context, in *pb.ChangeOwnPasswordRequest, opts ...grpc.CallOption) (*pb.ChangeOwnPasswordResponse, error)
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

func (f *fakeIdentityClient) ListUsers(context.Context, *pb.ListUsersRequest, ...grpc.CallOption) (*pb.ListUsersResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeIdentityClient) UpdateOwnProfile(ctx context.Context, in *pb.UpdateOwnProfileRequest, opts ...grpc.CallOption) (*pb.UpdateOwnProfileResponse, error) {
	if f.updateProfileFn == nil {
		return nil, status.Error(codes.Unimplemented, "not implemented")
	}
	return f.updateProfileFn(ctx, in, opts...)
}

func (f *fakeIdentityClient) ChangeOwnPassword(ctx context.Context, in *pb.ChangeOwnPasswordRequest, opts ...grpc.CallOption) (*pb.ChangeOwnPasswordResponse, error) {
	if f.changePasswordFn == nil {
		return nil, status.Error(codes.Unimplemented, "not implemented")
	}
	return f.changePasswordFn(ctx, in, opts...)
}

func (f *fakeIdentityClient) UpdateUserRole(context.Context, *pb.UpdateUserRoleRequest, ...grpc.CallOption) (*pb.UpdateUserRoleResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeIdentityClient) UpdateUserStatus(context.Context, *pb.UpdateUserStatusRequest, ...grpc.CallOption) (*pb.UpdateUserStatusResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeIdentityClient) ResetUserPassword(context.Context, *pb.ResetUserPasswordRequest, ...grpc.CallOption) (*pb.ResetUserPasswordResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeIdentityClient) ListIdentityAudits(context.Context, *pb.ListIdentityAuditsRequest, ...grpc.CallOption) (*pb.ListIdentityAuditsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
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
		registerFn: func(ctx context.Context, in *pb.RegisterLocalUserRequest, _ ...grpc.CallOption) (*pb.RegisterLocalUserResponse, error) {
			if in.GetUsername() != "alice" {
				t.Fatalf("unexpected username %s", in.GetUsername())
			}
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
			return &pb.RegisterLocalUserResponse{
				TokenPair: &pb.TokenPair{AccessToken: "a", RefreshToken: "r", ExpiresIn: 100, TokenType: "Bearer"},
				CurrentUser: &pb.CurrentUser{
					UserId: "1", Username: "alice", Email: "alice@example.com", Role: pb.Role_ROLE_MEMBER, Status: pb.AccountStatus_ACCOUNT_STATUS_ACTIVE,
				},
				SessionInfo: &pb.SessionInfo{SessionId: "s1", UserId: "1", AuthSource: pb.AuthSource_AUTH_SOURCE_LOCAL, Status: pb.SessionStatus_SESSION_STATUS_ACTIVE},
			}, nil
		},
	}
	logic := NewAuthRegisterLogic(context.Background(), &svc.ServiceContext{
		Config: config.Config{
			IdentityRPC: config.IdentityRPCConf{
				InternalAuthToken:  "secret",
				InternalCallerName: "gateway",
			},
		},
		IdentityClient: client,
	})
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
	logic := NewAuthSsoStartLogic(context.Background(), &svc.ServiceContext{
		Config: config.Config{
			IdentityRPC: config.IdentityRPCConf{
				InternalAuthToken:  "secret",
				InternalCallerName: "gateway",
			},
		},
		IdentityClient: client,
	})
	_, err := logic.AuthSsoStart(&types.AuthSsoStartReq{Provider: "qq", RedirectUri: "https://example.com/cb"})
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, errs.E(errs.CodeGatewayBadRequest)) && !errors.Is(err, errs.E(errs.CodeIdentitySSOProviderNotReady)) {
		t.Fatalf("expected mapped domain error, got %v", err)
	}
}

func TestAuthSsoStartMapsSuccessForQQ(t *testing.T) {
	t.Parallel()

	client := &fakeIdentityClient{
		ssoStartFn: func(_ context.Context, in *pb.StartSsoLoginRequest, _ ...grpc.CallOption) (*pb.StartSsoLoginResponse, error) {
			if in.GetProvider() != "qq" {
				t.Fatalf("expected qq provider, got %s", in.GetProvider())
			}
			return &pb.StartSsoLoginResponse{
				Provider: "qq",
				AuthUrl:  "https://graph.qq.com/oauth2.0/authorize?state=state-1",
				State:    "state-1",
			}, nil
		},
	}

	logic := NewAuthSsoStartLogic(context.Background(), &svc.ServiceContext{
		Config: config.Config{
			IdentityRPC: config.IdentityRPCConf{
				InternalAuthToken:  "secret",
				InternalCallerName: "gateway",
			},
		},
		IdentityClient: client,
	})
	resp, err := logic.AuthSsoStart(&types.AuthSsoStartReq{Provider: "qq", RedirectUri: "https://example.com/auth/qq/callback"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Provider != "qq" || resp.AuthUrl == "" || resp.State != "state-1" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestAuthSsoCallbackMapsSuccessForWeChat(t *testing.T) {
	t.Parallel()

	client := &fakeIdentityClient{
		finishSsoFn: func(_ context.Context, in *pb.FinishSsoLoginRequest, _ ...grpc.CallOption) (*pb.FinishSsoLoginResponse, error) {
			if in.GetProvider() != "wechat" {
				t.Fatalf("expected wechat provider, got %s", in.GetProvider())
			}
			return &pb.FinishSsoLoginResponse{
				TokenPair: &pb.TokenPair{AccessToken: "a", RefreshToken: "r", ExpiresIn: 300, TokenType: "Bearer"},
				CurrentUser: &pb.CurrentUser{
					UserId: "1", Username: "wechat_user", Role: pb.Role_ROLE_MEMBER, Status: pb.AccountStatus_ACCOUNT_STATUS_ACTIVE,
				},
				SessionInfo: &pb.SessionInfo{SessionId: "s1", UserId: "1", AuthSource: pb.AuthSource_AUTH_SOURCE_SSO, Status: pb.SessionStatus_SESSION_STATUS_ACTIVE},
			}, nil
		},
	}

	logic := NewAuthSsoCallbackLogic(context.Background(), &svc.ServiceContext{
		Config: config.Config{
			IdentityRPC: config.IdentityRPCConf{
				InternalAuthToken:  "secret",
				InternalCallerName: "gateway",
			},
		},
		IdentityClient: client,
	})
	resp, err := logic.AuthSsoCallback(&types.AuthSsoCallbackReq{
		Provider:    "wechat",
		Code:        "code-1",
		State:       "state-1",
		RedirectUri: "https://example.com/auth/wechat/callback",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.SessionId != "s1" || resp.User.UserId != "1" {
		t.Fatalf("unexpected response: %+v", resp)
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
	logic := NewAuthLogoutLogic(ctx, &svc.ServiceContext{
		Config: config.Config{
			IdentityRPC: config.IdentityRPCConf{
				InternalAuthToken:  "secret",
				InternalCallerName: "gateway",
			},
		},
		IdentityClient: client,
	})
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

	logic := NewAuthMeLogic(context.Background(), &svc.ServiceContext{Config: config.Config{
		IdentityRPC: config.IdentityRPCConf{
			InternalAuthToken:  "secret",
			InternalCallerName: "gateway",
		},
	}, IdentityClient: &fakeIdentityClient{
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

func TestAuthUpdateProfileUsesTrustedUser(t *testing.T) {
	t.Parallel()

	client := &fakeIdentityClient{
		updateProfileFn: func(_ context.Context, in *pb.UpdateOwnProfileRequest, _ ...grpc.CallOption) (*pb.UpdateOwnProfileResponse, error) {
			if in.GetUserId() != "1" || in.GetNickname() != "Alice" || in.GetAvatarUrl() != "https://cdn.example.com/a.png" {
				t.Fatalf("unexpected update profile request: %+v", in)
			}
			if in.Nickname == nil || in.AvatarUrl == nil {
				t.Fatalf("expected profile patch fields to preserve presence: %+v", in)
			}
			return &pb.UpdateOwnProfileResponse{CurrentUser: &pb.CurrentUser{
				UserId: "1", Username: "alice", Nickname: "Alice", AvatarUrl: "https://cdn.example.com/a.png", Role: pb.Role_ROLE_MEMBER, Status: pb.AccountStatus_ACCOUNT_STATUS_ACTIVE,
			}}, nil
		},
	}
	ctx := middleware.WithAuthContext(context.Background(), middleware.AuthContext{UserID: "1"})
	logic := NewAuthUpdateProfileLogic(ctx, &svc.ServiceContext{
		Config:         config.Config{IdentityRPC: config.IdentityRPCConf{InternalAuthToken: "secret", InternalCallerName: "gateway"}},
		IdentityClient: client,
	})

	resp, err := logic.AuthUpdateProfile(&types.AuthUpdateProfileReq{
		Nickname:  stringPtr("Alice"),
		AvatarUrl: stringPtr("https://cdn.example.com/a.png"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.User.UserId != "1" || resp.User.Role != "member" || resp.User.Status != "active" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestAuthChangePasswordUsesTrustedUser(t *testing.T) {
	t.Parallel()

	client := &fakeIdentityClient{
		changePasswordFn: func(_ context.Context, in *pb.ChangeOwnPasswordRequest, _ ...grpc.CallOption) (*pb.ChangeOwnPasswordResponse, error) {
			if in.GetUserId() != "1" || in.GetOldPassword() != "old" || in.GetNewPassword() != "new-password" {
				t.Fatalf("unexpected change password request: %+v", in)
			}
			return &pb.ChangeOwnPasswordResponse{Ok: true}, nil
		},
	}
	ctx := middleware.WithAuthContext(context.Background(), middleware.AuthContext{UserID: "1"})
	logic := NewAuthChangePasswordLogic(ctx, &svc.ServiceContext{
		Config:         config.Config{IdentityRPC: config.IdentityRPCConf{InternalAuthToken: "secret", InternalCallerName: "gateway"}},
		IdentityClient: client,
	})

	resp, err := logic.AuthChangePassword(&types.AuthChangePasswordReq{OldPassword: "old", NewPassword: "new-password"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Ok {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func stringPtr(value string) *string {
	return &value
}
