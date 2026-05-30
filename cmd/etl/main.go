package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"sipemapars-etl/config"
	"sipemapars-etl/internal/repository"
	"sipemapars-etl/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "config error:", err)
		os.Exit(1)
	}

	sourceDB, err := sql.Open("mysql", cfg.SourceDSN)
	if err != nil {
		fmt.Fprintln(os.Stderr, "source db open error:", err)
		os.Exit(1)
	}
	defer sourceDB.Close()

	dwDB, err := sql.Open("mysql", cfg.DWDSN)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dw db open error:", err)
		os.Exit(1)
	}
	defer dwDB.Close()

	sourceDB.SetMaxOpenConns(5)
	sourceDB.SetMaxIdleConns(5)
	dwDB.SetMaxOpenConns(5)
	dwDB.SetMaxIdleConns(5)

	if err := sourceDB.Ping(); err != nil {
		fmt.Fprintln(os.Stderr, "source db ping error:", err)
		os.Exit(1)
	}
	if err := dwDB.Ping(); err != nil {
		fmt.Fprintln(os.Stderr, "dw db ping error:", err)
		os.Exit(1)
	}

	sourceRepo := repository.NewSourceRepository(sourceDB)
	warehouseRepo := repository.NewWarehouseRepository(dwDB)
	etl := service.New(sourceRepo, warehouseRepo)

	if err := etl.Run(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, "etl failed:", err)
		os.Exit(1)
	}
}
