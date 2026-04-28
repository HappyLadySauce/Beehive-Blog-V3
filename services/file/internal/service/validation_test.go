package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/config"
)

func TestNormalizeScope(t *testing.T) {
	t.Parallel()

	scope, err := normalizeScope("content_cover")
	if err != nil {
		t.Fatalf("expected content_cover to pass, got %v", err)
	}
	if scope != ScopeContentCover {
		t.Fatalf("expected %s, got %s", ScopeContentCover, scope)
	}

	if _, err := normalizeScope("profile_banner"); !errors.Is(err, errs.E(errs.CodeFileInvalidScope)) {
		t.Fatalf("expected invalid scope error, got %v", err)
	}
}

func TestValidateUploadFile(t *testing.T) {
	t.Parallel()

	conf := config.ObjectStorageConf{
		MaxBytesByScope: map[string]int64{
			ScopeAvatar: 128,
		},
		AllowedContentTypesByScope: map[string][]string{
			ScopeAvatar: []string{"image/png"},
		},
	}

	contentType, maxBytes, err := validateUploadFile(conf, ScopeAvatar, "avatar.png", "image/png", 100)
	if err != nil {
		t.Fatalf("expected upload file validation to pass, got %v", err)
	}
	if contentType != "image/png" || maxBytes != 128 {
		t.Fatalf("unexpected normalized values: contentType=%s maxBytes=%d", contentType, maxBytes)
	}

	if _, _, err := validateUploadFile(conf, ScopeAvatar, "avatar.gif", "image/gif", 100); !errors.Is(err, errs.E(errs.CodeFileInvalidContentType)) {
		t.Fatalf("expected invalid content type error, got %v", err)
	}
	if _, _, err := validateUploadFile(conf, ScopeAvatar, "avatar.png", "image/png", 129); !errors.Is(err, errs.E(errs.CodeFileTooLarge)) {
		t.Fatalf("expected file too large error, got %v", err)
	}
}

func TestObjectKeyUsesScopePrefix(t *testing.T) {
	t.Parallel()

	key := objectKey(ScopeContentImage, 42, "cover.jpeg", "image/jpeg")
	if !strings.HasPrefix(key, "content/images/42/") {
		t.Fatalf("expected content image prefix, got %s", key)
	}
	if !strings.HasSuffix(key, ".jpg") {
		t.Fatalf("expected jpeg extension normalization, got %s", key)
	}
}
