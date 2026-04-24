package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListMigrationFilesSortsByFilenameBeforeDirectory(t *testing.T) {
	t.Helper()

	root := t.TempDir()
	migrationsDir := filepath.Join(root, "sql", "migrations", "v3")
	writeMigrationTestFile(t, migrationsDir, filepath.Join("content", "030_v3_content_items.sql"))
	writeMigrationTestFile(t, migrationsDir, filepath.Join("identity", "020_v3_identity_users.sql"))

	files, err := listMigrationFiles(migrationsDir, filepath.Join(root, "sql", "migrations"))
	if err != nil {
		t.Fatalf("list migration files failed: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("unexpected migration count: %d", len(files))
	}
	if files[0].Name != "identity/020_v3_identity_users.sql" {
		t.Fatalf("expected identity migration first, got %s", files[0].Name)
	}
	if files[1].Name != "content/030_v3_content_items.sql" {
		t.Fatalf("expected content migration second, got %s", files[1].Name)
	}
}

func writeMigrationTestFile(t *testing.T, root, relativePath string) {
	t.Helper()

	path := filepath.Join(root, relativePath)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create migration directory failed: %v", err)
	}
	if err := os.WriteFile(path, []byte("SELECT 1;\n"), 0o644); err != nil {
		t.Fatalf("write migration file failed: %v", err)
	}
}
