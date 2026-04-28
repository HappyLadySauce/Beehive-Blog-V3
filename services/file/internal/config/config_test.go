package config

import (
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
)

func TestExampleConfigLoadsAndValidates(t *testing.T) {
	t.Parallel()

	var c Config
	if err := conf.Load("../../etc/file.yaml", &c); err != nil {
		t.Fatalf("expected example config to load, got %v", err)
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("expected example config to validate, got %v", err)
	}
	if len(c.Storage.Local.AllowedOrigins) == 0 {
		t.Fatal("expected local upload CORS allowed origins in example config")
	}
}
