package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/pb"
)

type UpdateFileConfigLogic struct {
	baseLogic
}

func NewUpdateFileConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFileConfigLogic {
	return &UpdateFileConfigLogic{baseLogic: newBaseLogic(ctx, svcCtx)}
}

func (l *UpdateFileConfigLogic) UpdateFileConfig(in *pb.UpdateFileConfigRequest) (*pb.UpdateFileConfigResponse, error) {
	if in.GetMaxUploadBytes() > 0 {
		if err := l.svcCtx.ConfigCache.Update(l.ctx, "max_upload_bytes", fmt.Sprintf("%d", in.GetMaxUploadBytes())); err != nil {
			return nil, toStatus(err, "update max_upload_bytes failed")
		}
	}
	if len(in.GetAllowedContentTypes()) > 0 {
		jsonBytes, err := json.Marshal(in.GetAllowedContentTypes())
		if err != nil {
			return nil, toStatus(err, "marshal allowed_content_types failed")
		}
		if err := l.svcCtx.ConfigCache.Update(l.ctx, "allowed_content_types", string(jsonBytes)); err != nil {
			return nil, toStatus(err, "update allowed_content_types failed")
		}
	}
	if in.GetPresignTtlSeconds() > 0 {
		if err := l.svcCtx.ConfigCache.Update(l.ctx, "presign_ttl_seconds", fmt.Sprintf("%d", in.GetPresignTtlSeconds())); err != nil {
			return nil, toStatus(err, "update presign_ttl_seconds failed")
		}
	}

	snapshot := l.svcCtx.ConfigCache.Snapshot()
	return &pb.UpdateFileConfigResponse{
		Config: &pb.FileConfig{
			MaxUploadBytes:      snapshot.MaxUploadBytes,
			AllowedContentTypes: snapshot.AllowedContentTypes,
			PresignTtlSeconds:   int32(snapshot.PresignTTLSeconds),
		},
	}, nil
}
