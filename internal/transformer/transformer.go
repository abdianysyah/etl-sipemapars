package transformer

import (
	"time"

	source "sipemapars-etl/internal/models/source"
	warehouse "sipemapars-etl/internal/models/warehouse"
	"sipemapars-etl/internal/utils"
)

func Diet(id int64, row source.Diet) warehouse.Diet {
	return warehouse.Diet{
		ID:          id,
		Name:        row.Name,
		Description: row.Description.String,
	}
}

func Room(id int64, row source.Room) warehouse.Room {
	return warehouse.Room{
		ID:   id,
		Name: row.Name,
	}
}

func User(id int64, row source.User) warehouse.User {
	return warehouse.User{
		ID:     id,
		Name:   row.Name,
		Email:  row.Email,
		RoleID: row.RoleID.Int64,
	}
}

func Menu(id int64, row source.Menu) warehouse.Menu {
	status := row.Status
	if status != 1 {
		status = 0
	}
	return warehouse.Menu{
		ID:       id,
		Name:     row.Name,
		Calories: row.Calories,
		Price:    row.Price,
		Status:   status,
	}
}

func Patient(id int64, row source.Patient, roomID, dietID int64) warehouse.Patient {
	return warehouse.Patient{
		ID:            id,
		Name:          row.Name,
		RoomID:        roomID,
		Gender:        row.Gender,
		Diagnosis:     row.Diagnosis,
		AdmissionDate: row.AdmissionDate,
		Age:           row.Age,
		BirthDate:     row.BirthDate,
		NIK:           row.NIK,
		BirthPlace:    row.BirthPlace,
		DietID:        dietID,
	}
}

func TimeDim(id int64, ts time.Time) warehouse.TimeDim {
	return warehouse.TimeDim{
		ID:            id,
		FullTimestamp: ts,
		Date:          utils.DateOnly(ts),
		Day:           utils.IndonesianDayName(ts.Weekday()),
		Month:         utils.MonthNumber(ts),
		Year:          utils.YearString(ts),
		Hour:          utils.TimeOnlyString(ts),
	}
}

func FactTransaction(id int64, patientID, menuID, dietID, roomID, accountID, timeID int64, totalPrice int, totalItem, status int) warehouse.FactTransaction {
	return warehouse.FactTransaction{
		ID:              id,
		PatientID:       patientID,
		MenuID:          menuID,
		DietID:          dietID,
		RoomID:          roomID,
		AccountID:       accountID,
		TimeID:          timeID,
		TotalHarga:      float64(totalPrice),
		TotalItem:       totalItem,
		StatusTransaksi: status,
	}
}
