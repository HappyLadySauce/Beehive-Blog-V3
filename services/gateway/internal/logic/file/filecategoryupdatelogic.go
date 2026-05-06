package file

import (
	"context"

	fileadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/file"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type FileCategoryUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileCategoryUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileCategoryUpdateLogic {
	return &FileCategoryUpdateLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *FileCategoryUpdateLogic) FileCategoryUpdate(req *types.FileCategoryUpdateReq) (*types.FileCategoryResp, error) {
	rpcCtx, _, err := rpcContextWithAdminActor(l.ctx, l.svcCtx.Config.FileRPC)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.FileClient.UpdateFileCategory(rpcCtx, fileadapter.BuildUpdateFileCategoryRequest(req))
	if rpcErr != nil {
		return nil, mapFileError(l.ctx, "file_category_update", "/api/v3/studio/file/categories/:category_key", rpcErr)
	}
	return fileadapter.ToFileCategoryResp(rpcResp), nil
}
