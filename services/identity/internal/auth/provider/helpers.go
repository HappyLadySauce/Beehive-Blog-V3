package provider

import (
	"regexp"
	"strings"
)

// scopesPtr joins scopes into a comma-separated pointer value.
// scopesPtr 将 scopes 拼接为逗号分隔的指针值。
func scopesPtr(scopes []string) *string {
	if len(scopes) == 0 {
		return nil
	}

	value := strings.Join(scopes, ",")
	return &value
}

// stringPtr returns a trimmed string pointer for non-empty values.
// stringPtr 为非空值返回裁剪后的字符串指针。
func stringPtr(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

// defaultScopes returns configured scopes or falls back to defaults.
// defaultScopes 返回配置 scopes，若为空则回退到默认值。
func defaultScopes(configured []string, defaults []string) []string {
	configured = trimmedScopes(configured)
	if len(configured) > 0 {
		return configured
	}

	return defaults
}

// scopeStringPtr returns a trimmed scope string pointer for non-empty values.
// scopeStringPtr 为非空 scope 字符串返回裁剪后的指针。
func scopeStringPtr(value string) *string {
	trimmed := strings.Trim(strings.TrimSpace(value), ",")
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

var sensitivePreviewPattern = regexp.MustCompile(`(?i)"(access_token|refresh_token|client_secret|token)"\s*:\s*"[^"]*"`)

// trimmedScopes removes empty scope values and trims whitespace.
// trimmedScopes 去除空 scope 并裁剪空白。
func trimmedScopes(scopes []string) []string {
	if len(scopes) == 0 {
		return nil
	}

	result := make([]string, 0, len(scopes))
	for _, scope := range scopes {
		trimmed := strings.TrimSpace(scope)
		if trimmed == "" {
			continue
		}
		result = append(result, trimmed)
	}

	return result
}

// bodyPreview returns a masked, truncated response body preview for diagnostics.
// bodyPreview 返回用于诊断的脱敏且截断后的响应体预览。
func bodyPreview(body []byte) string {
	preview := strings.TrimSpace(string(body))
	if preview == "" {
		return ""
	}

	preview = sensitivePreviewPattern.ReplaceAllStringFunc(preview, func(match string) string {
		parts := strings.SplitN(match, ":", 2)
		if len(parts) != 2 {
			return `"masked":"***"`
		}
		return parts[0] + `:"***"`
	})
	if len(preview) > 256 {
		return preview[:256]
	}

	return preview
}
