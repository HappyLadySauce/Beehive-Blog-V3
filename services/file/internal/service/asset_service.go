package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/model/repo"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/storage"
	"github.com/google/uuid"
)

func (m *Manager) CreateUpload(ctx context.Context, in CreateUploadInput) (*CreateUploadResult, error) {
	if m == nil || m.store == nil || m.objectStorage == nil {
		return nil, serviceNotInitialized()
	}
	ownerUserID, err := parseActorUserID(in.ActorUserID)
	if err != nil {
		return nil, err
	}
	scope, err := normalizeScope(in.Scope)
	if err != nil {
		return nil, err
	}
	visibility, err := normalizeVisibility(in.Visibility)
	if err != nil {
		return nil, err
	}
	contentType, maxBytes, err := validateUploadFile(m.conf.ObjectStorage, scope, in.FileName, in.ContentType, in.ByteSize)
	if err != nil {
		return nil, err
	}

	ttl := time.Duration(m.conf.ObjectStorage.PresignTTLSeconds) * time.Second
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	now := time.Now().UTC()
	objectKey := objectKey(scope, ownerUserID, in.FileName, contentType)
	publicURL := publicURLForVisibility(m.conf.ObjectStorage.PublicBaseURL, visibility, objectKey)
	asset := &entity.FileAsset{
		AssetID:     "asset_" + uuid.NewString(),
		UploadID:    "upload_" + uuid.NewString(),
		OwnerUserID: ownerUserID,
		Scope:       scope,
		Visibility:  visibility,
		Status:      StatusPending,
		Bucket:      strings.TrimSpace(m.conf.ObjectStorage.Bucket),
		ObjectKey:   objectKey,
		PublicURL:   publicURL,
		FileName:    strings.TrimSpace(in.FileName),
		ContentType: contentType,
		ByteSize:    in.ByteSize,
		ExpiresAt:   now.Add(ttl),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	presign, err := m.objectStorage.PresignPut(ctx, storage.PresignPutInput{
		Bucket:      asset.Bucket,
		ObjectKey:   asset.ObjectKey,
		ContentType: asset.ContentType,
		Expires:     ttl,
	})
	if err != nil {
		return nil, dependencyUnavailable(err)
	}
	if err := m.store.Assets.Create(ctx, asset); err != nil {
		return nil, errs.Wrap(err, errs.CodeFileInternal, "create file asset failed")
	}
	return &CreateUploadResult{
		Asset:     toAssetView(asset),
		UploadURL: presign.UploadURL,
		Headers:   presign.Headers,
		ExpiresAt: asset.ExpiresAt,
		MaxBytes:  maxBytes,
	}, nil
}

func (m *Manager) CompleteUpload(ctx context.Context, actorUserID string, uploadID string) (*AssetView, error) {
	if m == nil || m.store == nil || m.objectStorage == nil {
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
	info, err := m.objectStorage.Head(ctx, asset.Bucket, asset.ObjectKey)
	if err != nil {
		return nil, dependencyUnavailable(err)
	}
	headContentType := normalizeContentType(info.ContentType)
	if info.ByteSize != asset.ByteSize || headContentType != asset.ContentType {
		return nil, errs.New(errs.CodeFileInvalidState, "uploaded object does not match requested file")
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

func (m *Manager) DeleteAsset(ctx context.Context, actorUserID string, assetID string) error {
	if m == nil || m.store == nil || m.objectStorage == nil {
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
	if err := m.objectStorage.Delete(ctx, asset.Bucket, asset.ObjectKey); err != nil {
		return dependencyUnavailable(err)
	}
	return nil
}

func objectKey(scope string, ownerUserID int64, fileName string, contentType string) string {
	prefix := map[string]string{
		ScopeAvatar:       "avatars",
		ScopeContentCover: "content/covers",
		ScopeContentImage: "content/images",
		ScopeAttachment:   "attachments",
	}[scope]
	return prefix + "/" + strconv.FormatInt(ownerUserID, 10) + "/" + uuid.NewString() + extensionFor(fileName, contentType)
}

func publicURLForVisibility(publicBaseURL string, visibility string, objectKey string) string {
	if visibility != VisibilityPublic {
		return ""
	}
	baseURL := strings.TrimRight(strings.TrimSpace(publicBaseURL), "/")
	objectKey = strings.TrimLeft(strings.TrimSpace(objectKey), "/")
	if baseURL == "" || objectKey == "" {
		return ""
	}
	return baseURL + "/" + objectKey
}
