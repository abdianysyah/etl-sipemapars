# Sipemapars ETL (Go)

ETL service untuk memindahkan data dari database operasional `sipemapars` ke data warehouse `sipemapars_dw`.

## Cara jalan

1. Copy `.env.example` menjadi `.env`
2. Sesuaikan DSN MySQL
3. Jalankan:

```bash
go mod tidy
go run ./cmd/etl
```

## Catatan

- ETL ini memakai full refresh karena schema DW yang ada tidak menyimpan kolom `source_*`.
- Semua tabel DW di-truncate dulu, lalu diisi ulang dari source.
- Faktanya:
  - `patient_id` diambil dari mapping pasien source -> DW
  - `menu_id` dipilih dari menu yang sesuai diet pasien jika ada, kalau tidak pakai record Unknown
  - `diet_id` dan `room_id` diambil dari pasien
  - `account_id` memakai Unknown User karena source transaction tidak menyimpan user pembuat transaksi
  - `total_item` di-set 1 karena source transaksi tidak punya detail item menu

## Struktur

- `internal/models/source`  -> model tabel operasional
- `internal/models/warehouse` -> model tabel DW
- `internal/repository` -> akses database
- `internal/service` -> orchestration ETL
- `internal/transformer` -> mapping data
