package repo

import (
	"context"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/model/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AssetRepository struct {
	db *gorm.DB
}

func (r *AssetRepository) Create(ctx context.Context, asset *entity.FileAsset) error {
	return r.db.WithContext(ctx).Create(asset).Error
}

func (r *AssetRepository) FindByAssetID(ctx context.Context, assetID string) (*entity.FileAsset, error) {
	var asset entity.FileAsset
	if err := r.db.WithContext(ctx).Where("asset_id = ?", assetID).First(&asset).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *AssetRepository) FindByUploadID(ctx context.Context, uploadID string) (*entity.FileAsset, error) {
	var asset entity.FileAsset
	if err := r.db.WithContext(ctx).Where("upload_id = ?", uploadID).First(&asset).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *AssetRepository) FindByUploadIDForUpdate(ctx context.Context, uploadID string) (*entity.FileAsset, error) {
	var asset entity.FileAsset
	err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("upload_id = ?", uploadID).
		First(&asset).
		Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *AssetRepository) MarkUploaded(ctx context.Context, assetID string, uploadedAt time.Time, byteSize int64, contentType string) (bool, error) {
	result := r.db.WithContext(ctx).
		Model(&entity.FileAsset{}).
		Where("asset_id = ? AND status = ?", assetID, "pending").
		Updates(map[string]any{
			"status":       "uploaded",
			"uploaded_at":  uploadedAt,
			"byte_size":    byteSize,
			"content_type": contentType,
			"updated_at":   uploadedAt,
		})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *AssetRepository) MarkDeleted(ctx context.Context, assetID string, deletedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&entity.FileAsset{}).
		Where("asset_id = ? AND status <> ?", assetID, "deleted").
		Updates(map[string]any{
			"status":     "deleted",
			"deleted_at": deletedAt,
			"updated_at": deletedAt,
		}).
		Error
}
