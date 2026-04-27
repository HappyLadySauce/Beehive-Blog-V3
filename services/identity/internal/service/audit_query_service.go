package service

import (
	"context"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/repo"
)

// ListIdentityAudits returns paged audits for active administrators.
// ListIdentityAudits 为活跃管理员返回分页审计记录。
func (s *UserManagementService) ListIdentityAudits(ctx context.Context, in ListIdentityAuditsInput) (*AuditListResult, error) {
	if _, err := s.requireActiveAdmin(ctx, in.ActorUserID); err != nil {
		return nil, err
	}
	if result := strings.TrimSpace(in.Result); result != "" && result != auth.AuditResultSuccess && result != auth.AuditResultFailure {
		return nil, errs.New(errs.CodeIdentityInvalidArgument, "audit result is invalid")
	}
	if len(strings.TrimSpace(in.EventType)) > maxAuditEventLength {
		return nil, errs.New(errs.CodeIdentityInvalidArgument, "event_type must not exceed 64 characters")
	}
	if in.StartedAt != nil && in.EndedAt != nil && in.StartedAt.After(*in.EndedAt) {
		return nil, errs.New(errs.CodeIdentityInvalidArgument, "started_at must be before ended_at")
	}

	page, pageSize := normalizePageInput(in.Page, in.PageSize)
	audits, total, err := s.deps.Store.IdentityAudits.List(ctx, repo.AuditListFilter{
		EventType: strings.TrimSpace(in.EventType),
		Result:    strings.TrimSpace(in.Result),
		UserID:    in.UserID,
		StartedAt: in.StartedAt,
		EndedAt:   in.EndedAt,
		Page:      page,
		PageSize:  pageSize,
	})
	if err != nil {
		return nil, errs.Wrap(err, errs.CodeIdentityInternal, "list identity audits failed")
	}

	return &AuditListResult{Items: audits, Total: total, Page: page, PageSize: pageSize}, nil
}
