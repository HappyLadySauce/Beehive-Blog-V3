package file

import (
	"context"

	fileadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/file"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type FileCategoryExtensionsUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileCategoryExtensionsUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileCategoryExtensionsUpdateLogic {
	return &FileCategoryExtensionsUpdateLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *FileCategoryExtensionsUpdateLogic) FileCategoryExtensionsUpdate(req *types.FileCategoryExtensionsUpdateReq) (*types.FileCategoryResp, error) {
	rpcCtx, _, err := rpcContextWithAdminActor(l.ctx, l.svcCtx.Config.FileRPC)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.FileClient.UpdateFileCategoryExtensions(rpcCtx, fileadapter.BuildUpdateFileCategoryExtensionsRequest(req))
	if rpcErr != nil {
		return nil, mapFileError(l.ctx, "file_category_extensions_update", "/api/v3/studio/file/categories/:category_key/extensions", rpcErr)
	}
	return fileadapter.ToFileCategoryResp(rpcResp), nil
}
