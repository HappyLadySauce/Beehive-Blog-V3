package file

import (
	"context"

	fileadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/file"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type FileCategoryDefaultSetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileCategoryDefaultSetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileCategoryDefaultSetLogic {
	return &FileCategoryDefaultSetLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *FileCategoryDefaultSetLogic) FileCategoryDefaultSet(req *types.FileCategoryDefaultSetReq) (*types.FileCategoryResp, error) {
	rpcCtx, _, err := rpcContextWithAdminActor(l.ctx, l.svcCtx.Config.FileRPC)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.FileClient.SetDefaultFileCategory(rpcCtx, fileadapter.BuildSetDefaultFileCategoryRequest(req))
	if rpcErr != nil {
		return nil, mapFileError(l.ctx, "file_category_default_set", "/api/v3/studio/file/categories/:category_key/default", rpcErr)
	}
	return fileadapter.ToFileCategoryResp(rpcResp), nil
}
