package warehouse

import "time"

type Diet struct {
	ID          int64
	Name        string
	Description string
}

type Room struct {
	ID   int64
	Name string
}

type User struct {
	ID     int64
	Name   string
	Email  string
	RoleID int64
}

type Menu struct {
	ID       int64
	Name     string
	Calories int
	Price    int
	Status   int
}

type Patient struct {
	ID            int64
	Name          string
	RoomID        int64
	Gender        int
	Diagnosis     string
	AdmissionDate time.Time
	Age           int
	BirthDate     time.Time
	NIK           string
	BirthPlace    string
	DietID        int64
}

type TimeDim struct {
	ID            int64
	FullTimestamp time.Time
	Date          time.Time
	Day           string
	Month         string
	Year          string
	Hour          string
}
