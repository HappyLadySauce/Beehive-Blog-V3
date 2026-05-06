package repo

import (
	"context"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/model/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FileCategoryRecord struct {
	entity.FileCategory
	AllowedExtensions []string
}

type FileCategoryRepository struct {
	db *gorm.DB
}

func (r *FileCategoryRepository) List(ctx context.Context, includeDisabled bool) ([]FileCategoryRecord, error) {
	var categories []entity.FileCategory
	query := r.db.WithContext(ctx).Order("is_default DESC, sort_order ASC, display_name ASC, category_key ASC")
	if !includeDisabled {
		query = query.Where("enabled = ?", true)
	}
	if err := query.Find(&categories).Error; err != nil {
		return nil, err
	}
	return r.withExtensions(ctx, categories)
}

func (r *FileCategoryRepository) FindByKey(ctx context.Context, categoryKey string) (*FileCategoryRecord, error) {
	var category entity.FileCategory
	if err := r.db.WithContext(ctx).Where("category_key = ?", categoryKey).First(&category).Error; err != nil {
		return nil, err
	}
	items, err := r.withExtensions(ctx, []entity.FileCategory{category})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &items[0], nil
}

func (r *FileCategoryRepository) Create(ctx context.Context, category *entity.FileCategory, extensions []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if category.IsDefault {
			if err := tx.Model(&entity.FileCategory{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
				return err
			}
		}
		if err := tx.Create(category).Error; err != nil {
			return err
		}
		return replaceCategoryExtensions(ctx, tx, category.CategoryKey, extensions)
	})
}

func (r *FileCategoryRepository) Update(ctx context.Context, categoryKey string, update map[string]any) error {
	if len(update) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Model(&entity.FileCategory{}).Where("category_key = ?", categoryKey).Updates(update).Error
}

func (r *FileCategoryRepository) ReplaceExtensions(ctx context.Context, categoryKey string, extensions []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return replaceCategoryExtensions(ctx, tx, categoryKey, extensions)
	})
}

func (r *FileCategoryRepository) SetDefault(ctx context.Context, categoryKey string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var category entity.FileCategory
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("category_key = ?", categoryKey).First(&category).Error; err != nil {
			return err
		}
		if err := tx.Model(&entity.FileCategory{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
			return err
		}
		return tx.Model(&entity.FileCategory{}).Where("category_key = ?", categoryKey).Updates(map[string]any{
			"is_default": true,
			"updated_at": time.Now().UTC(),
		}).Error
	})
}

func (r *FileCategoryRepository) withExtensions(ctx context.Context, categories []entity.FileCategory) ([]FileCategoryRecord, error) {
	if len(categories) == 0 {
		return nil, nil
	}

	keys := make([]string, 0, len(categories))
	for _, item := range categories {
		keys = append(keys, item.CategoryKey)
	}

	var rows []entity.FileCategoryExtension
	if err := r.db.WithContext(ctx).
		Where("category_key IN ?", keys).
		Order("extension ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	extensionsByKey := make(map[string][]string, len(keys))
	for _, row := range rows {
		extensionsByKey[row.CategoryKey] = append(extensionsByKey[row.CategoryKey], row.Extension)
	}

	result := make([]FileCategoryRecord, 0, len(categories))
	for _, item := range categories {
		result = append(result, FileCategoryRecord{
			FileCategory:      item,
			AllowedExtensions: append([]string(nil), extensionsByKey[item.CategoryKey]...),
		})
	}
	return result, nil
}

func replaceCategoryExtensions(ctx context.Context, tx *gorm.DB, categoryKey string, extensions []string) error {
	if err := tx.WithContext(ctx).Where("category_key = ?", categoryKey).Delete(&entity.FileCategoryExtension{}).Error; err != nil {
		return err
	}
	if len(extensions) == 0 {
		return nil
	}
	rows := make([]entity.FileCategoryExtension, 0, len(extensions))
	for _, extension := range extensions {
		rows = append(rows, entity.FileCategoryExtension{
			CategoryKey: categoryKey,
			Extension:   extension,
			CreatedAt:   time.Now().UTC(),
		})
	}
	return tx.WithContext(ctx).Create(&rows).Error
}
