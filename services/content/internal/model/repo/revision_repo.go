package repo

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/entity"
	"gorm.io/gorm"
)

type RevisionRepository struct {
	db *gorm.DB
}

func (r *RevisionRepository) Create(ctx context.Context, revision *entity.Revision) error {
	return r.db.WithContext(ctx).Create(revision).Error
}

func (r *RevisionRepository) GetByID(ctx context.Context, contentID, id int64) (*entity.Revision, error) {
	var revision entity.Revision
	if err := r.db.WithContext(ctx).First(&revision, "content_id = ? AND id = ?", contentID, id).Error; err != nil {
		return nil, err
	}
	return &revision, nil
}

func (r *RevisionRepository) GetCurrent(ctx context.Context, item *entity.Item) (*entity.Revision, error) {
	if item.CurrentRevisionID == nil {
		return nil, gorm.ErrRecordNotFound
	}
	var revision entity.Revision
	if err := r.db.WithContext(ctx).First(&revision, "id = ? AND content_id = ?", *item.CurrentRevisionID, item.ID).Error; err != nil {
		return nil, err
	}
	return &revision, nil
}

func (r *RevisionRepository) NextRevisionNo(ctx context.Context, contentID int64) (int32, error) {
	var maxNo int32
	if err := r.db.WithContext(ctx).Model(&entity.Revision{}).Where("content_id = ?", contentID).Select("COALESCE(MAX(revision_no), 0)").Scan(&maxNo).Error; err != nil {
		return 0, err
	}
	return maxNo + 1, nil
}

func (r *RevisionRepository) List(ctx context.Context, contentID int64, page, pageSize int) ([]entity.Revision, int64, error) {
	q := r.db.WithContext(ctx).Model(&entity.Revision{}).Where("content_id = ?", contentID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var revisions []entity.Revision
	if err := q.Order("revision_no DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&revisions).Error; err != nil {
		return nil, 0, err
	}
	return revisions, total, nil
}
