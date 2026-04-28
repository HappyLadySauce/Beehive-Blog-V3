package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/pb"
)

type CompleteUploadLogic struct {
	baseLogic
}

func NewCompleteUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteUploadLogic {
	return &CompleteUploadLogic{baseLogic: newBaseLogic(ctx, svcCtx)}
}

func (l *CompleteUploadLogic) CompleteUpload(in *pb.CompleteUploadRequest) (*pb.AssetResponse, error) {
	asset, err := l.svcCtx.Services.CompleteUpload(l.ctx, in.GetActorUserId(), in.GetUploadId())
	if err != nil {
		return nil, toStatus(err, "complete file upload failed")
	}
	return &pb.AssetResponse{Asset: toProtoAsset(asset)}, nil
}
