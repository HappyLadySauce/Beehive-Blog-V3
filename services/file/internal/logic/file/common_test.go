package file

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/storage"
)

func TestMapStorageReadError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want int
	}{
		{name: "not found", err: os.ErrNotExist, want: http.StatusNotFound},
		{name: "disabled", err: storage.ErrStorageDisabled, want: http.StatusServiceUnavailable},
		{name: "io failure", err: errors.New("disk read failed"), want: http.StatusServiceUnavailable},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mapped := mapStorageReadError(context.Background(), tt.err)
			dataErr, ok := mapped.(DataPlaneError)
			if !ok {
				t.Fatalf("expected DataPlaneError, got %T", mapped)
			}
			if dataErr.Status != tt.want {
				t.Fatalf("expected status %d, got %d", tt.want, dataErr.Status)
			}
		})
	}
}
