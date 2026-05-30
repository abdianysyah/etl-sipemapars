package source

import (
	"database/sql"
	"time"
)

type Diet struct {
	ID          int64          `db:"id"`
	Name        string         `db:"nama"`
	Description sql.NullString `db:"deskripsi"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
}

type Room struct {
	ID        int64        `db:"id"`
	Name      string       `db:"nama_room"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

type User struct {
	ID        int64         `db:"id"`
	Name      string        `db:"name"`
	Email     string        `db:"email"`
	RoleID    sql.NullInt64 `db:"role_id"`
	CreatedAt sql.NullTime  `db:"created_at"`
	UpdatedAt sql.NullTime  `db:"updated_at"`
}

type Menu struct {
	ID        string       `db:"id"`
	Name      string       `db:"nama_menu"`
	DietID    int64        `db:"id_diet"`
	Calories  int          `db:"kalori"`
	Price     int          `db:"harga"`
	Status    int          `db:"status"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

type Patient struct {
	ID            string       `db:"id"`
	Name          string       `db:"nama_pasien"`
	RoomID        int64        `db:"id_room"`
	Gender        int          `db:"jenis_kelamin"`
	Diagnosis     string       `db:"diagnosis"`
	AdmissionDate time.Time    `db:"tanggal_masuk"`
	Age           int          `db:"umur"`
	BirthDate     time.Time    `db:"tanggal_lahir"`
	NIK           string       `db:"nik"`
	BirthPlace    string       `db:"tempat_lahir"`
	DietID        int64        `db:"id_diet"`
	CreatedAt     sql.NullTime `db:"created_at"`
	UpdatedAt     sql.NullTime `db:"updated_at"`
}

type Transaction struct {
	ID         string       `db:"id"`
	PatientID  string       `db:"patient_id"`
	OrderDate  time.Time    `db:"order_date"`
	TotalPrice int          `db:"total_price"`
	Status     int          `db:"status"`
	CreatedAt  sql.NullTime `db:"created_at"`
	UpdatedAt  sql.NullTime `db:"updated_at"`
}
