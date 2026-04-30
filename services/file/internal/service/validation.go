package service

import (
	"mime"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/config"
)

var hardcodedContentTypes = []string{
	"image/png", "image/jpeg", "image/webp", "image/avif", "application/pdf",
}

const hardcodedMaxBytes int64 = 20 * 1024 * 1024

var validNamespaceRE = regexp.MustCompile(`^[a-z0-9_:.-]+$`)

func parseActorUserID(actorUserID string) (int64, error) {
	value, err := strconv.ParseInt(strings.TrimSpace(actorUserID), 10, 64)
	if err != nil || value <= 0 {
		return 0, errs.New(errs.CodeFileInvalidArgument, "actor_user_id is invalid")
	}
	return value, nil
}

func normalizeNamespace(namespace string) (string, error) {
	namespace = strings.ToLower(strings.TrimSpace(namespace))
	if namespace == "" || len(namespace) > 64 {
		return "", errs.New(errs.CodeFileInvalidScope, "namespace is invalid")
	}
	if !validNamespaceRE.MatchString(namespace) {
		return "", errs.New(errs.CodeFileInvalidScope, "namespace contains invalid characters")
	}
	return namespace, nil
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

func normalizeOptionalNamespace(namespace string) (string, error) {
	if strings.TrimSpace(namespace) == "" {
		return "", nil
	}
	return normalizeNamespace(namespace)
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

func validateUploadFile(conf config.StorageConf, fileName string, contentType string, byteSize int64) (string, int64, error) {
	fileName = strings.TrimSpace(fileName)
	if fileName == "" || len(fileName) > 255 {
		return "", 0, invalidArgument("file_name is invalid")
	}
	normalizedContentType := normalizeContentType(contentType)
	if normalizedContentType == "" || len(normalizedContentType) > 128 {
		return "", 0, invalidArgument("content_type is invalid")
	}

	types := conf.AllowedContentTypes
	if len(types) == 0 {
		types = hardcodedContentTypes
	}
	allowed := make(map[string]struct{}, len(types))
	for _, ct := range types {
		if n := normalizeContentType(ct); n != "" {
			allowed[n] = struct{}{}
		}
	}
	if _, ok := allowed[normalizedContentType]; !ok {
		return "", 0, errs.New(errs.CodeFileInvalidContentType, "content_type is not allowed")
	}

	maxBytes := conf.MaxUploadBytes
	if maxBytes <= 0 {
		maxBytes = hardcodedMaxBytes
	}
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
