// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package upload

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/middleware"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
	uploadsvc "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/upload"
)

type UploadAvatarPresignLogic struct {
	logger *logs.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadAvatarPresignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadAvatarPresignLogic {
	return &UploadAvatarPresignLogic{
		logger: logs.Ctx(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadAvatarPresignLogic) UploadAvatarPresign(req *types.UploadAvatarPresignReq) (resp *types.UploadAvatarPresignResp, err error) {
	authCtx, ok := middleware.AuthContextFrom(l.ctx)
	if !ok || authCtx.UserID == "" {
		return nil, errs.New(errs.CodeGatewayAuthorizationRequired, "trusted auth context is missing")
	}

	result, err := l.svcCtx.AvatarPresigner.Presign(l.ctx, uploadsvc.AvatarPresignInput{
		UserID:      authCtx.UserID,
		FileName:    req.FileName,
		ContentType: req.ContentType,
		ByteSize:    req.ByteSize,
	})
	if err != nil {
		return nil, err
	}

	l.logger.Info("gateway_avatar_presign_created", logs.String("object_key", result.ObjectKey))

	return &types.UploadAvatarPresignResp{
		UploadUrl: result.UploadURL,
		PublicUrl: result.PublicURL,
		ObjectKey: result.ObjectKey,
		Headers:   result.Headers,
		ExpiresAt: result.ExpiresAt.Unix(),
		MaxBytes:  result.MaxBytes,
	}, nil
}
