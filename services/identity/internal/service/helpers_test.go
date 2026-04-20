package service_test

import (
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth/provider"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/testkit"
)

func newDeps(t *testing.T, now time.Time) service.Dependencies {
	t.Helper()
	return testkit.NewServiceDependencies(t, func() time.Time { return now.UTC() })
}

func githubClientFromDeps(t *testing.T, deps service.Dependencies) *provider.GitHubClient {
	t.Helper()

	callbackProvider, ok := deps.Providers.GetCallback("github")
	if !ok {
		t.Fatalf("expected github callback provider to be registered")
	}

	client, ok := callbackProvider.(*provider.GitHubClient)
	if !ok {
		t.Fatalf("expected github callback provider type, got %T", callbackProvider)
	}

	return client
}
