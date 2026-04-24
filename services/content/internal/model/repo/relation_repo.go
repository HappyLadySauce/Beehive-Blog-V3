package repo

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/entity"
	"gorm.io/gorm"
)

type RelationRepository struct {
	db *gorm.DB
}

func (r *RelationRepository) Create(ctx context.Context, relation *entity.Relation) error {
	return r.db.WithContext(ctx).Create(relation).Error
}

func (r *RelationRepository) GetByIDForContent(ctx context.Context, contentID, relationID int64) (*entity.Relation, error) {
	var relation entity.Relation
	if err := r.db.WithContext(ctx).First(&relation, "id = ? AND from_content_id = ?", relationID, contentID).Error; err != nil {
		return nil, err
	}
	return &relation, nil
}

func (r *RelationRepository) Delete(ctx context.Context, relation *entity.Relation) error {
	result := r.db.WithContext(ctx).Delete(relation)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

type RelationFilter struct {
	FromContentID int64
	RelationType  string
}

func (r *RelationRepository) List(ctx context.Context, filter RelationFilter, page, pageSize int) ([]entity.Relation, int64, error) {
	q := r.db.WithContext(ctx).Model(&entity.Relation{}).Where("from_content_id = ?", filter.FromContentID)
	if filter.RelationType != "" {
		q = q.Where("relation_type = ?", filter.RelationType)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var relations []entity.Relation
	if err := q.Order("sort_order ASC, weight DESC, id DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&relations).Error; err != nil {
		return nil, 0, err
	}
	return relations, total, nil
}

func (r *RelationRepository) CountByContent(ctx context.Context, contentID int64) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&entity.Relation{}).Where("from_content_id = ? OR to_content_id = ?", contentID, contentID).Count(&total).Error
	return total, err
}
