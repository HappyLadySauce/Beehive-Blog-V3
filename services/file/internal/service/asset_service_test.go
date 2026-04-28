package service

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/model/repo"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/storage"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
		Config:        config.Config{ObjectStorage: config.ObjectStorageConf{Bucket: "beehive-test"}},
		Store:         store,
		ObjectStorage: objectStorage,
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

type fakeObjectStorage struct {
	deleteErrors []error
	deleteCalls  int
}

func (s *fakeObjectStorage) PresignPut(context.Context, storage.PresignPutInput) (*storage.PresignPutOutput, error) {
	return &storage.PresignPutOutput{UploadURL: "https://upload.example.com", Headers: map[string]string{}}, nil
}

func (s *fakeObjectStorage) Head(context.Context, string, string) (*storage.ObjectInfo, error) {
	return &storage.ObjectInfo{ByteSize: 1, ContentType: "text/plain"}, nil
}

func (s *fakeObjectStorage) Delete(context.Context, string, string) error {
	s.deleteCalls++
	if len(s.deleteErrors) == 0 {
		return nil
	}
	err := s.deleteErrors[0]
	s.deleteErrors = s.deleteErrors[1:]
	return err
}

func (s *fakeObjectStorage) Health(context.Context) error {
	return nil
}

func newAssetServiceTestStore(t *testing.T) *repo.Store {
	t.Helper()

	ctx := context.Background()
	container, err := tcpostgres.Run(
		ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("beehive_file_test"),
		tcpostgres.WithUsername("postgres"),
		tcpostgres.WithPassword("postgres"),
	)
	if err != nil {
		t.Skipf("skip PostgreSQL integration test: %v", err)
	}
	t.Cleanup(func() {
		_ = container.Terminate(context.Background())
	})

	dsn, err := container.ConnectionString(ctx, "sslmode=disable", "TimeZone=Asia/Shanghai")
	if err != nil {
		t.Fatalf("failed to build PostgreSQL DSN: %v", err)
	}
	db, sqlDB := openPostgresWithRetry(t, ctx, dsn)
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	migration, err := os.ReadFile(filepath.Join(repoRootDir(), "sql", "migrations", "v3", "file", "040_v3_file_assets.sql"))
	if err != nil {
		t.Fatalf("failed to read file migration: %v", err)
	}
	if _, err := sqlDB.ExecContext(ctx, string(migration)); err != nil {
		t.Fatalf("failed to apply file migration: %v", err)
	}

	return repo.NewStore(db)
}

func openPostgresWithRetry(t *testing.T, ctx context.Context, dsn string) (*gorm.DB, *sql.DB) {
	t.Helper()

	var lastErr error
	for attempt := 0; attempt < 20; attempt++ {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			lastErr = err
			time.Sleep(500 * time.Millisecond)
			continue
		}
		sqlDB, err := db.DB()
		if err != nil {
			lastErr = err
			time.Sleep(500 * time.Millisecond)
			continue
		}
		if err := sqlDB.PingContext(ctx); err != nil {
			lastErr = err
			_ = sqlDB.Close()
			time.Sleep(500 * time.Millisecond)
			continue
		}
		return db, sqlDB
	}

	t.Fatalf("failed to open PostgreSQL test DB after retries: %v", lastErr)
	return nil, nil
}

func repoRootDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "..", ".."))
}
