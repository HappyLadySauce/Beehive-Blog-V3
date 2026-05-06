package service

import (
	"strconv"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/model/repo"
)

func toAssetView(asset *entity.FileAsset) *AssetView {
	if asset == nil {
		return nil
	}
	return &AssetView{
		AssetID:     asset.AssetID,
		UploadID:    asset.UploadID,
		OwnerUserID: strconv.FormatInt(asset.OwnerUserID, 10),
		CategoryKey: asset.CategoryKey,
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

func toFileCategoryView(record repo.FileCategoryRecord) *FileCategoryView {
	return &FileCategoryView{
		CategoryKey:       record.CategoryKey,
		DisplayName:       record.DisplayName,
		Description:       record.Description,
		Enabled:           record.Enabled,
		IsDefault:         record.IsDefault,
		SortOrder:         record.SortOrder,
		AllowedExtensions: append([]string(nil), record.AllowedExtensions...),
		CreatedAt:         record.CreatedAt,
		UpdatedAt:         record.UpdatedAt,
	}
}
