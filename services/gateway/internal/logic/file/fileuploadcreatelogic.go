package file

import (
	"context"

	fileadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/file"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type FileUploadCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadCreateLogic {
	return &FileUploadCreateLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *FileUploadCreateLogic) FileUploadCreate(req *types.FileUploadCreateReq) (*types.FileUploadCreateResp, error) {
	rpcCtx, actorUserID, err := rpcContextWithActor(l.ctx, l.svcCtx.Config.FileRPC)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.FileClient.CreateUpload(rpcCtx, fileadapter.BuildCreateUploadRequest(actorUserID, req))
	if rpcErr != nil {
		return nil, mapFileError(l.ctx, "file_upload_create", "/api/v3/files/uploads", rpcErr)
	}
	return fileadapter.ToCreateUploadResp(rpcResp), nil
}
