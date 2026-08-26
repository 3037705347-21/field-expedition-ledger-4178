package config

import "os"

type Config struct {
	Address string
	Seed    bool
}

func Load() Config {
	address := os.Getenv("LEDGER_ADDR")
	if address == "" {
		address = "127.0.0.1:8090"
	}
	seed := os.Getenv("LEDGER_SEED") != "false"
	return Config{Address: address, Seed: seed}
}
