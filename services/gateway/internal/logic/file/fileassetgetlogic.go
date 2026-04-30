// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package file

import (
	"context"

	fileadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/file"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type FileAssetGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileAssetGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileAssetGetLogic {
	return &FileAssetGetLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *FileAssetGetLogic) FileAssetGet(req *types.FileAssetIdReq) (*types.FileAssetResp, error) {
	rpcCtx, actorUserID, err := rpcContextWithActor(l.ctx, l.svcCtx.Config.FileRPC)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.FileClient.GetAsset(rpcCtx, fileadapter.BuildGetAssetRequest(actorUserID, req))
	if rpcErr != nil {
		return nil, mapFileError(l.ctx, "file_asset_get", "/api/v3/files/assets/:asset_id", rpcErr)
	}
	return fileadapter.ToAssetResp(rpcResp), nil
}
