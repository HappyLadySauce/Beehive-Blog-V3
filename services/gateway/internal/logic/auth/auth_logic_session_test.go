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
	"google.golang.org/grpc/metadata"
)

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
