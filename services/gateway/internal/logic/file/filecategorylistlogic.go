package file

import (
	"context"

	filepb "github.com/HappyLadySauce/Beehive-Blog-V3/services/file/pb"
	fileadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/file"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type FileCategoryListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileCategoryListLogic {
	return &FileCategoryListLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *FileCategoryListLogic) FileCategoryList() (*types.FileCategoryListResp, error) {
	rpcCtx, _, err := rpcContextWithActor(l.ctx, l.svcCtx.Config.FileRPC)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.FileClient.ListFileCategories(rpcCtx, &filepb.ListFileCategoriesRequest{
		IncludeDisabled: false,
	})
	if rpcErr != nil {
		return nil, mapFileError(l.ctx, "file_category_list", "/api/v3/files/categories", rpcErr)
	}
	return fileadapter.ToFileCategoryListResp(rpcResp), nil
}
