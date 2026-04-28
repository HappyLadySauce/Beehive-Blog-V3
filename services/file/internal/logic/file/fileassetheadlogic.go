// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package file

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/types"
)

type FileAssetHeadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Read public uploaded asset metadata
func NewFileAssetHeadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileAssetHeadLogic {
	return &FileAssetHeadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileAssetHeadLogic) FileAssetHead(req *types.FileAssetReadReq) (*AssetHead, error) {
	if l == nil || l.svcCtx == nil || l.svcCtx.Store == nil || l.svcCtx.LocalStorage == nil {
		return nil, storageUnavailable()
	}
	asset, err := loadPublicUploadedAsset(l.ctx, l.svcCtx.Store.Assets, req.AssetId)
	if err != nil {
		return nil, err
	}
	reader, info, err := l.svcCtx.LocalStorage.OpenUploaded(context.WithoutCancel(l.ctx), asset.ObjectKey)
	if err != nil {
		return nil, mapStorageReadError(l.ctx, err)
	}
	_ = reader.Close()
	return &AssetHead{ContentType: asset.ContentType, ByteSize: info.ByteSize}, nil
}
