package repo

import (
	"context"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"gorm.io/gorm"
)

// IdentityAuditRepository encapsulates audit persistence.
// IdentityAuditRepository 封装审计记录持久化访问。
type IdentityAuditRepository struct {
	db *gorm.DB
}

// Create persists an audit record.
// Create 持久化审计记录。
func (r *IdentityAuditRepository) Create(ctx context.Context, audit *entity.IdentityAudit) error {
	return r.db.WithContext(ctx).Create(audit).Error
}

// AuditListFilter describes audit-list filtering and pagination.
// AuditListFilter 描述审计列表过滤与分页参数。
type AuditListFilter struct {
	EventType string
	Result    string
	UserID    *int64
	StartedAt *time.Time
	EndedAt   *time.Time
	Page      int
	PageSize  int
}

// List returns audit rows matching the filter and the total count.
// List 返回符合过滤条件的审计记录与总数。
func (r *IdentityAuditRepository) List(ctx context.Context, filter AuditListFilter) ([]entity.IdentityAudit, int64, error) {
	page, pageSize := normalizePagination(filter.Page, filter.PageSize)
	query := r.db.WithContext(ctx).Model(&entity.IdentityAudit{})

	if filter.EventType != "" {
		query = query.Where("event_type = ?", filter.EventType)
	}
	if filter.Result != "" {
		query = query.Where("result = ?", filter.Result)
	}
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.StartedAt != nil {
		query = query.Where("created_at >= ?", *filter.StartedAt)
	}
	if filter.EndedAt != nil {
		query = query.Where("created_at <= ?", *filter.EndedAt)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var audits []entity.IdentityAudit
	if err := query.
		Order("created_at DESC, id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&audits).Error; err != nil {
		return nil, 0, err
	}

	return audits, total, nil
}
