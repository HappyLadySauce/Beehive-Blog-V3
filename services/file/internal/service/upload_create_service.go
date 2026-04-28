package service

import (
	"context"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/storage"
	"github.com/google/uuid"
)

func (m *Manager) CreateUpload(ctx context.Context, in CreateUploadInput) (*CreateUploadResult, error) {
	if m == nil || m.store == nil || m.storage == nil {
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
	contentType, maxBytes, err := validateUploadFile(m.conf.Storage, scope, in.FileName, in.ContentType, in.ByteSize)
	if err != nil {
		return nil, err
	}

	ttl := time.Duration(m.conf.Storage.PresignTTLSeconds) * time.Second
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	now := time.Now().UTC()
	objectKey := objectKey(scope, ownerUserID, in.FileName, contentType)
	assetID := "asset_" + uuid.NewString()
	publicURL := publicURLForVisibility(m.conf.Storage, visibility, assetID, objectKey)
	asset := &entity.FileAsset{
		AssetID:     assetID,
		UploadID:    "upload_" + uuid.NewString(),
		OwnerUserID: ownerUserID,
		Scope:       scope,
		Visibility:  visibility,
		Status:      StatusPending,
		Bucket:      storageBucket(m.conf.Storage),
		ObjectKey:   objectKey,
		PublicURL:   publicURL,
		FileName:    strings.TrimSpace(in.FileName),
		ContentType: contentType,
		ByteSize:    in.ByteSize,
		ExpiresAt:   now.Add(ttl),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	presign, err := m.storage.PresignPut(ctx, storage.PresignPutInput{
		UploadID:    asset.UploadID,
		Bucket:      asset.Bucket,
		ObjectKey:   asset.ObjectKey,
		ContentType: asset.ContentType,
		ByteSize:    asset.ByteSize,
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
