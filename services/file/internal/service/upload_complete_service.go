package service

import (
	"context"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/model/repo"
)

func (m *Manager) CompleteUpload(ctx context.Context, actorUserID string, uploadID string) (*AssetView, error) {
	if m == nil || m.store == nil || m.storage == nil {
		return nil, serviceNotInitialized()
	}
	ownerUserID, err := parseActorUserID(actorUserID)
	if err != nil {
		return nil, err
	}
	uploadID = strings.TrimSpace(uploadID)
	if uploadID == "" {
		return nil, invalidArgument("upload_id is required")
	}

	asset, err := m.store.Assets.FindByUploadID(ctx, uploadID)
	if err != nil {
		if repo.IsNotFound(err) {
			return nil, errs.New(errs.CodeFileUploadNotFound, "upload not found")
		}
		return nil, errs.Wrap(err, errs.CodeFileInternal, "get file upload failed")
	}
	if asset.OwnerUserID != ownerUserID {
		return nil, errs.New(errs.CodeFileAccessForbidden, "asset access is forbidden")
	}
	if asset.Status == StatusUploaded {
		return toAssetView(asset), nil
	}
	if asset.Status != StatusPending {
		return nil, errs.New(errs.CodeFileInvalidState, "upload state is invalid")
	}
	if time.Now().UTC().After(asset.ExpiresAt) {
		return nil, errs.New(errs.CodeFileInvalidState, "upload is expired")
	}
	info, err := m.storage.Head(ctx, asset.Bucket, asset.ObjectKey)
	if err != nil {
		return nil, dependencyUnavailable(err)
	}
	headContentType := normalizeContentType(info.ContentType)
	if headContentType == "" {
		headContentType = asset.ContentType
	}
	if info.ByteSize != asset.ByteSize || headContentType != asset.ContentType {
		return nil, errs.New(errs.CodeFileInvalidState, "uploaded object does not match requested file")
	}
	if err := m.storage.Commit(ctx, asset.Bucket, asset.ObjectKey); err != nil {
		return nil, dependencyUnavailable(err)
	}
	updated, err := m.store.Assets.MarkUploaded(ctx, asset.AssetID, time.Now().UTC(), info.ByteSize, headContentType)
	if err != nil {
		return nil, errs.Wrap(err, errs.CodeFileInternal, "complete file upload failed")
	}
	completed, err := m.store.Assets.FindByAssetID(ctx, asset.AssetID)
	if err != nil {
		return nil, errs.Wrap(err, errs.CodeFileInternal, "get completed file asset failed")
	}
	if !updated && completed.Status != StatusUploaded {
		return nil, errs.New(errs.CodeFileInvalidState, "upload state is invalid")
	}
	return toAssetView(completed), nil
}
