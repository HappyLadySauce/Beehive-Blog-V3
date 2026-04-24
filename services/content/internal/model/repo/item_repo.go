package repo

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ItemRepository struct {
	db *gorm.DB
}

func (r *ItemRepository) Create(ctx context.Context, item *entity.Item) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *ItemRepository) Save(ctx context.Context, item *entity.Item) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *ItemRepository) GetByID(ctx context.Context, id int64) (*entity.Item, error) {
	var item entity.Item
	if err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ItemRepository) GetBySlug(ctx context.Context, slug string) (*entity.Item, error) {
	var item entity.Item
	if err := r.db.WithContext(ctx).First(&item, "slug = ?", slug).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ItemRepository) LockByID(ctx context.Context, id int64) (*entity.Item, error) {
	var item entity.Item
	if err := r.db.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

type ItemFilter struct {
	PublicOnly bool
	Type       string
	Status     string
	Visibility string
	Keyword    string
}

func (r *ItemRepository) List(ctx context.Context, filter ItemFilter, page, pageSize int) ([]entity.Item, int64, error) {
	q := r.db.WithContext(ctx).Model(&entity.Item{})
	if filter.PublicOnly {
		q = q.Where("status = ? AND visibility = ?", "published", "public")
	}
	if filter.Type != "" {
		q = q.Where("type = ?", filter.Type)
	}
	if filter.Status != "" {
		q = q.Where("status = ?", filter.Status)
	}
	if filter.Visibility != "" {
		q = q.Where("visibility = ?", filter.Visibility)
	}
	if filter.Keyword != "" {
		like := "%" + filter.Keyword + "%"
		q = q.Where("title ILIKE ? OR slug ILIKE ? OR COALESCE(summary, '') ILIKE ?", like, like, like)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []entity.Item
	if err := q.Order("sort_order ASC, published_at DESC NULLS LAST, id DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
