package warehouse

type FactTransaction struct {
	ID              int64
	PatientID       int64
	MenuID          int64
	DietID          int64
	RoomID          int64
	AccountID       int64
	TimeID          int64
	TotalHarga      float64
	TotalItem       int
	StatusTransaksi int
}
