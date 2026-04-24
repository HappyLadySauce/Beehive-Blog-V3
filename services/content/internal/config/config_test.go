package config

import "testing"

func TestValidateRequiresInfrastructureAndInternalAuth(t *testing.T) {
	t.Parallel()

	var c Config
	if err := c.Validate(); err == nil {
		t.Fatalf("expected empty config to fail")
	}

	c.Postgres.Host = "127.0.0.1"
	c.Postgres.User = "postgres"
	c.Postgres.DBName = "beehive"
	c.InternalAuthToken = "secret"
	c.AllowedCallers = []string{"gateway"}
	if err := c.Validate(); err != nil {
		t.Fatalf("expected valid config, got %v", err)
	}
}
