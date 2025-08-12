package config

import "testing"

func TestFromEnv(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost/db")
	t.Setenv("PORT", "9090")
	cfg := FromEnv()
	if cfg.DatabaseURL != "postgres://user:pass@localhost/db" {
		t.Fatalf("unexpected DatabaseURL: %s", cfg.DatabaseURL)
	}
	if cfg.Port != "9090" {
		t.Fatalf("unexpected Port: %s", cfg.Port)
	}
}

func TestFromEnvDefaultPort(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("PORT", "")
	cfg := FromEnv()
	if cfg.Port != "8080" {
		t.Fatalf("expected default port 8080, got %s", cfg.Port)
	}
}
