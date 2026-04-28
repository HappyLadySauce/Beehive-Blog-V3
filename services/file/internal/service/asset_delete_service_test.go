package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/model/entity"
)

func TestDeleteAssetAllowsRetryAfterObjectDeleteFailure(t *testing.T) {
	store := newAssetServiceTestStore(t)
	asset := &entity.FileAsset{
		AssetID:     "asset-delete-retry",
		UploadID:    "upload-delete-retry",
		OwnerUserID: 42,
		Scope:       ScopeAttachment,
		Visibility:  VisibilityPrivate,
		Status:      StatusUploaded,
		Bucket:      "beehive-test",
		ObjectKey:   "attachments/42/delete-retry.txt",
		PublicURL:   "",
		FileName:    "delete-retry.txt",
		ContentType: "text/plain",
		ByteSize:    16,
		ExpiresAt:   time.Now().UTC().Add(time.Hour),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	if err := store.Assets.Create(context.Background(), asset); err != nil {
		t.Fatalf("failed to create asset: %v", err)
	}

	objectStorage := &fakeObjectStorage{
		deleteErrors: []error{errors.New("object storage unavailable"), nil},
	}
	manager := NewManager(Dependencies{
		Config:  config.Config{Storage: config.StorageConf{Driver: "s3", S3: config.S3StorageConf{Bucket: "beehive-test"}}},
		Store:   store,
		Storage: objectStorage,
	})

	if err := manager.DeleteAsset(context.Background(), "42", asset.AssetID); !errors.Is(err, errs.E(errs.CodeFileDependencyUnavailable)) {
		t.Fatalf("expected object delete failure to be dependency unavailable, got %v", err)
	}
	afterFailure, err := store.Assets.FindByAssetID(context.Background(), asset.AssetID)
	if err != nil {
		t.Fatalf("failed to reload asset after failed delete: %v", err)
	}
	if afterFailure.Status != StatusDeleted || afterFailure.DeletedAt == nil {
		t.Fatalf("expected database marker to be deleted after object failure, got %+v", afterFailure)
	}

	if err := manager.DeleteAsset(context.Background(), "42", asset.AssetID); err != nil {
		t.Fatalf("expected retry delete to succeed, got %v", err)
	}
	if objectStorage.deleteCalls != 2 {
		t.Fatalf("expected deleted asset retry to call object delete again, got %d calls", objectStorage.deleteCalls)
	}
}
