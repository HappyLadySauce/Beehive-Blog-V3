package repo

import (
	"context"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/model/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AssetRepository struct {
	db *gorm.DB
}

type AssetListFilter struct {
	OwnerUserID int64
	CategoryKey string
	Status      string
	Visibility  string
	Keyword     string
	Page        int
	PageSize    int
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

func (r *AssetRepository) List(ctx context.Context, filter AssetListFilter) ([]entity.FileAsset, int64, error) {
	page, pageSize := normalizePagination(filter.Page, filter.PageSize)
	query := r.db.WithContext(ctx).Model(&entity.FileAsset{}).Where("owner_user_id = ?", filter.OwnerUserID)

	if categoryKey := strings.TrimSpace(filter.CategoryKey); categoryKey != "" {
		query = query.Where("category_key = ?", categoryKey)
	}
	if status := strings.TrimSpace(filter.Status); status != "" {
		query = query.Where("status = ?", status)
	}
	if visibility := strings.TrimSpace(filter.Visibility); visibility != "" {
		query = query.Where("visibility = ?", visibility)
	}
	if keyword := strings.ToLower(strings.TrimSpace(filter.Keyword)); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where(
			"LOWER(file_name) LIKE ? OR LOWER(content_type) LIKE ? OR LOWER(object_key) LIKE ?",
			like,
			like,
			like,
		)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var assets []entity.FileAsset
	if err := query.
		Order("updated_at DESC, created_at DESC, asset_id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&assets).Error; err != nil {
		return nil, 0, err
	}

	return assets, total, nil
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

func normalizePagination(page int, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	switch {
	case pageSize <= 0:
		pageSize = 20
	case pageSize > 100:
		pageSize = 100
	}
	return page, pageSize
}
