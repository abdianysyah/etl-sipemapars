package repository

import (
	"context"
	"database/sql"
	"fmt"

	warehouse "sipemapars-etl/internal/models/warehouse"
)

type execContext interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

type WarehouseRepository struct {
	db *sql.DB
}

func NewWarehouseRepository(db *sql.DB) *WarehouseRepository {
	return &WarehouseRepository{db: db}
}

func (r *WarehouseRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r *WarehouseRepository) Reset(ctx context.Context) error {
	stmts := []string{
		"SET FOREIGN_KEY_CHECKS=0",
		"TRUNCATE TABLE fact_transaction",
		"TRUNCATE TABLE dim_patient",
		"TRUNCATE TABLE dim_menu",
		"TRUNCATE TABLE dim_diet",
		"TRUNCATE TABLE dim_room",
		"TRUNCATE TABLE dim_users",
		"TRUNCATE TABLE dim_time",
		"SET FOREIGN_KEY_CHECKS=1",
	}

	for _, stmt := range stmts {
		if _, err := r.db.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("reset warehouse: %w", err)
		}
	}
	return nil
}

func (r *WarehouseRepository) InsertDiet(ctx context.Context, exec execContext, row warehouse.Diet) error {
	_, err := exec.ExecContext(ctx, `
		INSERT INTO dim_diet (id, name_diet, description, created_at, updated_at)
		VALUES (?, ?, ?, NOW(), NOW())
	`, row.ID, row.Name, row.Description)
	return err
}

func (r *WarehouseRepository) InsertRoom(ctx context.Context, exec execContext, row warehouse.Room) error {
	_, err := exec.ExecContext(ctx, `
		INSERT INTO dim_room (id, name_room, created_at, updated_at)
		VALUES (?, ?, NOW(), NOW())
	`, row.ID, row.Name)
	return err
}

func (r *WarehouseRepository) InsertUser(ctx context.Context, exec execContext, row warehouse.User) error {
	_, err := exec.ExecContext(ctx, `
		INSERT INTO dim_users (id, name, email, role_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`, row.ID, row.Name, row.Email, row.RoleID)
	return err
}

func (r *WarehouseRepository) InsertMenu(ctx context.Context, exec execContext, row warehouse.Menu) error {
	_, err := exec.ExecContext(ctx, `
		INSERT INTO dim_menu (id, nama_menu, kalori, harga, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`, row.ID, row.Name, row.Calories, row.Price, row.Status)
	return err
}

func (r *WarehouseRepository) InsertPatient(ctx context.Context, exec execContext, row warehouse.Patient) error {
	_, err := exec.ExecContext(ctx, `
		INSERT INTO dim_patient (
			id, nama_pasien, id_room, jenis_kelamin, diagnosis,
			tanggal_masuk, umur, tanggal_lahir, nik, tempat_lahir, id_diet
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, row.ID, row.Name, row.RoomID, row.Gender, row.Diagnosis, row.AdmissionDate, row.Age, row.BirthDate, row.NIK, row.BirthPlace, row.DietID)
	return err
}

func (r *WarehouseRepository) InsertTime(ctx context.Context, exec execContext, row warehouse.TimeDim) error {
	_, err := exec.ExecContext(ctx, `
		INSERT INTO dim_time (
			id, full_timestamp, tanggal, hari, bulan, tahun, jam, created_at, updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`, row.ID, row.FullTimestamp, row.Date, row.Day, row.Month, row.Year, row.Hour)
	return err
}

func (r *WarehouseRepository) InsertFact(ctx context.Context, exec execContext, row warehouse.FactTransaction) error {
	_, err := exec.ExecContext(ctx, `
		INSERT INTO fact_transaction (
			id, patient_id, menu_id, diet_id, room_id, account_id, time_id,
			total_harga, total_item, status_transaksi, created_at, updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`, row.ID, row.PatientID, row.MenuID, row.DietID, row.RoomID, row.AccountID, row.TimeID, row.TotalHarga, row.TotalItem, row.StatusTransaksi)
	return err
}
