package service

import (
	"context"
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/model/entity"
)

func TestListAssetsReturnsOnlyCurrentOwnerAssets(t *testing.T) {
	store := newAssetServiceTestStore(t)
	now := time.Now().UTC()
	assets := []*entity.FileAsset{
		{
			AssetID:     "asset-owner-1",
			UploadID:    "upload-owner-1",
			OwnerUserID: 42,
			CategoryKey: "default",
			Visibility:  VisibilityPublic,
			Status:      StatusUploaded,
			Bucket:      "beehive-test",
			ObjectKey:   "content/42/first.png",
			PublicURL:   "https://example.com/first.png",
			FileName:    "first.png",
			ContentType: "image/png",
			ByteSize:    16,
			ExpiresAt:   now.Add(time.Hour),
			UploadedAt:  &now,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			AssetID:     "asset-owner-2",
			UploadID:    "upload-owner-2",
			OwnerUserID: 99,
			CategoryKey: "default",
			Visibility:  VisibilityPublic,
			Status:      StatusUploaded,
			Bucket:      "beehive-test",
			ObjectKey:   "content/99/second.png",
			PublicURL:   "https://example.com/second.png",
			FileName:    "second.png",
			ContentType: "image/png",
			ByteSize:    16,
			ExpiresAt:   now.Add(time.Hour),
			UploadedAt:  &now,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
	for _, asset := range assets {
		if err := store.Assets.Create(context.Background(), asset); err != nil {
			t.Fatalf("failed to create asset %s: %v", asset.AssetID, err)
		}
	}

	manager := NewManager(Dependencies{
		Config:  config.Config{},
		Store:   store,
		Storage: &fakeObjectStorage{},
	})

	result, err := manager.ListAssets(context.Background(), ListAssetsInput{
		ActorUserID: "42",
		Status:      StatusUploaded,
		Page:        1,
		PageSize:    20,
	})
	if err != nil {
		t.Fatalf("expected list assets to succeed, got %v", err)
	}
	if result.Total != 1 {
		t.Fatalf("expected exactly one asset for owner 42, got %d", result.Total)
	}
	if len(result.Items) != 1 || result.Items[0].AssetID != "asset-owner-1" {
		t.Fatalf("unexpected listed assets: %+v", result.Items)
	}
}
