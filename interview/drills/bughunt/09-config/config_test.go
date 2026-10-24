//go:build bughunt

package config

import "testing"

func TestGetSeesLoadedConfig(t *testing.T) {
	go Load()

	cfg := Get()
	if cfg == nil || cfg.Addr != "db:5432" {
		t.Fatalf("Get = %+v, want Addr db:5432", cfg)
	}
}
