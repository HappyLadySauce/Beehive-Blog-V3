package file

import (
	"context"

	fileadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/file"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type FileAssetDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileAssetDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileAssetDeleteLogic {
	return &FileAssetDeleteLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *FileAssetDeleteLogic) FileAssetDelete(req *types.FileAssetIdReq) (*types.FileAssetDeleteResp, error) {
	rpcCtx, actorUserID, err := rpcContextWithActor(l.ctx, l.svcCtx.Config.FileRPC)
	if err != nil {
		return nil, err
	}
	_, rpcErr := l.svcCtx.FileClient.DeleteAsset(rpcCtx, fileadapter.BuildDeleteAssetRequest(actorUserID, req))
	if rpcErr != nil {
		return nil, mapFileError(l.ctx, "file_asset_delete", "/api/v3/files/assets/:asset_id", rpcErr)
	}
	return &types.FileAssetDeleteResp{Ok: true}, nil
}
