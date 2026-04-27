package identity

import (
	"context"
	"errors"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/middleware"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeIdentityClient struct {
	listUsersFn          func(context.Context, *pb.ListUsersRequest, ...grpc.CallOption) (*pb.ListUsersResponse, error)
	updateRoleFn         func(context.Context, *pb.UpdateUserRoleRequest, ...grpc.CallOption) (*pb.UpdateUserRoleResponse, error)
	updateStatusFn       func(context.Context, *pb.UpdateUserStatusRequest, ...grpc.CallOption) (*pb.UpdateUserStatusResponse, error)
	resetPasswordFn      func(context.Context, *pb.ResetUserPasswordRequest, ...grpc.CallOption) (*pb.ResetUserPasswordResponse, error)
	deleteUserFn         func(context.Context, *pb.DeleteUserRequest, ...grpc.CallOption) (*pb.DeleteUserResponse, error)
	listIdentityAuditsFn func(context.Context, *pb.ListIdentityAuditsRequest, ...grpc.CallOption) (*pb.ListIdentityAuditsResponse, error)
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

func (f *fakeIdentityClient) ListUsers(ctx context.Context, in *pb.ListUsersRequest, opts ...grpc.CallOption) (*pb.ListUsersResponse, error) {
	return f.listUsersFn(ctx, in, opts...)
}

func (f *fakeIdentityClient) UpdateOwnProfile(context.Context, *pb.UpdateOwnProfileRequest, ...grpc.CallOption) (*pb.UpdateOwnProfileResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeIdentityClient) ChangeOwnPassword(context.Context, *pb.ChangeOwnPasswordRequest, ...grpc.CallOption) (*pb.ChangeOwnPasswordResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeIdentityClient) UpdateUserRole(ctx context.Context, in *pb.UpdateUserRoleRequest, opts ...grpc.CallOption) (*pb.UpdateUserRoleResponse, error) {
	return f.updateRoleFn(ctx, in, opts...)
}

func (f *fakeIdentityClient) UpdateUserStatus(ctx context.Context, in *pb.UpdateUserStatusRequest, opts ...grpc.CallOption) (*pb.UpdateUserStatusResponse, error) {
	return f.updateStatusFn(ctx, in, opts...)
}

func (f *fakeIdentityClient) ResetUserPassword(ctx context.Context, in *pb.ResetUserPasswordRequest, opts ...grpc.CallOption) (*pb.ResetUserPasswordResponse, error) {
	return f.resetPasswordFn(ctx, in, opts...)
}

func (f *fakeIdentityClient) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest, opts ...grpc.CallOption) (*pb.DeleteUserResponse, error) {
	return f.deleteUserFn(ctx, in, opts...)
}

func (f *fakeIdentityClient) ListIdentityAudits(ctx context.Context, in *pb.ListIdentityAuditsRequest, opts ...grpc.CallOption) (*pb.ListIdentityAuditsResponse, error) {
	return f.listIdentityAuditsFn(ctx, in, opts...)
}

