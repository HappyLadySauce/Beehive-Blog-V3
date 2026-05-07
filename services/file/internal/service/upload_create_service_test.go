package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/model/repo"
)

func TestCreateUploadRejectsMismatchedContentType(t *testing.T) {
	t.Parallel()

	store := newAssetServiceTestStore(t)
	createTestCategory(t, store, "images", []string{".png"})
	objectStorage := &fakeObjectStorage{}
	manager := NewManager(Dependencies{
		Config: config.Config{
			Storage: config.StorageConf{
				MaxUploadBytes:    1024,
				PresignTTLSeconds: 300,
			},
		},
		Store:       store,
		Storage:     objectStorage,
		ConfigCache: config.NewConfigCache(nil, config.StorageConf{MaxUploadBytes: 1024, PresignTTLSeconds: 300}),
	})

	_, err := manager.CreateUpload(context.Background(), CreateUploadInput{
		ActorUserID: "42",
		CategoryKey: "images",
		FileName:    "avatar.png",
		ContentType: "text/html",
		ByteSize:    128,
		Visibility:  VisibilityPublic,
	})
	if !errors.Is(err, errs.E(errs.CodeFileInvalidArgument)) {
		t.Fatalf("expected invalid argument for mismatched content type, got %v", err)
	}
	if objectStorage.presignCalls != 0 {
		t.Fatalf("expected presign to be skipped, got %d calls", objectStorage.presignCalls)
	}
}

func TestCreateUploadInfersAndPersistsContentTypeForMatchingExtension(t *testing.T) {
	t.Parallel()

	store := newAssetServiceTestStore(t)
	createTestCategory(t, store, "images", []string{".png"})
	objectStorage := &fakeObjectStorage{}
	manager := NewManager(Dependencies{
		Config: config.Config{
			Storage: config.StorageConf{
				MaxUploadBytes:    1024,
				PresignTTLSeconds: 300,
			},
		},
		Store:       store,
		Storage:     objectStorage,
		ConfigCache: config.NewConfigCache(nil, config.StorageConf{MaxUploadBytes: 1024, PresignTTLSeconds: 300}),
	})

	result, err := manager.CreateUpload(context.Background(), CreateUploadInput{
		ActorUserID: "42",
		CategoryKey: "images",
		FileName:    "avatar.png",
		ContentType: "",
		ByteSize:    128,
		Visibility:  VisibilityPublic,
	})
	if err != nil {
		t.Fatalf("expected create upload to succeed, got %v", err)
	}
	if result.Asset.ContentType != "image/png" {
		t.Fatalf("expected persisted content type image/png, got %s", result.Asset.ContentType)
	}
	if objectStorage.presignCalls != 1 {
		t.Fatalf("expected one presign call, got %d", objectStorage.presignCalls)
	}
	if objectStorage.lastPresign.ContentType != "image/png" {
		t.Fatalf("expected presigned content type image/png, got %s", objectStorage.lastPresign.ContentType)
	}
}

func createTestCategory(t *testing.T, store *repo.Store, categoryKey string, allowedExtensions []string) {
	t.Helper()

	now := time.Now().UTC()
	if err := store.Categories.Create(context.Background(), &entity.FileCategory{
		CategoryKey: categoryKey,
		DisplayName: categoryKey,
		Description: "test category",
		Enabled:     true,
		IsDefault:   false,
		SortOrder:   10,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, allowedExtensions); err != nil {
		t.Fatalf("failed to create test category %s: %v", categoryKey, err)
	}
}
