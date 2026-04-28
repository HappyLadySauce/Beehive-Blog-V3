//go:build windows

package file

import (
	"net"
	"os"
	"syscall"
	"testing"
)

func TestClientStreamAbortDetectsWindowsSyscall(t *testing.T) {
	t.Parallel()

	err := &net.OpError{
		Op:  "write",
		Net: "tcp",
		Err: &os.SyscallError{Syscall: "wsasend", Err: syscall.WSAECONNRESET},
	}

	if !isClientStreamAbort(err) {
		t.Fatal("expected WSAECONNRESET to be treated as client abort")
	}
}
