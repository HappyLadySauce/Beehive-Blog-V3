package service

import (
	"strconv"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/model/entity"
)

func toAssetView(asset *entity.FileAsset) *AssetView {
	if asset == nil {
		return nil
	}
	return &AssetView{
		AssetID:     asset.AssetID,
		UploadID:    asset.UploadID,
		OwnerUserID: strconv.FormatInt(asset.OwnerUserID, 10),
		Scope:       asset.Scope,
		Visibility:  asset.Visibility,
		Status:      asset.Status,
		Bucket:      asset.Bucket,
		ObjectKey:   asset.ObjectKey,
		PublicURL:   asset.PublicURL,
		FileName:    asset.FileName,
		ContentType: asset.ContentType,
		ByteSize:    asset.ByteSize,
		CreatedAt:   asset.CreatedAt,
		ExpiresAt:   asset.ExpiresAt,
		UploadedAt:  asset.UploadedAt,
		DeletedAt:   asset.DeletedAt,
	}
}
