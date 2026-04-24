package repo

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TagRepository struct {
	db *gorm.DB
}

func (r *TagRepository) Create(ctx context.Context, tag *entity.Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

func (r *TagRepository) Save(ctx context.Context, tag *entity.Tag) error {
	return r.db.WithContext(ctx).Save(tag).Error
}

func (r *TagRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&entity.Tag{}, "id = ?", id).Error
}

func (r *TagRepository) GetByID(ctx context.Context, id int64) (*entity.Tag, error) {
	var tag entity.Tag
	if err := r.db.WithContext(ctx).First(&tag, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepository) LockByID(ctx context.Context, id int64) (*entity.Tag, error) {
	var tag entity.Tag
	if err := r.db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(&tag, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepository) List(ctx context.Context, keyword string, page, pageSize int) ([]entity.Tag, int64, error) {
	q := r.db.WithContext(ctx).Model(&entity.Tag{})
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("name ILIKE ? OR slug ILIKE ? OR COALESCE(description, '') ILIKE ?", like, like, like)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var tags []entity.Tag
	if err := q.Order("name ASC, id ASC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&tags).Error; err != nil {
		return nil, 0, err
	}
	return tags, total, nil
}

type ContentTagRepository struct {
	db *gorm.DB
}

func (r *ContentTagRepository) ReplaceForContent(ctx context.Context, contentID int64, tagIDs []int64) error {
	if err := r.db.WithContext(ctx).Delete(&entity.ContentTag{}, "content_id = ?", contentID).Error; err != nil {
		return err
	}
	for _, tagID := range tagIDs {
		if err := r.db.WithContext(ctx).Create(&entity.ContentTag{ContentID: contentID, TagID: tagID}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *ContentTagRepository) ListTags(ctx context.Context, contentID int64) ([]entity.Tag, error) {
	var tags []entity.Tag
	err := r.db.WithContext(ctx).
		Table("content.tags").
		Joins("JOIN content.content_tags ON content.content_tags.tag_id = content.tags.id").
		Where("content.content_tags.content_id = ?", contentID).
		Order("content.tags.name ASC").
		Find(&tags).Error
	return tags, err
}

func (r *ContentTagRepository) CountByTag(ctx context.Context, tagID int64) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&entity.ContentTag{}).Where("tag_id = ?", tagID).Count(&total).Error
	return total, err
}
