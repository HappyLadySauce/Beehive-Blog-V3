package config

import (
	"fmt"
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

func TestLocalStorageValidateRequiresUploadSecret(t *testing.T) {
	t.Parallel()

	conf := LocalStorageConf{
		ListenOn:      "127.0.0.1:8084",
		RootDir:       t.TempDir(),
		TempDir:       t.TempDir(),
		Bucket:        "beehive",
		UploadBaseURL: "http://127.0.0.1:8084/files/uploads",
		UploadSecret:  " ",
	}

	if err := conf.Validate(); err == nil || fmt.Sprint(err) != "Storage.Local.UploadSecret is required" {
		t.Fatalf("expected UploadSecret validation error, got %v", err)
	}
}
