// Package main 提供 Beehive-Blog 的 SQL 迁移 CLI。
//
// 模式说明：
//   - versioned（默认，全覆盖）：每个迁移文件在一个事务中原子执行，依赖 schema_migrations
//     的 checksum 做版本一致性校验；适合空库或严格与仓库迁移历史一致的环境。
//   - adaptive（适应）：将单个迁移文件按语句拆分后顺序执行；遇到「对象已存在」等
//     可预期 SQLSTATE 时跳过该语句并继续，用于已有手工表结构、重复执行或半旧库向前对齐。
package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	var (
		dsn           = flag.String("dsn", envOrDefault("DB_DSN", "postgres://Beehive-Blog:Beehive-Blog@127.0.0.1:5432/Beehive-Blog?sslmode=disable"), "PostgreSQL DSN")
		migrationsDir = flag.String("dir", "sql/migrations", "迁移 SQL 所在目录（相对当前工作目录或绝对路径）")
		mode          = flag.String("mode", "versioned", "迁移模式：versioned=全覆盖（整文件原子）；adaptive=适应（按语句，跳过已存在类错误）")
		verbose       = flag.Bool("v", false, "adaptive 模式下打印被跳过的语句错误")
	)
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	db, err := sql.Open("pgx", *dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err = db.PingContext(ctx); err != nil {
		panic(err)
	}

	if err = ensureSchemaMigrationsTable(ctx, db); err != nil {
		panic(err)
	}

	files, err := listMigrationFiles(*migrationsDir)
	if err != nil {
		panic(err)
	}

	m := strings.ToLower(strings.TrimSpace(*mode))
	if m != "versioned" && m != "adaptive" {
		panic(fmt.Errorf("unknown -mode %q (use versioned or adaptive)", *mode))
	}

	for _, f := range files {
		version := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
		path := filepath.Join(*migrationsDir, f.Name())
		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		checksum := sha256Hex(sqlBytes)
		applied, err := isApplied(ctx, db, version, checksum)
		if err != nil {
			panic(err)
		}
		if applied {
			fmt.Printf("skip %s\n", version)
			continue
		}

		switch m {
		case "versioned":
			if err := applyVersioned(ctx, db, version, checksum, string(sqlBytes)); err != nil {
				panic(err)
			}
		case "adaptive":
			if err := applyAdaptive(ctx, db, version, checksum, string(sqlBytes), *verbose); err != nil {
				panic(err)
			}
		default:
			panic("unreachable")
		}
		fmt.Printf("applied %s (%s)\n", version, m)
	}

	fmt.Println("migrations completed")
}

func applyVersioned(ctx context.Context, db *sql.DB, version, checksum, body string) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, body); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("apply %s failed: %w", version, err)
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO schema_migrations(version, checksum) VALUES ($1, $2)`, version, checksum); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("record %s failed: %w", version, err)
	}
	return tx.Commit()
}

func applyAdaptive(ctx context.Context, db *sql.DB, version, checksum, body string, verbose bool) error {
	stmts := splitMigrationStatements(body)
	if len(stmts) == 0 {
		return fmt.Errorf("%s: no executable statements after split", version)
	}

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	for i, stmt := range stmts {
		if _, err := tx.ExecContext(ctx, stmt); err != nil {
			code := pgSQLState(err)
			if _, ok := adaptiveSkipSQLSTATE[code]; ok {
				if verbose {
					fmt.Fprintf(os.Stderr, "adaptive skip %s stmt#%d sqlstate=%s: %v\n", version, i+1, code, err)
				}
				continue
			}
			_ = tx.Rollback()
			return fmt.Errorf("apply %s stmt#%d failed: %w", version, i+1, err)
		}
	}

	if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations(version, checksum) VALUES ($1, $2)`, version, checksum); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("record %s failed: %w", version, err)
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// adaptiveSkipSQLSTATE 为「对象已存在 / 重复定义」等可安全跳过、继续后续语句的状态码。
// 参考: https://www.postgresql.org/docs/current/errcodes-appendix.html
var adaptiveSkipSQLSTATE = map[string]struct{}{
	"42P07": {}, // duplicate_table
	"42701": {}, // duplicate_column
	"42710": {}, // duplicate_object
	"23505": {}, // unique_violation（唯一约束/索引已存在等）
}

func pgSQLState(err error) string {
	var e *pgconn.PgError
	if errors.As(err, &e) {
		return e.Code
	}
	return ""
}

var stmtSplitPattern = regexp.MustCompile(`;\s*\n`)

func splitMigrationStatements(body string) []string {
	body = strings.TrimSpace(strings.ReplaceAll(body, "\r\n", "\n"))
	if body == "" {
		return nil
	}
	raw := stmtSplitPattern.Split(body, -1)
	var out []string
	for _, chunk := range raw {
		s := stripSQLLineComments(strings.TrimSpace(chunk))
		if s == "" {
			continue
		}
		out = append(out, s)
	}
	return out
}

func stripSQLLineComments(s string) string {
	lines := strings.Split(s, "\n")
	var b strings.Builder
	for _, ln := range lines {
		t := strings.TrimSpace(ln)
		if t == "" || strings.HasPrefix(t, "--") {
			continue
		}
		b.WriteString(ln)
		b.WriteByte('\n')
	}
	return strings.TrimSpace(b.String())
}

func ensureSchemaMigrationsTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    checksum VARCHAR(64) NOT NULL,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)`)
	return err
}

func listMigrationFiles(dir string) ([]os.DirEntry, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	files := make([]os.DirEntry, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".sql") && len(name) >= 8 && name[0] >= '0' && name[0] <= '9' {
			files = append(files, e)
		}
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	return files, nil
}

func isApplied(ctx context.Context, db *sql.DB, version, checksum string) (bool, error) {
	var existing string
	err := db.QueryRowContext(ctx, `SELECT checksum FROM schema_migrations WHERE version = $1`, version).Scan(&existing)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	if existing != checksum {
		return false, fmt.Errorf("migration %s checksum mismatch (recorded vs file); 请检查是否手工改过迁移文件或需重建库", version)
	}
	return true, nil
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func envOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
