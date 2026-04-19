// Package main 提供 Beehive-Blog 的 SQL 迁移 CLI。
//
// 模式说明：
//   - versioned（默认，全覆盖）：每个迁移文件在一个事务中原子执行，并记录 checksum。
//   - adaptive（适应）：按 SQL 语句执行；仅跳过“对象已存在”类错误。
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
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	defaultDSN          = "postgres://Beehive-Blog-V3:Beehive-Blog-V3@127.0.0.1:5432/Beehive-Blog-V3?sslmode=disable"
	modeVersioned       = "versioned"
	modeAdaptive        = "adaptive"
	migrateLockID int64 = 903241127
)

type migrationFile struct {
	Name     string
	Version  string
	Path     string
	Body     string
	Checksum string
}

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "migration error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		dsn           = flag.String("dsn", envOrDefault("DB_DSN", defaultDSN), "PostgreSQL DSN")
		migrationsDir = flag.String("dir", "sql/migrations", "迁移 SQL 所在目录（相对当前工作目录或绝对路径）")
		mode          = flag.String("mode", modeVersioned, "迁移模式：versioned=全覆盖（整文件原子）；adaptive=适应（按语句，仅跳过已存在类错误）")
		verbose       = flag.Bool("v", false, "adaptive 模式下打印被跳过的语句错误")
	)
	flag.Parse()

	m := strings.ToLower(strings.TrimSpace(*mode))
	if m != modeVersioned && m != modeAdaptive {
		return fmt.Errorf("unknown -mode %q (use versioned or adaptive)", *mode)
	}

	files, err := listMigrationFiles(*migrationsDir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		fmt.Println("no migration files found, skipped")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	db, err := sql.Open("pgx", *dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if err = db.PingContext(ctx); err != nil {
		return err
	}
	if err = ensureSchemaMigrationsTable(ctx, db); err != nil {
		return err
	}
	if err = lockMigrations(ctx, db); err != nil {
		return err
	}
	defer unlockMigrations(db)

	for _, mf := range files {
		applied, err := isApplied(ctx, db, mf.Version, mf.Checksum)
		if err != nil {
			return err
		}
		if applied {
			fmt.Printf("skip %s\n", mf.Version)
			continue
		}

		switch m {
		case modeVersioned:
			err = applyVersioned(ctx, db, mf.Version, mf.Checksum, mf.Body)
		case modeAdaptive:
			err = applyAdaptive(ctx, db, mf.Version, mf.Checksum, mf.Body, *verbose)
		}
		if err != nil {
			return err
		}
		fmt.Printf("applied %s (%s)\n", mf.Version, m)
	}

	fmt.Println("migrations completed")
	return nil
}

func applyVersioned(ctx context.Context, db *sql.DB, version, checksum, body string) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx, body); err != nil {
		return fmt.Errorf("apply %s failed: %w", version, err)
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO schema_migrations(version, checksum) VALUES ($1, $2)`, version, checksum); err != nil {
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
	defer tx.Rollback()

	for i, stmt := range stmts {
		if _, err := tx.ExecContext(ctx, stmt); err != nil {
			code := pgSQLState(err)
			if _, ok := adaptiveSkipSQLSTATE[code]; ok {
				if verbose {
					fmt.Fprintf(os.Stderr, "adaptive skip %s stmt#%d sqlstate=%s: %v\n", version, i+1, code, err)
				}
				continue
			}
			return fmt.Errorf("apply %s stmt#%d failed: %w", version, i+1, err)
		}
	}

	if _, err = tx.ExecContext(ctx, `INSERT INTO schema_migrations(version, checksum) VALUES ($1, $2)`, version, checksum); err != nil {
		return fmt.Errorf("record %s failed: %w", version, err)
	}
	return tx.Commit()
}

// adaptiveSkipSQLSTATE 仅包含“对象重复定义”场景。
// 注意：不要加入 23505（unique_violation），它常常代表真实的数据冲突。
var adaptiveSkipSQLSTATE = map[string]struct{}{
	"42P07": {}, // duplicate_table
	"42701": {}, // duplicate_column
	"42710": {}, // duplicate_object
}

func pgSQLState(err error) string {
	var e *pgconn.PgError
	if errors.As(err, &e) {
		return e.Code
	}
	return ""
}

func splitMigrationStatements(body string) []string {
	body = strings.ReplaceAll(body, "\r\n", "\n")
	body = strings.TrimSpace(body)
	if body == "" {
		return nil
	}

	var (
		out            []string
		current        strings.Builder
		inSingleQuote  bool
		inDoubleQuote  bool
		inLineComment  bool
		inBlockComment bool
		dollarTag      string
	)

	runes := []rune(body)
	for i := 0; i < len(runes); i++ {
		ch := runes[i]
		var next rune
		if i+1 < len(runes) {
			next = runes[i+1]
		}

		if inLineComment {
			current.WriteRune(ch)
			if ch == '\n' {
				inLineComment = false
			}
			continue
		}
		if inBlockComment {
			current.WriteRune(ch)
			if ch == '*' && next == '/' {
				current.WriteRune(next)
				i++
				inBlockComment = false
			}
			continue
		}
		if inSingleQuote {
			current.WriteRune(ch)
			if ch == '\'' {
				if next == '\'' {
					current.WriteRune(next)
					i++
				} else {
					inSingleQuote = false
				}
			}
			continue
		}
		if inDoubleQuote {
			current.WriteRune(ch)
			if ch == '"' {
				inDoubleQuote = false
			}
			continue
		}
		if dollarTag != "" {
			current.WriteRune(ch)
			if ch == '$' && tailMatchTag(runes, i, dollarTag) {
				dollarTag = ""
			}
			continue
		}

		if ch == '-' && next == '-' {
			current.WriteRune(ch)
			current.WriteRune(next)
			i++
			inLineComment = true
			continue
		}
		if ch == '/' && next == '*' {
			current.WriteRune(ch)
			current.WriteRune(next)
			i++
			inBlockComment = true
			continue
		}
		if ch == '\'' {
			current.WriteRune(ch)
			inSingleQuote = true
			continue
		}
		if ch == '"' {
			current.WriteRune(ch)
			inDoubleQuote = true
			continue
		}
		if ch == '$' {
			if tag, ok := parseDollarTag(runes, i); ok {
				current.WriteString(tag)
				i += len([]rune(tag)) - 1
				dollarTag = tag
				continue
			}
		}
		if ch == ';' {
			stmt := strings.TrimSpace(current.String())
			if stmt != "" {
				out = append(out, stmt)
			}
			current.Reset()
			continue
		}

		current.WriteRune(ch)
	}

	if tail := strings.TrimSpace(current.String()); tail != "" {
		out = append(out, tail)
	}
	return out
}

func parseDollarTag(runes []rune, start int) (string, bool) {
	for i := start + 1; i < len(runes); i++ {
		if runes[i] == '$' {
			tag := string(runes[start : i+1])
			if isValidDollarTag(tag) {
				return tag, true
			}
			return "", false
		}
		if !isTagRune(runes[i]) {
			return "", false
		}
	}
	return "", false
}

func isValidDollarTag(tag string) bool {
	if len(tag) < 2 || tag[0] != '$' || tag[len(tag)-1] != '$' {
		return false
	}
	for _, r := range tag[1 : len(tag)-1] {
		if !isTagRune(r) {
			return false
		}
	}
	return true
}

func isTagRune(r rune) bool {
	return r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

func tailMatchTag(runes []rune, pos int, tag string) bool {
	tagRunes := []rune(tag)
	if len(tagRunes) == 0 {
		return false
	}
	start := pos - len(tagRunes) + 1
	if start < 0 {
		return false
	}
	for i := 0; i < len(tagRunes); i++ {
		if runes[start+i] != tagRunes[i] {
			return false
		}
	}
	return true
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

func lockMigrations(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `SELECT pg_advisory_lock($1)`, migrateLockID)
	return err
}

func unlockMigrations(db *sql.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, _ = db.ExecContext(ctx, `SELECT pg_advisory_unlock($1)`, migrateLockID)
}

func listMigrationFiles(dir string) ([]migrationFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(strings.ToLower(name), ".sql") && len(name) >= 8 && name[0] >= '0' && name[0] <= '9' {
			names = append(names, name)
		}
	}
	sort.Strings(names)

	files := make([]migrationFile, 0, len(names))
	for _, name := range names {
		path := filepath.Join(dir, name)
		bodyBytes, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		files = append(files, migrationFile{
			Name:     name,
			Version:  strings.TrimSuffix(name, filepath.Ext(name)),
			Path:     path,
			Body:     string(bodyBytes),
			Checksum: sha256Hex(bodyBytes),
		})
	}
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
