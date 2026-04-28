package file

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/file/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

func TestUploadCORSAllowsConfiguredOrigin(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodOptions, "/files/uploads/upload_1", http.NoBody)
	req.Header.Set("Origin", "http://127.0.0.1:5173")
	rec := httptest.NewRecorder()

	ok := addUploadCORS(rec, req, config.LocalStorageConf{
		AllowedOrigins: []string{"http://127.0.0.1:5173"},
	})

	if !ok {
		t.Fatal("expected configured upload origin to pass")
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://127.0.0.1:5173" {
		t.Fatalf("expected allowed origin header, got %q", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Headers"); !strings.Contains(got, "X-Upload-Token") {
		t.Fatalf("expected upload token header to be allowed, got %q", got)
	}
	if got := rec.Header().Get("Vary"); got != "Origin" {
		t.Fatalf("expected Vary: Origin, got %q", got)
	}
}

func TestUploadCORSRejectsUnconfiguredOrigin(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodOptions, "/files/uploads/upload_1", http.NoBody)
	req.Header.Set("Origin", "https://evil.example")
	rec := httptest.NewRecorder()

	ok := addUploadCORS(rec, req, config.LocalStorageConf{
		AllowedOrigins: []string{"http://127.0.0.1:5173"},
	})

	if ok {
		t.Fatal("expected unconfigured upload origin to fail")
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no allowed origin header, got %q", got)
	}
}

func TestPublicReadCORSAllowsAnyOrigin(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	addPublicReadCORS(rec)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Fatalf("expected public read CORS wildcard, got %q", got)
	}
	if got := rec.Header().Get("Access-Control-Allow-Methods"); strings.Contains(got, "PUT") {
		t.Fatalf("expected public read CORS to exclude PUT, got %q", got)
	}
}

func TestFileUploadPutReqParsesHeaderToken(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPut, "/files/uploads/upload_1", http.NoBody)
	req.Header.Set("X-Upload-Token", "signed-token")
	req = pathvar.WithVars(req, map[string]string{"upload_id": "upload_1"})

	var parsed types.FileUploadPutReq
	if err := httpx.Parse(req, &parsed); err != nil {
		t.Fatalf("expected request to parse, got %v", err)
	}
	if parsed.UploadId != "upload_1" || parsed.Token != "signed-token" {
		t.Fatalf("unexpected parsed request: %+v", parsed)
	}
}

func TestClientStreamAbortDetection(t *testing.T) {
	t.Parallel()

	if !isClientStreamAbort(context.Canceled) {
		t.Fatal("expected context cancellation to be treated as client abort")
	}
	if !isClientStreamAbort(net.ErrClosed) {
		t.Fatal("expected closed network connection to be treated as client abort")
	}
	if isClientStreamAbort(errors.New("write tcp: broken pipe")) {
		t.Fatal("expected plain error text not to be treated as client abort")
	}
	if isClientStreamAbort(nil) {
		t.Fatal("expected nil error not to be client abort")
	}
}
