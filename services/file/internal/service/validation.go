package service

import (
	"mime"
	"path"
	"strconv"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/config"
)

var defaultMaxBytesByScope = map[string]int64{
	ScopeAvatar:       2 * 1024 * 1024,
	ScopeContentCover: 5 * 1024 * 1024,
	ScopeContentImage: 5 * 1024 * 1024,
	ScopeAttachment:   20 * 1024 * 1024,
}

var defaultContentTypesByScope = map[string][]string{
	ScopeAvatar:       {"image/png", "image/jpeg", "image/webp", "image/avif"},
	ScopeContentCover: {"image/png", "image/jpeg", "image/webp", "image/avif"},
	ScopeContentImage: {"image/png", "image/jpeg", "image/webp", "image/avif"},
	ScopeAttachment:   {"image/png", "image/jpeg", "image/webp", "image/avif", "application/pdf"},
}

func parseActorUserID(actorUserID string) (int64, error) {
	value, err := strconv.ParseInt(strings.TrimSpace(actorUserID), 10, 64)
	if err != nil || value <= 0 {
		return 0, errs.New(errs.CodeFileInvalidArgument, "actor_user_id is invalid")
	}
	return value, nil
}

func normalizeScope(scope string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(scope)) {
	case ScopeAvatar:
		return ScopeAvatar, nil
	case ScopeContentCover:
		return ScopeContentCover, nil
	case ScopeContentImage:
		return ScopeContentImage, nil
	case ScopeAttachment:
		return ScopeAttachment, nil
	default:
		return "", errs.New(errs.CodeFileInvalidScope, "file scope is invalid")
	}
}

func normalizeVisibility(visibility string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(visibility)) {
	case "", VisibilityPublic:
		return VisibilityPublic, nil
	case VisibilityPrivate:
		return VisibilityPrivate, nil
	default:
		return "", errs.New(errs.CodeFileInvalidArgument, "asset visibility is invalid")
	}
}

func normalizeOptionalScope(scope string) (string, error) {
	if strings.TrimSpace(scope) == "" {
		return "", nil
	}
	return normalizeScope(scope)
}

func normalizeOptionalStatus(status string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "":
		return "", nil
	case StatusPending:
		return StatusPending, nil
	case StatusUploaded:
		return StatusUploaded, nil
	case StatusDeleted:
		return StatusDeleted, nil
	default:
		return "", errs.New(errs.CodeFileInvalidArgument, "asset status is invalid")
	}
}

func normalizeOptionalVisibility(visibility string) (string, error) {
	if strings.TrimSpace(visibility) == "" {
		return "", nil
	}
	return normalizeVisibility(visibility)
}

func normalizeContentType(contentType string) string {
	return strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
}

func maxBytesForScope(conf config.StorageConf, scope string) int64 {
	if conf.MaxBytesByScope != nil && conf.MaxBytesByScope[scope] > 0 {
		return conf.MaxBytesByScope[scope]
	}
	return defaultMaxBytesByScope[scope]
}

func allowedContentTypesForScope(conf config.StorageConf, scope string) map[string]struct{} {
	values := defaultContentTypesByScope[scope]
	if conf.AllowedContentTypesByScope != nil && len(conf.AllowedContentTypesByScope[scope]) > 0 {
		values = conf.AllowedContentTypesByScope[scope]
	}
	allowlist := make(map[string]struct{}, len(values))
	for _, value := range values {
		normalized := normalizeContentType(value)
		if normalized != "" {
			allowlist[normalized] = struct{}{}
		}
	}
	return allowlist
}

func validateUploadFile(conf config.StorageConf, scope string, fileName string, contentType string, byteSize int64) (string, int64, error) {
	fileName = strings.TrimSpace(fileName)
	if fileName == "" || len(fileName) > 255 {
		return "", 0, invalidArgument("file_name is invalid")
	}
	normalizedContentType := normalizeContentType(contentType)
	if normalizedContentType == "" || len(normalizedContentType) > 128 {
		return "", 0, invalidArgument("content_type is invalid")
	}
	if _, ok := allowedContentTypesForScope(conf, scope)[normalizedContentType]; !ok {
		return "", 0, errs.New(errs.CodeFileInvalidContentType, "content_type is not allowed")
	}
	maxBytes := maxBytesForScope(conf, scope)
	if byteSize <= 0 || byteSize > maxBytes {
		return "", 0, errs.New(errs.CodeFileTooLarge, "file byte_size is invalid")
	}
	return normalizedContentType, maxBytes, nil
}

func extensionFor(fileName string, contentType string) string {
	ext := strings.ToLower(path.Ext(strings.TrimSpace(fileName)))
	if ext == "" {
		extensions, _ := mime.ExtensionsByType(contentType)
		if len(extensions) > 0 {
			ext = extensions[0]
		}
	}
	if ext == ".jpeg" {
		return ".jpg"
	}
	return ext
}
