//go:build windows

package file

import (
	"errors"
	"syscall"
)

func isStreamAbortSyscall(err error) bool {
	return errors.Is(err, syscall.WSAECONNRESET) ||
		errors.Is(err, syscall.WSAECONNABORTED)
}
