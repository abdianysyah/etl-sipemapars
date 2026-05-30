package repository

import (
	"context"
	"database/sql"
	"fmt"

	source "sipemapars-etl/internal/models/source"
)

type SourceRepository struct {
	db *sql.DB
}

func NewSourceRepository(db *sql.DB) *SourceRepository {
	return &SourceRepository{db: db}
}

func scanDiets(rows *sql.Rows) ([]source.Diet, error) {
	var out []source.Diet
	for rows.Next() {
		var row source.Diet
		if err := rows.Scan(&row.ID, &row.Name, &row.Description, &row.CreatedAt, &row.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func (r *SourceRepository) GetDiets(ctx context.Context) ([]source.Diet, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, nama, COALESCE(deskripsi, '') AS deskripsi, created_at, updated_at
		FROM diets
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("query diets: %w", err)
	}
	defer rows.Close()
	return scanDiets(rows)
}

func (r *SourceRepository) GetRooms(ctx context.Context) ([]source.Room, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, nama_room, created_at, updated_at
		FROM rooms
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("query rooms: %w", err)
	}
	defer rows.Close()

	var out []source.Room
	for rows.Next() {
		var row source.Room
		if err := rows.Scan(&row.ID, &row.Name, &row.CreatedAt, &row.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func (r *SourceRepository) GetUsers(ctx context.Context) ([]source.User, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, email, COALESCE(role_id, 0) AS role_id, created_at, updated_at
		FROM users
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	var out []source.User
	for rows.Next() {
		var row source.User
		if err := rows.Scan(&row.ID, &row.Name, &row.Email, &row.RoleID, &row.CreatedAt, &row.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func (r *SourceRepository) GetMenus(ctx context.Context) ([]source.Menu, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, nama_menu, id_diet, kalori, harga, status, created_at, updated_at
		FROM menus
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("query menus: %w", err)
	}
	defer rows.Close()

	var out []source.Menu
	for rows.Next() {
		var row source.Menu
		if err := rows.Scan(&row.ID, &row.Name, &row.DietID, &row.Calories, &row.Price, &row.Status, &row.CreatedAt, &row.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func (r *SourceRepository) GetPatients(ctx context.Context) ([]source.Patient, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, nama_pasien, id_room, jenis_kelamin, diagnosis, tanggal_masuk, umur,
		       tanggal_lahir, nik, tempat_lahir, id_diet, created_at, updated_at
		FROM patients
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("query patients: %w", err)
	}
	defer rows.Close()

	var out []source.Patient
	for rows.Next() {
		var row source.Patient
		if err := rows.Scan(&row.ID, &row.Name, &row.RoomID, &row.Gender, &row.Diagnosis, &row.AdmissionDate, &row.Age, &row.BirthDate, &row.NIK, &row.BirthPlace, &row.DietID, &row.CreatedAt, &row.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func (r *SourceRepository) GetTransactions(ctx context.Context) ([]source.Transaction, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, patient_id, order_date, total_price, status, created_at, updated_at
		FROM transactions
		ORDER BY order_date ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("query transactions: %w", err)
	}
	defer rows.Close()

	var out []source.Transaction
	for rows.Next() {
		var row source.Transaction
		if err := rows.Scan(&row.ID, &row.PatientID, &row.OrderDate, &row.TotalPrice, &row.Status, &row.CreatedAt, &row.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}
