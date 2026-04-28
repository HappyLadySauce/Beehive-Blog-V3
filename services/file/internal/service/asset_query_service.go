package service

import (
	"context"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/model/repo"
)

func (m *Manager) GetAsset(ctx context.Context, actorUserID string, assetID string) (*AssetView, error) {
	if m == nil || m.store == nil {
		return nil, serviceNotInitialized()
	}
	ownerUserID, err := parseActorUserID(actorUserID)
	if err != nil {
		return nil, err
	}
	assetID = strings.TrimSpace(assetID)
	if assetID == "" {
		return nil, invalidArgument("asset_id is required")
	}
	asset, err := m.store.Assets.FindByAssetID(ctx, assetID)
	if err != nil {
		if repo.IsNotFound(err) {
			return nil, errs.New(errs.CodeFileAssetNotFound, "asset not found")
		}
		return nil, errs.Wrap(err, errs.CodeFileInternal, "get file asset failed")
	}
	if asset.OwnerUserID != ownerUserID {
		return nil, errs.New(errs.CodeFileAccessForbidden, "asset access is forbidden")
	}
	return toAssetView(asset), nil
}
