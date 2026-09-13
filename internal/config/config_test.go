package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadValidConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(`
stores:
  aws:
    - regions:
        - us-east-1
  vault:
    enabled: false
  github:
    enabled: false
rotation_policy:
  default_max_age_days: 60
`), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Stores.AWS) != 1 {
		t.Errorf("expected 1 AWS entry, got %d", len(cfg.Stores.AWS))
	}
	if cfg.Stores.AWS[0].Regions[0] != "us-east-1" {
		t.Errorf("unexpected region: %s", cfg.Stores.AWS[0].Regions[0])
	}
	if cfg.RotationPolicy.DefaultMaxAgeDays != 60 {
		t.Errorf("expected max_age_days 60, got %d", cfg.RotationPolicy.DefaultMaxAgeDays)
	}
}

func TestLoadRejectsOldObjectShape(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	// Old broken shape: aws as a map, not a list
	if err := os.WriteFile(path, []byte(`
stores:
  aws:
    enabled: true
    regions:
      - us-east-1
`), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected unmarshal error for old object-shaped aws config, got nil")
	}
	if !strings.Contains(err.Error(), "unmarshal") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestTildeExpansionInDBPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(`
db_path: ~/.secretwatch/cache.db
`), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	home, _ := os.UserHomeDir()
	want := filepath.Join(home, ".secretwatch", "cache.db")
	if cfg.DBPath != want {
		t.Errorf("DBPath = %q, want %q", cfg.DBPath, want)
	}
}

func TestDefaultDBPathNeverContainsTilde(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	// No db_path in config — default should be applied with real home dir
	if err := os.WriteFile(path, []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.HasPrefix(cfg.DBPath, "~") {
		t.Errorf("default DBPath contains unexpanded tilde: %s", cfg.DBPath)
	}
}

func TestDefaultWebHost(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}

	os.Unsetenv("SECRETWATCH_WEB_HOST")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Web.Host != "127.0.0.1" {
		t.Errorf("Web.Host = %q, want 127.0.0.1", cfg.Web.Host)
	}
}

func TestWebHostEnvOverride(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(`{}`), 0600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("SECRETWATCH_WEB_HOST", "0.0.0.0")
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Web.Host != "0.0.0.0" {
		t.Errorf("Web.Host = %q, want 0.0.0.0", cfg.Web.Host)
	}
}
