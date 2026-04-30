package file

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/pb"
	fileadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/file"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type FileConfigGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileConfigGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileConfigGetLogic {
	return &FileConfigGetLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *FileConfigGetLogic) FileConfigGet() (*types.FileConfigGetResp, error) {
	rpcCtx, _, err := rpcContextWithAdminActor(l.ctx, l.svcCtx.Config.FileRPC)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.FileClient.GetFileConfig(rpcCtx, &pb.GetFileConfigRequest{})
	if rpcErr != nil {
		return nil, mapFileError(l.ctx, "file_config_get", "/api/v3/studio/file/config", rpcErr)
	}
	return fileadapter.ToFileConfigGetResp(rpcResp), nil
}
