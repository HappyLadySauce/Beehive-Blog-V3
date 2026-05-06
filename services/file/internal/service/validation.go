package service

import (
	"mime"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
)

const hardcodedMaxBytes int64 = 20 * 1024 * 1024

var validCategoryKeyRE = regexp.MustCompile(`^[a-z0-9_:.-]+$`)

func parseActorUserID(actorUserID string) (int64, error) {
	value, err := strconv.ParseInt(strings.TrimSpace(actorUserID), 10, 64)
	if err != nil || value <= 0 {
		return 0, errs.New(errs.CodeFileInvalidArgument, "actor_user_id is invalid")
	}
	return value, nil
}

func normalizeCategoryKey(categoryKey string) (string, error) {
	categoryKey = strings.ToLower(strings.TrimSpace(categoryKey))
	if categoryKey == "" || len(categoryKey) > 64 {
		return "", errs.New(errs.CodeFileInvalidScope, "category_key is invalid")
	}
	if !validCategoryKeyRE.MatchString(categoryKey) {
		return "", errs.New(errs.CodeFileInvalidScope, "category_key contains invalid characters")
	}
	return categoryKey, nil
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

func normalizeOptionalCategoryKey(categoryKey string) (string, error) {
	if strings.TrimSpace(categoryKey) == "" {
		return "", nil
	}
	return normalizeCategoryKey(categoryKey)
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

func normalizeAllowedExtensions(extensions []string) []string {
	if len(extensions) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(extensions))
	normalized := make([]string, 0, len(extensions))
	for _, item := range extensions {
		value := normalizeExtension(item)
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		normalized = append(normalized, value)
	}
	slices.Sort(normalized)
	return normalized
}

func normalizeExtension(extension string) string {
	value := strings.ToLower(strings.TrimSpace(extension))
	if value == "" {
		return ""
	}
	if !strings.HasPrefix(value, ".") {
		value = "." + value
	}
	if value == ".jpeg" {
		return ".jpg"
	}
	return value
}

func validateUploadFile(allowedExtensions []string, maxUploadBytes int64, fileName string, contentType string, byteSize int64) (string, string, int64, error) {
	fileName = strings.TrimSpace(fileName)
	if fileName == "" || len(fileName) > 255 {
		return "", "", 0, invalidArgument("file_name is invalid")
	}

	extension := extensionFor(fileName, contentType)
	if extension == "" {
		return "", "", 0, errs.New(errs.CodeFileInvalidExtension, "file extension is invalid")
	}

	allowed := normalizeAllowedExtensions(allowedExtensions)
	if len(allowed) == 0 {
		return "", "", 0, errs.New(errs.CodeFileInvalidExtension, "file extension is not allowed")
	}
	if !slices.Contains(allowed, extension) {
		return "", "", 0, errs.New(errs.CodeFileInvalidExtension, "file extension is not allowed")
	}

	normalizedContentType := resolveContentType(contentType, extension)
	if len(normalizedContentType) > 128 {
		return "", "", 0, invalidArgument("content_type is invalid")
	}

	if maxUploadBytes <= 0 {
		maxUploadBytes = hardcodedMaxBytes
	}
	if byteSize <= 0 || byteSize > maxUploadBytes {
		return "", "", 0, errs.New(errs.CodeFileTooLarge, "file byte_size is invalid")
	}
	return normalizedContentType, extension, maxUploadBytes, nil
}

func extensionFor(fileName string, contentType string) string {
	ext := normalizeExtension(path.Ext(strings.TrimSpace(fileName)))
	if ext == "" {
		extensions, _ := mime.ExtensionsByType(normalizeContentType(contentType))
		if len(extensions) > 0 {
			ext = normalizeExtension(extensions[0])
		}
	}
	return ext
}

func resolveContentType(contentType string, extension string) string {
	normalized := normalizeContentType(contentType)
	if normalized != "" {
		return normalized
	}
	if inferred := normalizeContentType(mime.TypeByExtension(extension)); inferred != "" {
		return inferred
	}
	return "application/octet-stream"
}

func normalizeDisplayName(displayName string) (string, error) {
	displayName = strings.TrimSpace(displayName)
	if displayName == "" || len(displayName) > 128 {
		return "", invalidArgument("display_name is invalid")
	}
	return displayName, nil
}

func normalizeDescription(description string) (string, error) {
	description = strings.TrimSpace(description)
	if len(description) > 2048 {
		return "", invalidArgument("description is invalid")
	}
	return description, nil
}
