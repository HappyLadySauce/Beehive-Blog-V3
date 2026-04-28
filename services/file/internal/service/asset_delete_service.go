package service

import (
	"context"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/model/repo"
)

func (m *Manager) DeleteAsset(ctx context.Context, actorUserID string, assetID string) error {
	if m == nil || m.store == nil || m.storage == nil {
		return serviceNotInitialized()
	}
	ownerUserID, err := parseActorUserID(actorUserID)
	if err != nil {
		return err
	}
	assetID = strings.TrimSpace(assetID)
	if assetID == "" {
		return invalidArgument("asset_id is required")
	}
	asset, err := m.store.Assets.FindByAssetID(ctx, assetID)
	if err != nil {
		if repo.IsNotFound(err) {
			return errs.New(errs.CodeFileAssetNotFound, "asset not found")
		}
		return errs.Wrap(err, errs.CodeFileInternal, "get file asset failed")
	}
	if asset.OwnerUserID != ownerUserID {
		return errs.New(errs.CodeFileAccessForbidden, "asset access is forbidden")
	}
	if asset.Status != StatusDeleted {
		if err := m.store.Assets.MarkDeleted(ctx, asset.AssetID, time.Now().UTC()); err != nil {
			return errs.Wrap(err, errs.CodeFileInternal, "mark file asset deleted failed")
		}
	}
	if err := m.storage.Delete(ctx, asset.Bucket, asset.ObjectKey); err != nil {
		return dependencyUnavailable(err)
	}
	return nil
}
