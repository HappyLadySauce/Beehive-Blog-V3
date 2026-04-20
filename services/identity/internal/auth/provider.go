package auth

import (
	"strings"

	"github.com/google/uuid"
)

// EnsureState returns the provided state or generates a new one when it is empty.
// EnsureState 返回传入 state，若为空则自动生成。
func EnsureState(state string) string {
	if trimmed := strings.TrimSpace(state); trimmed != "" {
		return trimmed
	}

	return uuid.NewString()
}
