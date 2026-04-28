//go:build !windows

package file

import (
	"net"
	"os"
	"syscall"
	"testing"
)

func TestClientStreamAbortDetectsUnixSyscall(t *testing.T) {
	t.Parallel()

	err := &net.OpError{
		Op:  "write",
		Net: "tcp",
		Err: &os.SyscallError{Syscall: "write", Err: syscall.EPIPE},
	}

	if !isClientStreamAbort(err) {
		t.Fatal("expected EPIPE to be treated as client abort")
	}
}
