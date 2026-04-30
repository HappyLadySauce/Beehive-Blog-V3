package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/pb"
)

type GetFileConfigLogic struct {
	baseLogic
}

func NewGetFileConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileConfigLogic {
	return &GetFileConfigLogic{baseLogic: newBaseLogic(ctx, svcCtx)}
}

func (l *GetFileConfigLogic) GetFileConfig(in *pb.GetFileConfigRequest) (*pb.GetFileConfigResponse, error) {
	snapshot := l.svcCtx.ConfigCache.Snapshot()
	return &pb.GetFileConfigResponse{
		Config: &pb.FileConfig{
			MaxUploadBytes:      snapshot.MaxUploadBytes,
			AllowedContentTypes: snapshot.AllowedContentTypes,
			PresignTtlSeconds:   int32(snapshot.PresignTTLSeconds),
		},
	}, nil
}
