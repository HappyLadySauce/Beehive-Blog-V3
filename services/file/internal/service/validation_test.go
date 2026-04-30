package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/config"
)

func TestNormalizeNamespace(t *testing.T) {
	t.Parallel()

	ns, err := normalizeNamespace("content_cover")
	if err != nil {
		t.Fatalf("expected content_cover to pass, got %v", err)
	}
	if ns != "content_cover" {
		t.Fatalf("expected content_cover, got %s", ns)
	}

	// Empty namespace
	if _, err := normalizeNamespace(""); !errors.Is(err, errs.E(errs.CodeFileInvalidScope)) {
		t.Fatalf("expected empty namespace error, got %v", err)
	}

	// Too long namespace
	if _, err := normalizeNamespace(strings.Repeat("a", 65)); !errors.Is(err, errs.E(errs.CodeFileInvalidScope)) {
		t.Fatalf("expected too-long namespace error, got %v", err)
	}

	// Invalid characters
	if _, err := normalizeNamespace("bad namespace!"); !errors.Is(err, errs.E(errs.CodeFileInvalidScope)) {
		t.Fatalf("expected invalid char error, got %v", err)
	}

	// Valid: any alphanumeric + separator string passes
	ns2, err := normalizeNamespace("banner_homepage:en")
	if err != nil {
		t.Fatalf("expected banner_homepage:en to pass, got %v", err)
	}
	if ns2 != "banner_homepage:en" {
		t.Fatalf("expected banner_homepage:en, got %s", ns2)
	}
}

func TestValidateUploadFile(t *testing.T) {
	t.Parallel()

	conf := config.StorageConf{
		AllowedContentTypes: []string{"image/png", "image/jpeg"},
		MaxUploadBytes:      128,
	}

	contentType, maxBytes, err := validateUploadFile(conf, "avatar.png", "image/png", 100)
	if err != nil {
		t.Fatalf("expected upload file validation to pass, got %v", err)
	}
	if contentType != "image/png" || maxBytes != 128 {
		t.Fatalf("unexpected normalized values: contentType=%s maxBytes=%d", contentType, maxBytes)
	}

	if _, _, err := validateUploadFile(conf, "avatar.gif", "image/gif", 100); !errors.Is(err, errs.E(errs.CodeFileInvalidContentType)) {
		t.Fatalf("expected invalid content type error, got %v", err)
	}
	if _, _, err := validateUploadFile(conf, "avatar.png", "image/png", 129); !errors.Is(err, errs.E(errs.CodeFileTooLarge)) {
		t.Fatalf("expected file too large error, got %v", err)
	}

	// Empty config uses hardcoded fallback
	emptyConf := config.StorageConf{}
	_, fallbackMaxBytes, err := validateUploadFile(emptyConf, "file.pdf", "application/pdf", 100)
	if err != nil {
		t.Fatalf("expected PDF to pass with hardcoded fallback, got %v", err)
	}
	if fallbackMaxBytes != hardcodedMaxBytes {
		t.Fatalf("expected hardcoded max bytes %d, got %d", hardcodedMaxBytes, fallbackMaxBytes)
	}
}

func TestObjectKeyUsesNamespaceAsPrefix(t *testing.T) {
	t.Parallel()

	key := objectKey("content_image", 42, "cover.jpeg", "image/jpeg")
	if !strings.HasPrefix(key, "content_image/42/") {
		t.Fatalf("expected content_image prefix, got %s", key)
	}
	if !strings.HasSuffix(key, ".jpg") {
		t.Fatalf("expected jpeg extension normalization, got %s", key)
	}

	// Any namespace string becomes the path prefix
	key2 := objectKey("banner", 42, "hero.png", "image/png")
	if !strings.HasPrefix(key2, "banner/42/") {
		t.Fatalf("expected banner prefix, got %s", key2)
	}

	// Empty namespace defaults to "files"
	key3 := objectKey("", 42, "file.txt", "application/pdf")
	if !strings.HasPrefix(key3, "files/42/") {
		t.Fatalf("expected files prefix for empty namespace, got %s", key3)
	}
}

func TestPublicURLForVisibility(t *testing.T) {
	t.Parallel()

	publicURL := publicURLForVisibility(config.StorageConf{
		Driver:        "s3",
		PublicBaseURL: "https://cdn.example.com/files/",
	}, VisibilityPublic, "asset_1", "/avatars/42/avatar.png")
	if publicURL != "https://cdn.example.com/files/avatars/42/avatar.png" {
		t.Fatalf("unexpected public URL: %s", publicURL)
	}

	privateURL := publicURLForVisibility(config.StorageConf{
		Driver:        "s3",
		PublicBaseURL: "https://cdn.example.com/files/",
	}, VisibilityPrivate, "asset_1", "avatars/42/avatar.png")
	if privateURL != "" {
		t.Fatalf("expected private asset public URL to be empty, got %s", privateURL)
	}

	localURL := publicURLForVisibility(config.StorageConf{
		Driver:        "local",
		PublicBaseURL: "",
	}, VisibilityPublic, "asset_1", "avatars/42/avatar.png")
	if localURL != "" {
		t.Fatalf("expected local public URL to be empty when PublicBaseURL is empty, got %s", localURL)
	}
}
