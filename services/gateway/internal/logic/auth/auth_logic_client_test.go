package auth

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
