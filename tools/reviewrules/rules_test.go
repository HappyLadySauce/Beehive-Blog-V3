package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunPassesWithAllowedUsage(t *testing.T) {
	t.Helper()

	root := t.TempDir()
	writeTestFile(t, root, filepath.Join("services", "demo", "main.go"), `package demo

import (
	"errors"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
)

func test(err error) bool {
	logs.Ctx(nil).Info("demo")
	return errors.Is(err, errs.E(errs.CodeIdentityInvalidCredentials))
}
`)

	runWithTempRepo(t, root, func() {
		if err := Run(); err != nil {
			t.Fatalf("expected review rules to pass, got %v", err)
		}
	})
}

func TestRunRejectsDirectLogxImport(t *testing.T) {
	t.Helper()

	root := t.TempDir()
	writeTestFile(t, root, filepath.Join("services", "demo", "main.go"), `package demo

import "github.com/zeromicro/go-zero/core/logx"

func test() {
	logx.Info("bad")
}
`)

	runWithTempRepo(t, root, func() {
		err := Run()
		if err == nil {
			t.Fatal("expected review rules to fail for direct logx usage")
		}
		if !strings.Contains(err.Error(), "no_direct_logx_import") {
			t.Fatalf("expected logx import violation, got %v", err)
		}
	})
}

func TestRunRejectsErrsIsCode(t *testing.T) {
	t.Helper()

	root := t.TempDir()
	writeTestFile(t, root, filepath.Join("services", "demo", "main.go"), `package demo

import "github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"

func test(err error) bool {
	return errs.IsCode(err, errs.CodeIdentityInvalidCredentials)
}
`)

	runWithTempRepo(t, root, func() {
		err := Run()
		if err == nil {
			t.Fatal("expected review rules to fail for errs.IsCode")
		}
		if !strings.Contains(err.Error(), "prefer_errors_is_over_errs_iscode") {
			t.Fatalf("expected errs.IsCode violation, got %v", err)
		}
	})
}

func TestRunRejectsStringBasedErrorMatching(t *testing.T) {
	t.Helper()

	root := t.TempDir()
	writeTestFile(t, root, filepath.Join("services", "demo", "main.go"), `package demo

import "strings"

func test(err error) bool {
	if strings.Contains(err.Error(), "bad") {
		return true
	}
	return err.Error() == "bad"
}
`)

	runWithTempRepo(t, root, func() {
		err := Run()
		if err == nil {
			t.Fatal("expected review rules to fail for string-based error matching")
		}
		if !strings.Contains(err.Error(), "no_string_match_on_error_contains") {
			t.Fatalf("expected contains violation, got %v", err)
		}
		if !strings.Contains(err.Error(), "no_string_match_on_error_equals") {
			t.Fatalf("expected equals violation, got %v", err)
		}
	})
}

func TestRunAllowsPkgLogsToUseLogx(t *testing.T) {
	t.Helper()

	root := t.TempDir()
	writeTestFile(t, root, filepath.Join("pkg", "logs", "logs.go"), `package logs

import "github.com/zeromicro/go-zero/core/logx"

func test() {
	logx.Info("allowed")
}
`)

	runWithTempRepo(t, root, func() {
		if err := Run(); err != nil {
			t.Fatalf("expected pkg/logs to be excluded, got %v", err)
		}
	})
}

func TestRunAllowsGeneratedContentPB(t *testing.T) {
	t.Helper()

	root := t.TempDir()
	writeTestFile(t, root, filepath.Join("services", "content", "pb", "content.pb.go"), `package pb

import "github.com/zeromicro/go-zero/core/logx"

func test() {
	logx.Info("generated code is excluded")
}
`)

	runWithTempRepo(t, root, func() {
		if err := Run(); err != nil {
			t.Fatalf("expected services/content/pb to be excluded, got %v", err)
		}
	})
}

func runWithTempRepo(t *testing.T, root string, fn func()) {
	t.Helper()

	previous, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	if err := os.Chdir(root); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(previous)
	})

	fn()
}

func writeTestFile(t *testing.T, root, relativePath, content string) {
	t.Helper()

	path := filepath.Join(root, relativePath)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("failed to create directory for %s: %v", relativePath, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write %s: %v", relativePath, err)
	}
}