func (f *fakeIdentityClient) IntrospectAccessToken(context.Context, *pb.IntrospectAccessTokenRequest, ...grpc.CallOption) (*pb.IntrospectAccessTokenResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeIdentityClient) Ping(context.Context, *pb.PingRequest, ...grpc.CallOption) (*pb.PingResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func TestIdentityUserListRequiresAdmin(t *testing.T) {
	t.Parallel()

	logic := NewIdentityUserListLogic(trustedContext("member"), testServiceContext(&fakeIdentityClient{}))
	_, err := logic.IdentityUserList(&types.AdminUserListReq{})
	if !errors.Is(err, errs.E(errs.CodeGatewayAccessForbidden)) {
		t.Fatalf("expected gateway access forbidden, got %v", err)
	}
}

func TestIdentityUserListRejectsInvalidFilters(t *testing.T) {
	t.Parallel()

	called := false
	client := &fakeIdentityClient{
		listUsersFn: func(context.Context, *pb.ListUsersRequest, ...grpc.CallOption) (*pb.ListUsersResponse, error) {
			called = true
			return &pb.ListUsersResponse{}, nil
		},
	}
	logic := NewIdentityUserListLogic(trustedContext("admin"), testServiceContext(client))

	_, err := logic.IdentityUserList(&types.AdminUserListReq{Role: "guest", Status: "active"})
	if !errors.Is(err, errs.E(errs.CodeGatewayBadRequest)) {
		t.Fatalf("expected bad request for invalid role filter, got %v", err)
	}
	if called {
		t.Fatalf("expected invalid filters to be rejected before RPC")
	}

	_, err = logic.IdentityUserList(&types.AdminUserListReq{Role: "admin", Status: "typo"})
	if !errors.Is(err, errs.E(errs.CodeGatewayBadRequest)) {
		t.Fatalf("expected bad request for invalid status filter, got %v", err)
	}
	if called {
		t.Fatalf("expected invalid status filter to be rejected before RPC")
	}
}

func TestIdentityLogicMapsAdminRPCs(t *testing.T) {
	t.Parallel()

	client := &fakeIdentityClient{
		listUsersFn: func(_ context.Context, in *pb.ListUsersRequest, _ ...grpc.CallOption) (*pb.ListUsersResponse, error) {
			if in.GetActorUserId() != "1" || in.GetRole() != pb.Role_ROLE_ADMIN || in.GetStatus() != pb.AccountStatus_ACCOUNT_STATUS_ACTIVE {
				t.Fatalf("unexpected list users request: %+v", in)
			}
			return &pb.ListUsersResponse{
				Items: []*pb.AdminUserView{{UserId: "2", Username: "bob", Role: pb.Role_ROLE_ADMIN, Status: pb.AccountStatus_ACCOUNT_STATUS_ACTIVE}},
				Total: 1, Page: 1, PageSize: 20,
			}, nil
		},
		updateRoleFn: func(_ context.Context, in *pb.UpdateUserRoleRequest, _ ...grpc.CallOption) (*pb.UpdateUserRoleResponse, error) {
			if in.GetActorUserId() != "1" || in.GetTargetUserId() != "2" || in.GetRole() != pb.Role_ROLE_MEMBER {
				t.Fatalf("unexpected update role request: %+v", in)
			}
			return &pb.UpdateUserRoleResponse{User: &pb.AdminUserView{UserId: "2", Role: pb.Role_ROLE_MEMBER, Status: pb.AccountStatus_ACCOUNT_STATUS_ACTIVE}}, nil
		},
		updateStatusFn: func(_ context.Context, in *pb.UpdateUserStatusRequest, _ ...grpc.CallOption) (*pb.UpdateUserStatusResponse, error) {
			if in.GetStatus() != pb.AccountStatus_ACCOUNT_STATUS_DISABLED {
				t.Fatalf("unexpected update status request: %+v", in)
			}
			return &pb.UpdateUserStatusResponse{User: &pb.AdminUserView{UserId: "2", Role: pb.Role_ROLE_MEMBER, Status: pb.AccountStatus_ACCOUNT_STATUS_DISABLED}}, nil
		},
		resetPasswordFn: func(_ context.Context, in *pb.ResetUserPasswordRequest, _ ...grpc.CallOption) (*pb.ResetUserPasswordResponse, error) {
			if in.GetTargetUserId() != "2" || in.GetNewPassword() != "new-password" {
				t.Fatalf("unexpected reset password request: %+v", in)
			}
			return &pb.ResetUserPasswordResponse{Ok: true}, nil
		},
		deleteUserFn: func(_ context.Context, in *pb.DeleteUserRequest, _ ...grpc.CallOption) (*pb.DeleteUserResponse, error) {
			if in.GetActorUserId() != "1" || in.GetTargetUserId() != "2" {
				t.Fatalf("unexpected delete user request: %+v", in)
			}
			return &pb.DeleteUserResponse{Ok: true}, nil
		},
		listIdentityAuditsFn: func(_ context.Context, in *pb.ListIdentityAuditsRequest, _ ...grpc.CallOption) (*pb.ListIdentityAuditsResponse, error) {
			if in.GetActorUserId() != "1" || in.GetEventType() != "admin_update_user_status" {
				t.Fatalf("unexpected audit list request: %+v", in)
			}
			return &pb.ListIdentityAuditsResponse{
				Items: []*pb.IdentityAuditView{{AuditId: "9001", EventType: "admin_update_user_status", AuthSource: pb.AuthSource_AUTH_SOURCE_LOCAL}},
				Total: 1, Page: 1, PageSize: 20,
			}, nil
		},
	}
	svcCtx := testServiceContext(client)
	ctx := trustedContext("ROLE_ADMIN")

	if resp, err := NewIdentityUserListLogic(ctx, svcCtx).IdentityUserList(&types.AdminUserListReq{Role: "admin", Status: "active", Page: 1, PageSize: 20}); err != nil || resp.Items[0].Role != "admin" {
		t.Fatalf("user list failed: resp=%+v err=%v", resp, err)
	}
	if resp, err := NewIdentityUserRoleUpdateLogic(ctx, svcCtx).IdentityUserRoleUpdate(&types.AdminUpdateUserRoleReq{UserId: "2", Role: "member"}); err != nil || resp.User.Role != "member" {
		t.Fatalf("role update failed: resp=%+v err=%v", resp, err)
	}
	if resp, err := NewIdentityUserStatusUpdateLogic(ctx, svcCtx).IdentityUserStatusUpdate(&types.AdminUpdateUserStatusReq{UserId: "2", Status: "disabled"}); err != nil || resp.User.Status != "disabled" {
		t.Fatalf("status update failed: resp=%+v err=%v", resp, err)
	}
	if resp, err := NewIdentityUserPasswordResetLogic(ctx, svcCtx).IdentityUserPasswordReset(&types.AdminResetUserPasswordReq{UserId: "2", NewPassword: "new-password"}); err != nil || !resp.Ok {
		t.Fatalf("password reset failed: resp=%+v err=%v", resp, err)
	}
	if resp, err := NewIdentityUserDeleteLogic(ctx, svcCtx).IdentityUserDelete(&types.AdminUserIdReq{UserId: "2"}); err != nil || !resp.Ok {
		t.Fatalf("delete user failed: resp=%+v err=%v", resp, err)
	}
	if resp, err := NewIdentityAuditListLogic(ctx, svcCtx).IdentityAuditList(&types.IdentityAuditListReq{EventType: "admin_update_user_status", Page: 1, PageSize: 20}); err != nil || resp.Items[0].AuthSource != "local" {
		t.Fatalf("audit list failed: resp=%+v err=%v", resp, err)
	}
}

func trustedContext(role string) context.Context {
	return middleware.WithAuthContext(context.Background(), middleware.AuthContext{UserID: "1", SessionID: "s1", Role: role})
}

func testServiceContext(client pb.IdentityClient) *svc.ServiceContext {
	return &svc.ServiceContext{
		Config:         config.Config{IdentityRPC: config.IdentityRPCConf{InternalAuthToken: "secret", InternalCallerName: "gateway"}},
		IdentityClient: client,
	}
}
