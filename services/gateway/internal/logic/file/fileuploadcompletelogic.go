package file

import (
	"context"

	fileadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/file"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type FileUploadCompleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileUploadCompleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadCompleteLogic {
	return &FileUploadCompleteLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *FileUploadCompleteLogic) FileUploadComplete(req *types.FileUploadCompleteReq) (*types.FileAssetResp, error) {
	rpcCtx, actorUserID, err := rpcContextWithActor(l.ctx, l.svcCtx.Config.FileRPC)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.FileClient.CompleteUpload(rpcCtx, fileadapter.BuildCompleteUploadRequest(actorUserID, req))
	if rpcErr != nil {
		return nil, mapFileError(l.ctx, "file_upload_complete", "/api/v3/files/uploads/:upload_id/complete", rpcErr)
	}
	return fileadapter.ToAssetResp(rpcResp), nil
}
