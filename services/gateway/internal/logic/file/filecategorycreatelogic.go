package file

import (
	"context"

	fileadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/file"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type FileCategoryCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileCategoryCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileCategoryCreateLogic {
	return &FileCategoryCreateLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *FileCategoryCreateLogic) FileCategoryCreate(req *types.FileCategoryCreateReq) (*types.FileCategoryResp, error) {
	rpcCtx, _, err := rpcContextWithAdminActor(l.ctx, l.svcCtx.Config.FileRPC)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.FileClient.CreateFileCategory(rpcCtx, fileadapter.BuildCreateFileCategoryRequest(req))
	if rpcErr != nil {
		return nil, mapFileError(l.ctx, "file_category_create", "/api/v3/studio/file/categories", rpcErr)
	}
	return fileadapter.ToFileCategoryResp(rpcResp), nil
}
