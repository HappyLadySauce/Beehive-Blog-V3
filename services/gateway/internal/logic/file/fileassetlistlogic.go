// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package file

import (
	"context"

	fileadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/file"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type FileAssetListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileAssetListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileAssetListLogic {
	return &FileAssetListLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *FileAssetListLogic) FileAssetList(req *types.FileAssetListReq) (*types.FileAssetListResp, error) {
	rpcCtx, actorUserID, err := rpcContextWithActor(l.ctx, l.svcCtx.Config.FileRPC)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.FileClient.ListAssets(rpcCtx, fileadapter.BuildListAssetsRequest(actorUserID, req))
	if rpcErr != nil {
		return nil, mapFileError(l.ctx, "file_asset_list", "/api/v3/files/assets", rpcErr)
	}
	return fileadapter.ToAssetListResp(rpcResp), nil
}
