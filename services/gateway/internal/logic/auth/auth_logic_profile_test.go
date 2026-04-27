package auth

import (
	"context"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/middleware"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
	"google.golang.org/grpc"
)

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
