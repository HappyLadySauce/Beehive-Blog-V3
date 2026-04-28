package storage

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/url"
	"strings"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/config"
)

func TestLocalStorageLifecycle(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	temp := t.TempDir()
	store, err := NewLocalStorage(config.LocalStorageConf{
		RootDir:       root,
		TempDir:       temp,
		Bucket:        "local-test",
		UploadBaseURL: "http://127.0.0.1:8084/files/uploads",
		UploadSecret:  "test-secret",
	})
	if err != nil {
		t.Fatalf("expected local storage to initialize, got %v", err)
	}

	presign, err := store.PresignPut(context.Background(), PresignPutInput{
		UploadID:    "upload_1",
		Bucket:      "local-test",
		ObjectKey:   "avatars/42/avatar.png",
		ContentType: "image/png",
		ByteSize:    4,
	})
	if err != nil {
		t.Fatalf("expected presign to pass, got %v", err)
	}
	parsedUploadURL, err := url.Parse(presign.UploadURL)
	if err != nil {
		t.Fatalf("expected upload URL to parse, got %v", err)
	}
	if !strings.HasSuffix(parsedUploadURL.Path, "/upload_1") {
		t.Fatalf("unexpected upload URL path: %s", presign.UploadURL)
	}
	if parsedUploadURL.RawQuery != "" {
		t.Fatalf("expected upload URL not to expose query credentials, got %q", parsedUploadURL.RawQuery)
	}
	token := presign.Headers[UploadTokenHeader]
	if token == "" {
		t.Fatal("expected presigned local upload token header")
	}
	if !store.VerifyUploadToken("upload_1", "avatars/42/avatar.png", token) {
		t.Fatal("expected presigned local upload token to verify")
	}
	if store.VerifyUploadToken("upload_1", "avatars/42/other.png", token) {
		t.Fatal("expected token to be bound to the object key")
	}

	info, err := store.PutPending(context.Background(), "avatars/42/avatar.png", bytes.NewBufferString("data"), 4)
	if err != nil {
		t.Fatalf("expected pending write to pass, got %v", err)
	}
	if info.ByteSize != 4 {
		t.Fatalf("expected 4 bytes, got %d", info.ByteSize)
	}
	if err := store.Commit(context.Background(), "local-test", "avatars/42/avatar.png"); err != nil {
		t.Fatalf("expected commit to pass, got %v", err)
	}

	reader, uploaded, err := store.OpenUploaded(context.Background(), "avatars/42/avatar.png")
	if err != nil {
		t.Fatalf("expected uploaded object to open, got %v", err)
	}
	body, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("expected uploaded object to read, got %v", err)
	}
	if string(body) != "data" || uploaded.ByteSize != 4 {
		t.Fatalf("unexpected uploaded object: body=%q info=%+v", string(body), uploaded)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("expected uploaded reader to close, got %v", err)
	}

	if err := store.Delete(context.Background(), "local-test", "avatars/42/avatar.png"); err != nil {
		t.Fatalf("expected delete to pass, got %v", err)
	}
	if _, _, err := store.OpenUploaded(context.Background(), "avatars/42/avatar.png"); err == nil {
		t.Fatal("expected deleted object to be unavailable")
	}
}

func TestLocalStorageRejectsTraversal(t *testing.T) {
	t.Parallel()

	store, err := NewLocalStorage(config.LocalStorageConf{
		RootDir:       t.TempDir(),
		TempDir:       t.TempDir(),
		Bucket:        "local-test",
		UploadBaseURL: "http://127.0.0.1:8084/files/uploads",
		UploadSecret:  "test-secret",
	})
	if err != nil {
		t.Fatalf("expected local storage to initialize, got %v", err)
	}

	if _, err := store.PutPending(context.Background(), "../escape.png", bytes.NewBufferString("data"), 4); !errors.Is(err, ErrStorageInvalidInput) {
		t.Fatalf("expected traversal to be rejected, got %v", err)
	}
}

func TestLocalStorageRejectsOversizedWrite(t *testing.T) {
	t.Parallel()

	store, err := NewLocalStorage(config.LocalStorageConf{
		RootDir:       t.TempDir(),
		TempDir:       t.TempDir(),
		Bucket:        "local-test",
		UploadBaseURL: "http://127.0.0.1:8084/files/uploads",
		UploadSecret:  "test-secret",
	})
	if err != nil {
		t.Fatalf("expected local storage to initialize, got %v", err)
	}

	if _, err := store.PutPending(context.Background(), "attachments/42/file.txt", bytes.NewBufferString("toolarge"), 4); !errors.Is(err, ErrStorageObjectTooLarge) {
		t.Fatalf("expected oversized write to be rejected, got %v", err)
	}
}
