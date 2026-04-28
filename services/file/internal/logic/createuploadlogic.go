package logic

import (
	"context"

	fileservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/pb"
)

type CreateUploadLogic struct {
	baseLogic
}

func NewCreateUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUploadLogic {
	return &CreateUploadLogic{baseLogic: newBaseLogic(ctx, svcCtx)}
}

func (l *CreateUploadLogic) CreateUpload(in *pb.CreateUploadRequest) (*pb.CreateUploadResponse, error) {
	result, err := l.svcCtx.Services.CreateUpload(l.ctx, fileservice.CreateUploadInput{
		ActorUserID: in.GetActorUserId(),
		Scope:       toServiceScope(in.GetScope()),
		FileName:    in.GetFileName(),
		ContentType: in.GetContentType(),
		ByteSize:    in.GetByteSize(),
		Visibility:  toServiceVisibility(in.GetVisibility()),
	})
	if err != nil {
		return nil, toStatus(err, "create file upload failed")
	}
	return &pb.CreateUploadResponse{
		Asset:     toProtoAsset(result.Asset),
		UploadUrl: result.UploadURL,
		Headers:   result.Headers,
		ExpiresAt: result.ExpiresAt.Unix(),
		MaxBytes:  result.MaxBytes,
	}, nil
}
