package service

import (
	"mime"
	"path"
	"strconv"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/config"
)

func parseActorUserID(actorUserID string) (int64, error) {
	value, err := strconv.ParseInt(strings.TrimSpace(actorUserID), 10, 64)
	if err != nil || value <= 0 {
		return 0, errs.New(errs.CodeFileInvalidArgument, "actor_user_id is invalid")
	}
	return value, nil
}

func normalizeNamespace(conf config.StorageConf, namespace string) (string, error) {
	namespace = strings.ToLower(strings.TrimSpace(namespace))
	if namespace == "" || len(namespace) > 64 {
		return "", errs.New(errs.CodeFileInvalidScope, "namespace is invalid")
	}
	if _, ok := conf.NamespaceRule(namespace); !ok {
		return "", errs.New(errs.CodeFileInvalidScope, "namespace is not configured")
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

func normalizeOptionalNamespace(conf config.StorageConf, namespace string) (string, error) {
	if strings.TrimSpace(namespace) == "" {
		return "", nil
	}
	return normalizeNamespace(conf, namespace)
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

func maxBytesForNamespace(conf config.StorageConf, namespace string) int64 {
	if rule, ok := conf.NamespaceRule(namespace); ok && rule.MaxBytes > 0 {
		return rule.MaxBytes
	}
	if conf.MaxBytesByScope != nil && conf.MaxBytesByScope[namespace] > 0 {
		return conf.MaxBytesByScope[namespace]
	}
	return 0
}

func allowedContentTypesForNamespace(conf config.StorageConf, namespace string) map[string]struct{} {
	var values []string
	if rule, ok := conf.NamespaceRule(namespace); ok && len(rule.AllowedContentTypes) > 0 {
		values = rule.AllowedContentTypes
	} else if conf.AllowedContentTypesByScope != nil && len(conf.AllowedContentTypesByScope[namespace]) > 0 {
		values = conf.AllowedContentTypesByScope[namespace]
	}
	if len(values) == 0 {
		return nil
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

func validateUploadFile(conf config.StorageConf, namespace string, fileName string, contentType string, byteSize int64) (string, int64, error) {
	fileName = strings.TrimSpace(fileName)
	if fileName == "" || len(fileName) > 255 {
		return "", 0, invalidArgument("file_name is invalid")
	}
	normalizedContentType := normalizeContentType(contentType)
	if normalizedContentType == "" || len(normalizedContentType) > 128 {
		return "", 0, invalidArgument("content_type is invalid")
	}
	if allowed := allowedContentTypesForNamespace(conf, namespace); allowed != nil {
		if _, ok := allowed[normalizedContentType]; !ok {
			return "", 0, errs.New(errs.CodeFileInvalidContentType, "content_type is not allowed")
		}
	}
	maxBytes := maxBytesForNamespace(conf, namespace)
	if maxBytes <= 0 {
		maxBytes = 5 * 1024 * 1024
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
