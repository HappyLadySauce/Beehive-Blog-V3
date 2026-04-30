package storage

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/config"
)

func TestLocalStorageLifecycle(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	temp := t.TempDir()
	store, err := NewLocalStorage(config.LocalStorageConf{
		RootDir: root,
		TempDir: temp,
		Bucket:  "local-test",
	})
	if err != nil {
		t.Fatalf("expected local storage to initialize, got %v", err)
	}

	objectKey := "avatars/42/avatar.png"

	// PresignPut returns ErrStorageDisabled for local storage without HTTP server.
	_, err = store.PresignPut(context.Background(), PresignPutInput{
		UploadID:    "upload_1",
		Bucket:      "local-test",
		ObjectKey:   objectKey,
		ContentType: "image/png",
		ByteSize:    4,
	})
	if err != ErrStorageDisabled {
		t.Fatalf("expected PresignPut to return ErrStorageDisabled, got %v", err)
	}

	// Write file to temp dir manually (simulating external upload).
	pendingPath := filepath.Join(temp, objectKey)
	if err := os.MkdirAll(filepath.Dir(pendingPath), 0o755); err != nil {
		t.Fatalf("expected mkdir for pending to pass, got %v", err)
	}
	if err := os.WriteFile(pendingPath, []byte("data"), 0o644); err != nil {
		t.Fatalf("expected pending file write to pass, got %v", err)
	}

	// Head finds file in pending dir.
	info, err := store.Head(context.Background(), "local-test", objectKey)
	if err != nil {
		t.Fatalf("expected head of pending file to pass, got %v", err)
	}
	if info.ByteSize != 4 {
		t.Fatalf("expected 4 bytes, got %d", info.ByteSize)
	}

	// Commit moves from temp to root.
	if err := store.Commit(context.Background(), "local-test", objectKey); err != nil {
		t.Fatalf("expected commit to pass, got %v", err)
	}

	// File should now be in root.
	uploadedPath := filepath.Join(root, objectKey)
	if _, err := os.Stat(uploadedPath); err != nil {
		t.Fatalf("expected committed file in root, got %v", err)
	}

	// Commit is idempotent.
	if err := store.Commit(context.Background(), "local-test", objectKey); err != nil {
		t.Fatalf("expected second commit to pass, got %v", err)
	}

	// Delete removes from both dirs.
	if err := store.Delete(context.Background(), "local-test", objectKey); err != nil {
		t.Fatalf("expected delete to pass, got %v", err)
	}
	if _, err := os.Stat(uploadedPath); !os.IsNotExist(err) {
		t.Fatal("expected deleted object to be unavailable")
	}
}

func TestLocalStorageHealth(t *testing.T) {
	t.Parallel()

	store, err := NewLocalStorage(config.LocalStorageConf{
		RootDir: t.TempDir(),
		TempDir: t.TempDir(),
		Bucket:  "local-test",
	})
	if err != nil {
		t.Fatalf("expected local storage to initialize, got %v", err)
	}

	if err := store.Health(context.Background()); err != nil {
		t.Fatalf("expected health to pass, got %v", err)
	}
}
