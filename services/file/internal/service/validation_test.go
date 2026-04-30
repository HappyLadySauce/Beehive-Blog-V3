package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/config"
)

func mkNamespacesConf() config.StorageConf {
	return config.StorageConf{
		Namespaces: map[string]config.NamespaceRule{
			"content_cover": {MaxBytes: 5 * 1024 * 1024, AllowedContentTypes: []string{"image/png", "image/jpeg", "image/webp", "image/avif"}, StoragePrefix: "content/covers"},
			"content_image": {MaxBytes: 5 * 1024 * 1024, AllowedContentTypes: []string{"image/png", "image/jpeg", "image/webp", "image/avif"}, StoragePrefix: "content/images"},
			"avatar":        {MaxBytes: 128, AllowedContentTypes: []string{"image/png"}, StoragePrefix: "avatars"},
			"*":             {MaxBytes: 5 * 1024 * 1024, AllowedContentTypes: []string{"image/png", "image/jpeg"}, StoragePrefix: "misc"},
		},
	}
}

func TestNormalizeNamespace(t *testing.T) {
	t.Parallel()

	conf := mkNamespacesConf()

	ns, err := normalizeNamespace(conf, "content_cover")
	if err != nil {
		t.Fatalf("expected content_cover to pass, got %v", err)
	}
	if ns != "content_cover" {
		t.Fatalf("expected content_cover, got %s", ns)
	}

	// Unknown namespace without wildcard fallback
	confNoWildcard := config.StorageConf{
		Namespaces: map[string]config.NamespaceRule{
			"avatar": {},
		},
	}
	if _, err := normalizeNamespace(confNoWildcard, "profile_banner"); !errors.Is(err, errs.E(errs.CodeFileInvalidScope)) {
		t.Fatalf("expected invalid namespace error, got %v", err)
	}

	// Unknown namespace with wildcard fallback should pass
	ns2, err := normalizeNamespace(conf, "banner")
	if err != nil {
		t.Fatalf("expected banner to pass via wildcard, got %v", err)
	}
	if ns2 != "banner" {
		t.Fatalf("expected banner, got %s", ns2)
	}
}

func TestValidateUploadFileNamespace(t *testing.T) {
	t.Parallel()

	conf := mkNamespacesConf()

	contentType, maxBytes, err := validateUploadFile(conf, "avatar", "avatar.png", "image/png", 100)
	if err != nil {
		t.Fatalf("expected upload file validation to pass, got %v", err)
	}
	if contentType != "image/png" || maxBytes != 128 {
		t.Fatalf("unexpected normalized values: contentType=%s maxBytes=%d", contentType, maxBytes)
	}

	if _, _, err := validateUploadFile(conf, "avatar", "avatar.gif", "image/gif", 100); !errors.Is(err, errs.E(errs.CodeFileInvalidContentType)) {
		t.Fatalf("expected invalid content type error, got %v", err)
	}
	if _, _, err := validateUploadFile(conf, "avatar", "avatar.png", "image/png", 129); !errors.Is(err, errs.E(errs.CodeFileTooLarge)) {
		t.Fatalf("expected file too large error, got %v", err)
	}
}

func TestObjectKeyUsesNamespacePrefix(t *testing.T) {
	t.Parallel()

	conf := mkNamespacesConf()
	key := objectKey(conf, "content_image", 42, "cover.jpeg", "image/jpeg")
	if !strings.HasPrefix(key, "content/images/42/") {
		t.Fatalf("expected content image prefix, got %s", key)
	}
	if !strings.HasSuffix(key, ".jpg") {
		t.Fatalf("expected jpeg extension normalization, got %s", key)
	}

	// Wildcard namespace uses "misc" prefix
	key2 := objectKey(conf, "banner", 42, "hero.png", "image/png")
	if !strings.HasPrefix(key2, "misc/42/") {
		t.Fatalf("expected misc prefix for wildcard namespace, got %s", key2)
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
