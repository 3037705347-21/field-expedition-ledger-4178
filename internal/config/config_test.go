package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("LEDGER_ADDR", "")
	t.Setenv("LEDGER_SEED", "")
	value := Load()
	if value.Address == "" || !value.Seed {
		t.Fatalf("unexpected defaults: %+v", value)
	}
}
