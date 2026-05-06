package server

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/logic"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/pb"
)

type FileServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedFileServiceServer
}

func NewFileServer(svcCtx *svc.ServiceContext) *FileServer {
	return &FileServer{svcCtx: svcCtx}
}

func (s *FileServer) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	l := logic.NewPingLogic(ctx, s.svcCtx)
	return l.Ping(in)
}

func (s *FileServer) CreateUpload(ctx context.Context, in *pb.CreateUploadRequest) (*pb.CreateUploadResponse, error) {
	l := logic.NewCreateUploadLogic(ctx, s.svcCtx)
	return l.CreateUpload(in)
}

func (s *FileServer) CompleteUpload(ctx context.Context, in *pb.CompleteUploadRequest) (*pb.AssetResponse, error) {
	l := logic.NewCompleteUploadLogic(ctx, s.svcCtx)
	return l.CompleteUpload(in)
}

func (s *FileServer) ListAssets(ctx context.Context, in *pb.ListAssetsRequest) (*pb.ListAssetsResponse, error) {
	l := logic.NewListAssetsLogic(ctx, s.svcCtx)
	return l.ListAssets(in)
}

func (s *FileServer) GetAsset(ctx context.Context, in *pb.GetAssetRequest) (*pb.AssetResponse, error) {
	l := logic.NewGetAssetLogic(ctx, s.svcCtx)
	return l.GetAsset(in)
}

func (s *FileServer) DeleteAsset(ctx context.Context, in *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error) {
	l := logic.NewDeleteAssetLogic(ctx, s.svcCtx)
	return l.DeleteAsset(in)
}

func (s *FileServer) GetFileConfig(ctx context.Context, in *pb.GetFileConfigRequest) (*pb.GetFileConfigResponse, error) {
	l := logic.NewGetFileConfigLogic(ctx, s.svcCtx)
	return l.GetFileConfig(in)
}

func (s *FileServer) UpdateFileConfig(ctx context.Context, in *pb.UpdateFileConfigRequest) (*pb.UpdateFileConfigResponse, error) {
	l := logic.NewUpdateFileConfigLogic(ctx, s.svcCtx)
	return l.UpdateFileConfig(in)
}

func (s *FileServer) ListFileCategories(ctx context.Context, in *pb.ListFileCategoriesRequest) (*pb.ListFileCategoriesResponse, error) {
	l := logic.NewListFileCategoriesLogic(ctx, s.svcCtx)
	return l.ListFileCategories(in)
}

func (s *FileServer) CreateFileCategory(ctx context.Context, in *pb.CreateFileCategoryRequest) (*pb.FileCategoryResponse, error) {
	l := logic.NewCreateFileCategoryLogic(ctx, s.svcCtx)
	return l.CreateFileCategory(in)
}

func (s *FileServer) UpdateFileCategory(ctx context.Context, in *pb.UpdateFileCategoryRequest) (*pb.FileCategoryResponse, error) {
	l := logic.NewUpdateFileCategoryLogic(ctx, s.svcCtx)
	return l.UpdateFileCategory(in)
}

func (s *FileServer) UpdateFileCategoryExtensions(ctx context.Context, in *pb.UpdateFileCategoryExtensionsRequest) (*pb.FileCategoryResponse, error) {
	l := logic.NewUpdateFileCategoryExtensionsLogic(ctx, s.svcCtx)
	return l.UpdateFileCategoryExtensions(in)
}

func (s *FileServer) SetDefaultFileCategory(ctx context.Context, in *pb.SetDefaultFileCategoryRequest) (*pb.FileCategoryResponse, error) {
	l := logic.NewSetDefaultFileCategoryLogic(ctx, s.svcCtx)
	return l.SetDefaultFileCategory(in)
}
