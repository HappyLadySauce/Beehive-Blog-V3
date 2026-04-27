package middleware

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
func (f *fakeIdentityClient) ListUsers(context.Context, *pb.ListUsersRequest, ...grpc.CallOption) (*pb.ListUsersResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
func (f *fakeIdentityClient) UpdateOwnProfile(context.Context, *pb.UpdateOwnProfileRequest, ...grpc.CallOption) (*pb.UpdateOwnProfileResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
func (f *fakeIdentityClient) ChangeOwnPassword(context.Context, *pb.ChangeOwnPasswordRequest, ...grpc.CallOption) (*pb.ChangeOwnPasswordResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
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
func (f *fakeIdentityClient) DeleteUser(context.Context, *pb.DeleteUserRequest, ...grpc.CallOption) (*pb.DeleteUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
func (f *fakeIdentityClient) ListIdentityAudits(context.Context, *pb.ListIdentityAuditsRequest, ...grpc.CallOption) (*pb.ListIdentityAuditsResponse, error) {
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
