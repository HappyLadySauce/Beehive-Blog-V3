package file

import (
	"context"
	"errors"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	filepb "github.com/HappyLadySauce/Beehive-Blog-V3/services/file/pb"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/middleware"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeFileClient struct {
	getFileConfigFn    func(context.Context, *filepb.GetFileConfigRequest, ...grpc.CallOption) (*filepb.GetFileConfigResponse, error)
	updateFileConfigFn func(context.Context, *filepb.UpdateFileConfigRequest, ...grpc.CallOption) (*filepb.UpdateFileConfigResponse, error)
}

func (f *fakeFileClient) Ping(context.Context, *filepb.PingRequest, ...grpc.CallOption) (*filepb.PingResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeFileClient) CreateUpload(context.Context, *filepb.CreateUploadRequest, ...grpc.CallOption) (*filepb.CreateUploadResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeFileClient) CompleteUpload(context.Context, *filepb.CompleteUploadRequest, ...grpc.CallOption) (*filepb.AssetResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeFileClient) ListAssets(context.Context, *filepb.ListAssetsRequest, ...grpc.CallOption) (*filepb.ListAssetsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeFileClient) GetAsset(context.Context, *filepb.GetAssetRequest, ...grpc.CallOption) (*filepb.AssetResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeFileClient) DeleteAsset(context.Context, *filepb.DeleteAssetRequest, ...grpc.CallOption) (*filepb.DeleteAssetResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (f *fakeFileClient) GetFileConfig(ctx context.Context, in *filepb.GetFileConfigRequest, opts ...grpc.CallOption) (*filepb.GetFileConfigResponse, error) {
	if f.getFileConfigFn == nil {
		return nil, status.Error(codes.Unimplemented, "not implemented")
	}
	return f.getFileConfigFn(ctx, in, opts...)
}

func (f *fakeFileClient) UpdateFileConfig(ctx context.Context, in *filepb.UpdateFileConfigRequest, opts ...grpc.CallOption) (*filepb.UpdateFileConfigResponse, error) {
	if f.updateFileConfigFn == nil {
		return nil, status.Error(codes.Unimplemented, "not implemented")
	}
	return f.updateFileConfigFn(ctx, in, opts...)
}

func TestFileConfigGetRequiresAdmin(t *testing.T) {
	t.Parallel()

	logic := NewFileConfigGetLogic(fileTrustedContext("member"), testFileServiceContext(&fakeFileClient{}))
	_, err := logic.FileConfigGet()
	if !errors.Is(err, errs.E(errs.CodeGatewayAccessForbidden)) {
		t.Fatalf("expected gateway access forbidden, got %v", err)
	}
}

func TestFileConfigUpdateRequiresAdmin(t *testing.T) {
	t.Parallel()

	logic := NewFileConfigUpdateLogic(fileTrustedContext("member"), testFileServiceContext(&fakeFileClient{}))
	_, err := logic.FileConfigUpdate(&types.FileConfigUpdateReq{MaxUploadBytes: 1024})
	if !errors.Is(err, errs.E(errs.CodeGatewayAccessForbidden)) {
		t.Fatalf("expected gateway access forbidden, got %v", err)
	}
}

func TestFileConfigUpdatePassesThroughForAdmin(t *testing.T) {
	t.Parallel()

	client := &fakeFileClient{
		updateFileConfigFn: func(_ context.Context, in *filepb.UpdateFileConfigRequest, _ ...grpc.CallOption) (*filepb.UpdateFileConfigResponse, error) {
			if in.GetMaxUploadBytes() != 1024 {
				t.Fatalf("unexpected update request: %+v", in)
			}
			return &filepb.UpdateFileConfigResponse{
				Config: &filepb.FileConfig{
					MaxUploadBytes:      1024,
					AllowedContentTypes: []string{"image/png"},
					PresignTtlSeconds:   300,
				},
			}, nil
		},
	}

	resp, err := NewFileConfigUpdateLogic(fileTrustedContext("admin"), testFileServiceContext(client)).FileConfigUpdate(
		&types.FileConfigUpdateReq{MaxUploadBytes: 1024},
	)
	if err != nil {
		t.Fatalf("expected update success, got %v", err)
	}
	if resp.Config.MaxUploadBytes != 1024 {
		t.Fatalf("expected updated max upload bytes, got %+v", resp.Config)
	}
}

func fileTrustedContext(role string) context.Context {
	ctx := middleware.WithRequestMeta(context.Background(), ctxmeta.RequestMeta{RequestID: "req-1"})
	return middleware.WithAuthContext(ctx, middleware.AuthContext{UserID: "1", SessionID: "s1", Role: role})
}

func testFileServiceContext(client filepb.FileServiceClient) *svc.ServiceContext {
	return &svc.ServiceContext{
		Config:     config.Config{FileRPC: config.FileRPCConf{InternalAuthToken: "secret", InternalCallerName: "gateway"}},
		FileClient: client,
	}
}
