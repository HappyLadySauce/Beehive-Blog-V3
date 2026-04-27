package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
