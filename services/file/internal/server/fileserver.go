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

func (s *FileServer) GetAsset(ctx context.Context, in *pb.GetAssetRequest) (*pb.AssetResponse, error) {
	l := logic.NewGetAssetLogic(ctx, s.svcCtx)
	return l.GetAsset(in)
}

func (s *FileServer) DeleteAsset(ctx context.Context, in *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error) {
	l := logic.NewDeleteAssetLogic(ctx, s.svcCtx)
	return l.DeleteAsset(in)
}
