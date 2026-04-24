package handler

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestRegisterHandlersKeepsManualMiddlewareGuards(t *testing.T) {
	source := readRoutesSource(t)

	requiredSnippets := []string{
		"rest.WithMiddleware(serverCtx.RequestMetaMiddleware, authPublicRoutes...)",
		"serverCtx.RequestMetaMiddleware,\n\t\t\tserverCtx.AuthMiddleware",
		"}, authProtectedRoutes...)",
		"}, studioContentRoutes...)",
		"rest.WithMiddleware(serverCtx.RequestMetaMiddleware, opsRoutes...)",
		"rest.WithMiddleware(serverCtx.RequestMetaMiddleware, publicContentRoutes...)",
		"Path:    \"/items/:content_id/relations\"",
		"Path:    \"/items/:content_id/relations/:relation_id\"",
	}

	for _, snippet := range requiredSnippets {
		if !strings.Contains(source, snippet) {
			t.Fatalf("routes.go lost required middleware or route snippet %q", snippet)
		}
	}

	if got := strings.Count(source, "serverCtx.AuthMiddleware"); got != 2 {
		t.Fatalf("protected auth middleware should be registered for auth and studio routes, got %d registrations", got)
	}
}

func readRoutesSource(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve current test file")
	}

	routesPath := filepath.Join(filepath.Dir(currentFile), "routes.go")
	source, err := os.ReadFile(routesPath)
	if err != nil {
		t.Fatalf("read routes.go: %v", err)
	}

	return string(source)
}
