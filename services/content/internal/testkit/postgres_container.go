package testkit

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/model/repo"
	contentservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/content/internal/service"
	_ "github.com/jackc/pgx/v5/stdlib"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresState struct {
	container *tcpostgres.PostgresContainer
	adminDSN  string
	adminDB   *sql.DB
	err       error
}

var (
	postgresOnce   sync.Once
	sharedPostgres postgresState
	databaseSeq    uint64
)

func NewServiceDependencies(t *testing.T) contentservice.Dependencies {
	t.Helper()

	db, sqlDB := newIsolatedDB(t)
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	fixedNow := time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)
	return contentservice.Dependencies{
		Config: config.Config{
			Postgres: config.PostgresConf{
				Host:   "test",
				User:   "postgres",
				DBName: "content_test",
			},
			InternalAuthToken: "secret",
			AllowedCallers:    []string{"gateway"},
		},
		Store: repo.NewStore(db),
		Clock: func() time.Time {
			return fixedNow
		},
		CheckReadiness: func(ctx context.Context) error {
			return sqlDB.PingContext(ctx)
		},
	}
}

func newIsolatedDB(t *testing.T) (*gorm.DB, *sql.DB) {
	t.Helper()

	state := SharedPostgres(t)
	dbName := fmt.Sprintf("content_test_%d_%d", time.Now().UnixNano(), atomic.AddUint64(&databaseSeq, 1))
	if _, err := state.adminDB.ExecContext(t.Context(), "CREATE DATABASE "+quoteIdent(dbName)); err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	t.Cleanup(func() {
		_, _ = state.adminDB.ExecContext(context.Background(), "DROP DATABASE IF EXISTS "+quoteIdent(dbName)+" WITH (FORCE)")
	})

	dsn, err := dsnForDatabase(state.adminDSN, dbName)
	if err != nil {
		t.Fatalf("failed to build test database dsn: %v", err)
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get sql database: %v", err)
	}
	if err := runMigrations(t.Context(), sqlDB); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}
	if err := seedUsers(t.Context(), sqlDB); err != nil {
		t.Fatalf("failed to seed users: %v", err)
	}

	return db, sqlDB
}

func SharedPostgres(t *testing.T) *postgresState {
	t.Helper()

	postgresOnce.Do(func() {
		ctx := t.Context()
		container, err := tcpostgres.Run(
			ctx,
			"postgres:16-alpine",
			tcpostgres.WithDatabase("beehive_content_test"),
			tcpostgres.WithUsername("postgres"),
			tcpostgres.WithPassword("postgres"),
		)
		if err != nil {
			sharedPostgres.err = err
			return
		}

		dsn, err := container.ConnectionString(ctx, "sslmode=disable", "TimeZone=Asia/Shanghai")
		if err != nil {
			sharedPostgres.err = err
			return
		}
		adminDB, err := sql.Open("pgx", dsn)
		if err != nil {
			sharedPostgres.err = err
			return
		}
		if err := pingWithRetry(ctx, adminDB); err != nil {
			sharedPostgres.err = err
			return
		}

		sharedPostgres.container = container
		sharedPostgres.adminDSN = dsn
		sharedPostgres.adminDB = adminDB
	})

	if sharedPostgres.err != nil {
		t.Skipf("skip PostgreSQL integration test: %v", sharedPostgres.err)
	}
	return &sharedPostgres
}

func pingWithRetry(ctx context.Context, db *sql.DB) error {
	deadline := time.Now().Add(10 * time.Second)
	var lastErr error
	for {
		if err := db.PingContext(ctx); err == nil {
			return nil
		} else {
			lastErr = err
		}
		if time.Now().After(deadline) {
			return lastErr
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(200 * time.Millisecond):
		}
	}
}

func runMigrations(ctx context.Context, db *sql.DB) error {
	files, err := migrationFiles()
	if err != nil {
		return err
	}
	for _, path := range files {
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read migration %s failed: %w", path, err)
		}
		if _, err := db.ExecContext(ctx, string(content)); err != nil {
			return fmt.Errorf("apply migration %s failed: %w", path, err)
		}
	}
	return nil
}

func migrationFiles() ([]string, error) {
	var files []string
	for _, name := range []string{"identity", "content"} {
		matches, err := filepath.Glob(filepath.Join(repoRootDir(), "sql", "migrations", "v3", name, "*.sql"))
		if err != nil {
			return nil, fmt.Errorf("glob %s migrations failed: %w", name, err)
		}
		files = append(files, matches...)
	}
	sort.Slice(files, func(i, j int) bool {
		return filepath.Base(files[i]) < filepath.Base(files[j])
	})
	return files, nil
}

func seedUsers(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
INSERT INTO identity.users (id, username, email, nickname, role, status)
VALUES (1, 'content_actor', 'content-actor@example.test', 'Content Actor', 'admin', 'active');
SELECT setval(pg_get_serial_sequence('identity.users', 'id'), 1, true);
`)
	return err
}

func dsnForDatabase(rawDSN, dbName string) (string, error) {
	parsed, err := url.Parse(rawDSN)
	if err != nil {
		return "", err
	}
	parsed.Path = path.Join("/", dbName)
	return parsed.String(), nil
}

func quoteIdent(value string) string {
	return `"` + strings.ReplaceAll(value, `"`, `""`) + `"`
}

func repoRootDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "..", ".."))
}
