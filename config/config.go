package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SourceDSN string
	DWDSN     string
}

func Load() (Config, error) {

	// load .env
	_ = godotenv.Load()

	cfg := Config{
		SourceDSN: os.Getenv("SOURCE_DSN"),
		DWDSN:     os.Getenv("DW_DSN"),
	}

	if cfg.SourceDSN == "" {
		return Config{}, fmt.Errorf("SOURCE_DSN is required")
	}

	if cfg.DWDSN == "" {
		return Config{}, fmt.Errorf("DW_DSN is required")
	}

	return cfg, nil
}