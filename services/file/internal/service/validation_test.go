package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/config"
)

func TestNormalizeCategoryKey(t *testing.T) {
	t.Parallel()

	ns, err := normalizeCategoryKey("content_cover")
	if err != nil {
		t.Fatalf("expected content_cover to pass, got %v", err)
	}
	if ns != "content_cover" {
		t.Fatalf("expected content_cover, got %s", ns)
	}

	// Empty category key
	if _, err := normalizeCategoryKey(""); !errors.Is(err, errs.E(errs.CodeFileInvalidScope)) {
		t.Fatalf("expected empty category key error, got %v", err)
	}

	// Too long category key
	if _, err := normalizeCategoryKey(strings.Repeat("a", 65)); !errors.Is(err, errs.E(errs.CodeFileInvalidScope)) {
		t.Fatalf("expected too-long category key error, got %v", err)
	}

	// Invalid characters
	if _, err := normalizeCategoryKey("bad category!"); !errors.Is(err, errs.E(errs.CodeFileInvalidScope)) {
		t.Fatalf("expected invalid char error, got %v", err)
	}

	// Valid: any alphanumeric + separator string passes
	ns2, err := normalizeCategoryKey("banner_homepage:en")
	if err != nil {
		t.Fatalf("expected banner_homepage:en to pass, got %v", err)
	}
	if ns2 != "banner_homepage:en" {
		t.Fatalf("expected banner_homepage:en, got %s", ns2)
	}
}

func TestValidateUploadFile(t *testing.T) {
	t.Parallel()

	allowedExtensions := []string{".png", ".jpg"}
	maxUploadBytes := int64(128)

	contentType, extension, maxBytes, err := validateUploadFile(allowedExtensions, maxUploadBytes, "avatar.png", "image/png", 100)
	if err != nil {
		t.Fatalf("expected upload file validation to pass, got %v", err)
	}
	if contentType != "image/png" || extension != ".png" || maxBytes != 128 {
		t.Fatalf("unexpected normalized values: contentType=%s extension=%s maxBytes=%d", contentType, extension, maxBytes)
	}

	if _, _, _, err := validateUploadFile(allowedExtensions, maxUploadBytes, "avatar.gif", "image/gif", 100); !errors.Is(err, errs.E(errs.CodeFileInvalidExtension)) {
		t.Fatalf("expected invalid extension error, got %v", err)
	}
	if _, _, _, err := validateUploadFile(allowedExtensions, maxUploadBytes, "avatar.png", "text/html", 100); !errors.Is(err, errs.E(errs.CodeFileInvalidArgument)) {
		t.Fatalf("expected mismatched content type error, got %v", err)
	}
	if _, _, _, err := validateUploadFile(allowedExtensions, maxUploadBytes, "avatar.png", "image/png", 129); !errors.Is(err, errs.E(errs.CodeFileTooLarge)) {
		t.Fatalf("expected file too large error, got %v", err)
	}

	// Missing content type should fall back to extension inference.
	contentType, extension, maxBytes, err = validateUploadFile([]string{".pdf"}, 0, "file.pdf", "", 100)
	if err != nil {
		t.Fatalf("expected PDF to pass with extension inference, got %v", err)
	}
	if contentType != "application/pdf" || extension != ".pdf" {
		t.Fatalf("expected inferred pdf values, got contentType=%s extension=%s", contentType, extension)
	}
	if maxBytes != hardcodedMaxBytes {
		t.Fatalf("expected hardcoded max bytes %d, got %d", hardcodedMaxBytes, maxBytes)
	}

	// Unknown extensions only accept octet-stream metadata to prevent arbitrary MIME publication.
	contentType, extension, _, err = validateUploadFile([]string{".bin"}, maxUploadBytes, "archive.bin", "application/octet-stream", 100)
	if err != nil {
		t.Fatalf("expected octet-stream to pass for unknown extension, got %v", err)
	}
	if contentType != "application/octet-stream" || extension != ".bin" {
		t.Fatalf("expected octet-stream values, got contentType=%s extension=%s", contentType, extension)
	}
	if _, _, _, err := validateUploadFile([]string{".bin"}, maxUploadBytes, "archive.bin", "text/plain", 100); !errors.Is(err, errs.E(errs.CodeFileInvalidArgument)) {
		t.Fatalf("expected invalid argument for arbitrary unknown-extension content type, got %v", err)
	}
}

func TestObjectKeyUsesCategoryKeyAsPrefix(t *testing.T) {
	t.Parallel()

	key := objectKey("content_image", 42, ".jpg")
	if !strings.HasPrefix(key, "content_image/42/") {
		t.Fatalf("expected content_image prefix, got %s", key)
	}
	if !strings.HasSuffix(key, ".jpg") {
		t.Fatalf("expected jpeg extension normalization, got %s", key)
	}

	// Any category key becomes the path prefix.
	key2 := objectKey("banner", 42, ".png")
	if !strings.HasPrefix(key2, "banner/42/") {
		t.Fatalf("expected banner prefix, got %s", key2)
	}

	// Empty category key falls back to default.
	key3 := objectKey("", 42, ".pdf")
	if !strings.HasPrefix(key3, "default/42/") {
		t.Fatalf("expected default prefix for empty category key, got %s", key3)
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
