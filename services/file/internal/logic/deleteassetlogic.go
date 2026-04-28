package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/pb"
)

type DeleteAssetLogic struct {
	baseLogic
}

func NewDeleteAssetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAssetLogic {
	return &DeleteAssetLogic{baseLogic: newBaseLogic(ctx, svcCtx)}
}

func (l *DeleteAssetLogic) DeleteAsset(in *pb.DeleteAssetRequest) (*pb.DeleteAssetResponse, error) {
	if err := l.svcCtx.Services.DeleteAsset(l.ctx, in.GetActorUserId(), in.GetAssetId()); err != nil {
		return nil, toStatus(err, "delete file asset failed")
	}
	return &pb.DeleteAssetResponse{Ok: true}, nil
}
