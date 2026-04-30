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
	if c.Storage.Local.Bucket != "beehive-local" {
		t.Fatal("expected local bucket in example config")
	}
}

func TestLocalStorageValidateRequiresRootDir(t *testing.T) {
	t.Parallel()

	conf := LocalStorageConf{
		RootDir: "",
		TempDir: t.TempDir(),
		Bucket:  "beehive",
	}

	if err := conf.Validate(); err == nil {
		t.Fatal("expected RootDir validation error")
	}
}
