package service

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/repo"
)

// ListUsers returns paged users for active administrators.
// ListUsers 为活跃管理员返回分页用户列表。
func (s *UserManagementService) ListUsers(ctx context.Context, in ListUsersInput) (*UserListResult, error) {
	if _, err := s.requireActiveAdmin(ctx, in.ActorUserID); err != nil {
		return nil, err
	}

	role, err := normalizeOptionalRole(in.Role)
	if err != nil {
		return nil, err
	}
	status, err := normalizeOptionalStatus(in.Status, true)
	if err != nil {
		return nil, err
	}

	page, pageSize := normalizePageInput(in.Page, in.PageSize)
	users, total, err := s.deps.Store.Users.List(ctx, repo.ListFilter{
		Keyword:        in.Keyword,
		Role:           role,
		Status:         status,
		IncludeDeleted: in.IncludeDeleted,
		Page:           page,
		PageSize:       pageSize,
	})
	if err != nil {
		return nil, errs.Wrap(err, errs.CodeIdentityInternal, "list users failed")
	}

	return &UserListResult{Items: users, Total: total, Page: page, PageSize: pageSize}, nil
}
