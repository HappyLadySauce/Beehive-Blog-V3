package file

import (
	"context"

	filepb "github.com/HappyLadySauce/Beehive-Blog-V3/services/file/pb"
	fileadapter "github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/file"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
)

type StudioFileCategoryListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStudioFileCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StudioFileCategoryListLogic {
	return &StudioFileCategoryListLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *StudioFileCategoryListLogic) StudioFileCategoryList() (*types.FileCategoryListResp, error) {
	rpcCtx, _, err := rpcContextWithAdminActor(l.ctx, l.svcCtx.Config.FileRPC)
	if err != nil {
		return nil, err
	}
	rpcResp, rpcErr := l.svcCtx.FileClient.ListFileCategories(rpcCtx, &filepb.ListFileCategoriesRequest{
		IncludeDisabled: true,
	})
	if rpcErr != nil {
		return nil, mapFileError(l.ctx, "studio_file_category_list", "/api/v3/studio/file/categories", rpcErr)
	}
	return fileadapter.ToFileCategoryListResp(rpcResp), nil
}
