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

	_ "github.com/lib/pq"
)

const (
	modeVersioned = "versioned"
	modeAdaptive  = "adaptive"

	lockID int64 = 803514217
)

type migrationFile struct {
	Name     string
	FullPath string
	Content  string
	Checksum string
}

func main() {
	dsn := flag.String("dsn", "", "PostgreSQL DSN, e.g. postgres://user:pass@host:5432/db?sslmode=disable")
	dir := flag.String("dir", "", "migration directory")
	mode := flag.String("mode", modeVersioned, "migration mode: versioned | adaptive")
	verbose := flag.Bool("v", false, "enable verbose logging")
	flag.Parse()

	if strings.TrimSpace(*dsn) == "" {
		fatalf("missing required -dsn")
	}
	if strings.TrimSpace(*dir) == "" {
		fatalf("missing required -dir")
	}
	if *mode != modeVersioned && *mode != modeAdaptive {
		fatalf("invalid -mode: %s", *mode)
	}

	files, err := discoverMigrations(*dir)
	if err != nil {
		fatalf("discover migrations failed: %v", err)
	}
	if len(files) == 0 {
		fmt.Println("no migration files found, skipped")
		return
	}

	db, err := sql.Open("postgres", *dsn)
	if err != nil {
		fatalf("open database failed: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		fatalf("connect database failed: %v", err)
	}

	if err := ensureMigrationTable(ctx, db); err != nil {
		fatalf("ensure schema_migrations failed: %v", err)
	}

	if err := lock(ctx, db); err != nil {
		fatalf("acquire migration lock failed: %v", err)
	}
	defer unlock(db)

	for _, f := range files {
		applied, appliedChecksum, err := getMigrationState(ctx, db, f.Name)
		if err != nil {
			fatalf("query migration state failed (%s): %v", f.Name, err)
		}

		if applied {
			if !strings.EqualFold(appliedChecksum, f.Checksum) {
				fatalf("checksum mismatch for %s: db=%s file=%s", f.Name, appliedChecksum, f.Checksum)
			}
			if *verbose {
				fmt.Printf("skip applied migration: %s\n", f.Name)
			}
			continue
		}

		start := time.Now()
		switch *mode {
		case modeVersioned:
			err = runVersioned(ctx, db, f)
		case modeAdaptive:
			err = runAdaptive(ctx, db, f)
		}
		if err != nil {
			fatalf("apply migration failed (%s): %v", f.Name, err)
		}

		if err := markMigrationApplied(ctx, db, f.Name, f.Checksum); err != nil {
			fatalf("mark migration applied failed (%s): %v", f.Name, err)
		}

		fmt.Printf("applied %s (%s)\n", f.Name, time.Since(start).Round(time.Millisecond))
	}

	fmt.Println("all migrations completed")
}

func discoverMigrations(dir string) ([]migrationFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	files := make([]migrationFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".sql") {
			continue
		}

		full := filepath.Join(dir, name)
		content, err := os.ReadFile(full)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", name, err)
		}

		hash := sha256.Sum256(content)
		files = append(files, migrationFile{
			Name:     name,
			FullPath: full,
			Content:  string(content),
			Checksum: hex.EncodeToString(hash[:]),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})
	return files, nil
}

func ensureMigrationTable(ctx context.Context, db *sql.DB) error {
	const ddl = `
CREATE TABLE IF NOT EXISTS schema_migrations (
    filename VARCHAR(255) PRIMARY KEY,
    checksum CHAR(64) NOT NULL,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);`
	_, err := db.ExecContext(ctx, ddl)
	return err
}

func lock(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `SELECT pg_advisory_lock($1)`, lockID)
	return err
}

func unlock(db *sql.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, _ = db.ExecContext(ctx, `SELECT pg_advisory_unlock($1)`, lockID)
}

func getMigrationState(ctx context.Context, db *sql.DB, filename string) (bool, string, error) {
	var checksum string
	err := db.QueryRowContext(ctx, `SELECT checksum FROM schema_migrations WHERE filename = $1`, filename).Scan(&checksum)
	if errors.Is(err, sql.ErrNoRows) {
		return false, "", nil
	}
	if err != nil {
		return false, "", err
	}
	return true, checksum, nil
}

func markMigrationApplied(ctx context.Context, db *sql.DB, filename, checksum string) error {
	_, err := db.ExecContext(ctx, `
INSERT INTO schema_migrations (filename, checksum, applied_at)
VALUES ($1, $2, NOW())
ON CONFLICT (filename) DO UPDATE
SET checksum = EXCLUDED.checksum,
    applied_at = EXCLUDED.applied_at;`, filename, checksum)
	return err
}

func runVersioned(ctx context.Context, db *sql.DB, mf migrationFile) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, mf.Content); err != nil {
		return err
	}

	return tx.Commit()
}

func runAdaptive(ctx context.Context, db *sql.DB, mf migrationFile) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmts := splitSQLStatements(mf.Content)
	for _, stmt := range stmts {
		if strings.TrimSpace(stmt) == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, stmt); err != nil {
			if isSkippableError(err) {
				continue
			}
			return err
		}
	}

	return tx.Commit()
}

func splitSQLStatements(sqlText string) []string {
	out := make([]string, 0, 16)
	var current strings.Builder

	inSingle := false
	inDouble := false
	inLineComment := false
	inBlockComment := false
	inDollar := false
	dollarTag := ""

	runes := []rune(sqlText)
	for i := 0; i < len(runes); i++ {
		ch := runes[i]
		next := rune(0)
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
		if inSingle {
			current.WriteRune(ch)
			if ch == '\'' {
				if next == '\'' {
					current.WriteRune(next)
					i++
				} else {
					inSingle = false
				}
			}
			continue
		}
		if inDouble {
			current.WriteRune(ch)
			if ch == '"' {
				inDouble = false
			}
			continue
		}
		if inDollar {
			current.WriteRune(ch)
			if ch == '$' {
				segment := tailMatchTag(runes, i, dollarTag)
				if segment {
					inDollar = false
				}
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
			inSingle = true
			continue
		}
		if ch == '"' {
			current.WriteRune(ch)
			inDouble = true
			continue
		}
		if ch == '$' {
			if tag, ok := parseDollarTag(runes, i); ok {
				current.WriteString(tag)
				i += len([]rune(tag)) - 1
				inDollar = true
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
			if matched, _ := regexp.MatchString(`^\$[A-Za-z0-9_]*\$$`, tag); matched {
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

func isTagRune(r rune) bool {
	return r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

func isSkippableError(err error) bool {
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "already exists") {
		return true
	}
	if strings.Contains(msg, "duplicate key value violates unique constraint") {
		return true
	}
	pgCodeHints := []string{"sqlstate 42p07", "sqlstate 42701", "sqlstate 42710", "sqlstate 23505"}
	for _, hint := range pgCodeHints {
		if strings.Contains(msg, hint) {
			return true
		}
	}
	return false
}

func fatalf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, "migration error: "+format+"\n", args...)
	os.Exit(1)
}
