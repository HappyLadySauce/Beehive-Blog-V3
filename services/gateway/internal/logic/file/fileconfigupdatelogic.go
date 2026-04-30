package file

import (
	"context"

	fileadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/file"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type FileConfigUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileConfigUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileConfigUpdateLogic {
	return &FileConfigUpdateLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *FileConfigUpdateLogic) FileConfigUpdate(req *types.FileConfigUpdateReq) (*types.FileConfigUpdateResp, error) {
	rpcCtx, _, err := rpcContextWithActor(l.ctx, l.svcCtx.Config.FileRPC)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.FileClient.UpdateFileConfig(rpcCtx, fileadapter.BuildUpdateFileConfigRequest(req))
	if rpcErr != nil {
		return nil, mapFileError(l.ctx, "file_config_update", "/api/v3/studio/file/config", rpcErr)
	}
	return fileadapter.ToFileConfigUpdateResp(rpcResp), nil
}
