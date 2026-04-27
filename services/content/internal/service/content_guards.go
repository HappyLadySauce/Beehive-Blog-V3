package service

import (
	"context"
	"strconv"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/repo"
)

func requireActor(actor Actor) error {
	if actor.UserID <= 0 || normalizeRole(actor.Role) != RoleAdmin {
		return errs.New(errs.CodeContentAccessForbidden, "content access forbidden")
	}
	return nil
}

func normalizeRole(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	return strings.TrimPrefix(normalized, "role_")
}

func parseID(value string, code errs.Code, message string) (int64, error) {
	id, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err != nil || id <= 0 {
		return 0, errs.New(code, message)
	}
	return id, nil
}

func validateTitleSlug(title, slug string) error {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(slug) == "" {
		return errs.New(errs.CodeContentInvalidArgument, "title and slug are required")
	}
	return nil
}

func loadTags(ctx context.Context, store *repo.Store, rawIDs []string) ([]int64, error) {
	result := make([]int64, 0, len(rawIDs))
	seen := map[int64]struct{}{}
	for _, raw := range rawIDs {
		id, err := parseID(raw, errs.CodeContentInvalidArgument, "invalid tag id")
		if err != nil {
			return nil, err
		}
		if _, ok := seen[id]; ok {
			continue
		}
		if _, err := store.Tags.GetByID(ctx, id); err != nil {
			return nil, mapRepoErr(err, errs.CodeContentTagNotFound, "tag not found")
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result, nil
}
