package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"

	"sipemapars-etl/config"
	"sipemapars-etl/internal/callback"
	"sipemapars-etl/internal/repository"
	"sipemapars-etl/internal/service"
)

type Runner struct {
	mu		sync.Mutex
	running bool
	etl 	*service.Service
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
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

	cb := callback.New(cfg.LaravelBaseURL, cfg.CallbackSecret)
	etl := service.New(sourceRepo, warehouseRepo, cb)

	runner := &Runner{etl: etl}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /api/etl/run", runner.runHandler)

	addr := cfg.APIAddr
	if addr == "" {
		addr = ":9000"
	}

	log.Println("ETL API listening on", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	_, _= w.Write([]byte(`{"status":"ok"}`))
}

func (r *Runner) runHandler(w http.ResponseWriter, req *http.Request)  {
	r.mu.Lock()
	if r.running {
		r.mu.Unlock()
		http.Error(w, "etl already running", http.StatusConflict)
		return
	}
	r.running = true
	r.mu.Unlock()

	type runRequest struct {
		JobUUID		string `json:"job_uuid"`
		Name		string `json:"name"`
		TargetTable string `json:"target_table"`
	}

	var body runRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		r.mu.Lock()
		r.running = false
		r.mu.Unlock()
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if body.JobUUID == "" {
		r.mu.Lock()
		r.running = false
		r.mu.Unlock()
		http.Error(w, "job_uuid is required", http.StatusBadRequest)
		return
	}

	go func(jobUUID string)  {
		defer func()  {
			r.mu.Lock()
			r.running = false
			r.mu.Unlock()
		}()

		if err := r.etl.Run(context.Background(), jobUUID); err != nil {
			log.Println("etl failed:", err)
		}
	}(body.JobUUID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write([]byte(`{"message":"accepted"}`))
}
