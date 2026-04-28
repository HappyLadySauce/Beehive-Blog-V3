package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/pb"
)

type GetAssetLogic struct {
	baseLogic
}

func NewGetAssetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAssetLogic {
	return &GetAssetLogic{baseLogic: newBaseLogic(ctx, svcCtx)}
}

func (l *GetAssetLogic) GetAsset(in *pb.GetAssetRequest) (*pb.AssetResponse, error) {
	asset, err := l.svcCtx.Services.GetAsset(l.ctx, in.GetActorUserId(), in.GetAssetId())
	if err != nil {
		return nil, toStatus(err, "get file asset failed")
	}
	return &pb.AssetResponse{Asset: toProtoAsset(asset)}, nil
}
