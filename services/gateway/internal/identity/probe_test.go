package identity

import (
	"context"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type fakeIdentityClient struct {
	pingFn func(ctx context.Context, in *pb.PingRequest, opts ...grpc.CallOption) (*pb.PingResponse, error)
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
func (f *fakeIdentityClient) IntrospectAccessToken(context.Context, *pb.IntrospectAccessTokenRequest, ...grpc.CallOption) (*pb.IntrospectAccessTokenResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
func (f *fakeIdentityClient) Ping(ctx context.Context, in *pb.PingRequest, opts ...grpc.CallOption) (*pb.PingResponse, error) {
	return f.pingFn(ctx, in, opts...)
}

// TestProbeCheckInjectsInternalAuth verifies readiness probe uses authenticated internal metadata.
// TestProbeCheckInjectsInternalAuth 验证就绪探针会注入已认证的内部 metadata。
func TestProbeCheckInjectsInternalAuth(t *testing.T) {
	t.Parallel()

	probe := NewProbe(&fakeIdentityClient{
		pingFn: func(ctx context.Context, _ *pb.PingRequest, _ ...grpc.CallOption) (*pb.PingResponse, error) {
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

			return &pb.PingResponse{Ok: true, Service: "identity"}, nil
		},
	}, config.IdentityRPCConf{
		InternalAuthToken:  "secret",
		InternalCallerName: "gateway",
	})

	if err := probe.Check(context.Background()); err != nil {
		t.Fatalf("expected probe check to succeed, got %v", err)
	}
}
