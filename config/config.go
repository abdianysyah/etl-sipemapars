package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SourceDSN		string
	DWDSN			string
	APIAddr			string
	LaravelBaseURL	string
	CallbackSecret	string
}

func Load() (Config, error) {

	// load .env
	_ = godotenv.Load()

	cfg := Config{
		SourceDSN: 		os.Getenv("SOURCE_DSN"),
		DWDSN:     		os.Getenv("DW_DSN"),
		APIAddr:		os.Getenv("API_ADDR"),
		LaravelBaseURL:	os.Getenv("LARAVEL_BASE_URL"),
		CallbackSecret: os.Getenv("ETL_CALLBACK_SECRET"),
	}

	if cfg.SourceDSN == "" {
		return Config{}, fmt.Errorf("SOURCE_DSN is required")
	}

	if cfg.DWDSN == "" {
		return Config{}, fmt.Errorf("DW_DSN is required")
	}

	if cfg.APIAddr == "" {
		cfg.APIAddr = ": 9000"
	}

	return cfg, nil
}