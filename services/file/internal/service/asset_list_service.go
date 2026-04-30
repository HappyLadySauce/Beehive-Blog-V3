package service

import (
	"context"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/model/repo"
)

func (m *Manager) ListAssets(ctx context.Context, in ListAssetsInput) (*AssetListResult, error) {
	if m == nil || m.store == nil {
		return nil, serviceNotInitialized()
	}

	ownerUserID, err := parseActorUserID(in.ActorUserID)
	if err != nil {
		return nil, err
	}
	if requestedOwner := strings.TrimSpace(in.OwnerUserID); requestedOwner != "" && requestedOwner != in.ActorUserID {
		return nil, errs.New(errs.CodeFileAccessForbidden, "owner_user_id is forbidden")
	}

	namespace, err := normalizeOptionalNamespace(in.Namespace)
	if err != nil {
		return nil, err
	}
	status, err := normalizeOptionalStatus(in.Status)
	if err != nil {
		return nil, err
	}
	visibility, err := normalizeOptionalVisibility(in.Visibility)
	if err != nil {
		return nil, err
	}

	assets, total, err := m.store.Assets.List(ctx, repoAssetFilter(ownerUserID, namespace, status, visibility, in.Keyword, in.Page, in.PageSize))
	if err != nil {
		return nil, errs.Wrap(err, errs.CodeFileInternal, "list file assets failed")
	}

	page, pageSize := normalizeListPage(in.Page, in.PageSize)
	items := make([]*AssetView, 0, len(assets))
	for i := range assets {
		items = append(items, toAssetView(&assets[i]))
	}
	return &AssetListResult{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func repoAssetFilter(ownerUserID int64, namespace string, status string, visibility string, keyword string, page int, pageSize int) repo.AssetListFilter {
	return repo.AssetListFilter{
		OwnerUserID: ownerUserID,
		Namespace:   namespace,
		Status:      status,
		Visibility:  visibility,
		Keyword:     strings.TrimSpace(keyword),
		Page:        page,
		PageSize:    pageSize,
	}
}

func normalizeListPage(page int, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	switch {
	case pageSize <= 0:
		pageSize = 20
	case pageSize > 100:
		pageSize = 100
	}
	return page, pageSize
}
