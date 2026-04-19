package repo

import (
	"context"

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
