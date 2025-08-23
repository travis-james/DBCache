package config

import "testing"

func TestLoadConfigParallel(t *testing.T) {
	t.Setenv("POSTGRES_HOST", "localhost")
	t.Setenv("POSTGRES_USER", "user")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.DatastoreDBHost != "localhost" {
		t.Errorf("expected 'localhost', got '%s'", cfg.DatastoreDBHost)
	}
}
