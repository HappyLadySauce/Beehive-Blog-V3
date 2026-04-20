package testkit

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresState struct {
	container *tcpostgres.PostgresContainer
	dsn       string
	db        *gorm.DB
	sqlDB     *sql.DB
	err       error
}

var (
	postgresOnce   sync.Once
	sharedPostgres postgresState
)

// SharedPostgres returns the shared PostgreSQL state for integration tests.
// SharedPostgres 返回集成测试共享的 PostgreSQL 状态。
func SharedPostgres(t *testing.T) *postgresState {
	t.Helper()

	postgresOnce.Do(func() {
		ctx := t.Context()
		container, err := tcpostgres.Run(
			ctx,
			"postgres:16-alpine",
			tcpostgres.WithDatabase("beehive_identity_test"),
			tcpostgres.WithUsername("postgres"),
			tcpostgres.WithPassword("postgres"),
		)
		if err != nil {
			sharedPostgres.err = fallbackPostgres(ctx, err)
			return
		}

		dsn, err := container.ConnectionString(ctx, "sslmode=disable", "TimeZone=Asia/Shanghai")
		if err != nil {
			sharedPostgres.err = err
			return
		}

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			sharedPostgres.err = err
			return
		}

		sqlDB, err := db.DB()
		if err != nil {
			sharedPostgres.err = err
			return
		}

		if err := runMigrations(ctx, sqlDB); err != nil {
			sharedPostgres.err = err
			return
		}

		sharedPostgres.container = container
		sharedPostgres.dsn = dsn
		sharedPostgres.db = db
		sharedPostgres.sqlDB = sqlDB
	})

	if sharedPostgres.err != nil {
		t.Skipf("skip PostgreSQL integration test: %v", sharedPostgres.err)
	}

	return &sharedPostgres
}

// ResetPostgres truncates all identity tables before a test.
// ResetPostgres 在测试前清空所有 identity 表。
func ResetPostgres(t *testing.T) {
	t.Helper()

	state := SharedPostgres(t)
	if state.sqlDB == nil {
		t.Fatalf("postgres sql DB is not initialized")
	}

	const truncateSQL = `
TRUNCATE TABLE
  identity.identity_audits,
  identity.refresh_tokens,
  identity.user_sessions,
  identity.oauth_login_states,
  identity.federated_identities,
  identity.credential_locals,
  identity.users
RESTART IDENTITY CASCADE;
`
	if _, err := state.sqlDB.ExecContext(context.Background(), truncateSQL); err != nil {
		t.Fatalf("failed to reset identity tables: %v", err)
	}
}

// PostgresDB returns the shared gorm DB for integration tests.
// PostgresDB 返回集成测试共享的 gorm DB。
func PostgresDB(t *testing.T) *gorm.DB {
	t.Helper()
	return SharedPostgres(t).db
}

func fallbackPostgres(ctx context.Context, containerErr error) error {
	cfg, ok, err := loadPostgresEnv()
	if err != nil {
		return fmt.Errorf("load PostgreSQL fallback environment failed after container error %v: %w", containerErr, err)
	}
	if !ok {
		return containerErr
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
		cfg.TimeZone,
	)

	db, openErr := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if openErr != nil {
		return fmt.Errorf("start PostgreSQL container failed: %v; fallback PostgreSQL open failed: %w", containerErr, openErr)
	}

	sqlDB, dbErr := db.DB()
	if dbErr != nil {
		return fmt.Errorf("start PostgreSQL container failed: %v; fallback PostgreSQL sql DB failed: %w", containerErr, dbErr)
	}
	if pingErr := sqlDB.PingContext(ctx); pingErr != nil {
		return fmt.Errorf("start PostgreSQL container failed: %v; fallback PostgreSQL ping failed: %w", containerErr, pingErr)
	}
	if migrateErr := runMigrations(ctx, sqlDB); migrateErr != nil {
		return fmt.Errorf("start PostgreSQL container failed: %v; fallback PostgreSQL migration failed: %w", containerErr, migrateErr)
	}

	sharedPostgres.dsn = dsn
	sharedPostgres.db = db
	sharedPostgres.sqlDB = sqlDB

	return nil
}

func runMigrations(ctx context.Context, db *sql.DB) error {
	files, err := filepath.Glob(filepath.Join(repoRootDir(), "sql", "migrations", "v3", "identity", "*.sql"))
	if err != nil {
		return fmt.Errorf("glob identity migrations failed: %w", err)
	}
	for _, path := range files {
		content, err := osReadFile(path)
		if err != nil {
			return fmt.Errorf("read migration %s failed: %w", path, err)
		}
		if _, err := db.ExecContext(ctx, string(content)); err != nil {
			return fmt.Errorf("apply migration %s failed: %w", path, err)
		}
	}

	return nil
}

func repoRootDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "..", ".."))
}
