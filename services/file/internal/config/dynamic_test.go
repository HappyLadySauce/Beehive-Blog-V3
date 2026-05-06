package config

import "testing"

func TestConfigCacheApplyUpdatedValueRefreshesSnapshot(t *testing.T) {
	t.Parallel()

	cache := NewConfigCache(nil, StorageConf{
		MaxUploadBytes:    1024,
		PresignTTLSeconds: 300,
	})

	snapshot := cache.applyUpdatedValue("max_upload_bytes", "2048", 2)
	if snapshot.MaxUploadBytes != 2048 {
		t.Fatalf("expected updated max upload bytes, got %d", snapshot.MaxUploadBytes)
	}

	snapshot = cache.applyUpdatedValue("presign_ttl_seconds", "600", 3)
	if snapshot.PresignTTLSeconds != 600 {
		t.Fatalf("expected updated presign ttl seconds, got %d", snapshot.PresignTTLSeconds)
	}
}

func TestConfigCacheRejectsStaleRevisionWriteBack(t *testing.T) {
	t.Parallel()

	cache := NewConfigCache(nil, StorageConf{
		MaxUploadBytes:    1024,
		PresignTTLSeconds: 300,
	})

	newer := cache.applyUpdatedValue("max_upload_bytes", "4096", 10)
	if newer.MaxUploadBytes != 4096 {
		t.Fatalf("expected newer revision to win, got %d", newer.MaxUploadBytes)
	}

	stale := cache.applyUpdatedValue("max_upload_bytes", "2048", 9)
	if stale.MaxUploadBytes != 4096 {
		t.Fatalf("expected stale revision to be ignored, got %d", stale.MaxUploadBytes)
	}

	current := cache.Snapshot()
	if current.MaxUploadBytes != 4096 {
		t.Fatalf("expected cache to keep newer revision value, got %d", current.MaxUploadBytes)
	}
}
