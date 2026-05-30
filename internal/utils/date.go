package utils

import (
	"fmt"
	"time"
)

func IndonesianDayName(wd time.Weekday) string {
	switch wd {
	case time.Sunday:
		return "Minggu"
	case time.Monday:
		return "Senin"
	case time.Tuesday:
		return "Selasa"
	case time.Wednesday:
		return "Rabu"
	case time.Thursday:
		return "Kamis"
	case time.Friday:
		return "Jumat"
	default:
		return "Sabtu"
	}
}

func MonthNumber(ts time.Time) string {
	return fmt.Sprintf("%02d", int(ts.Month()))
}

func YearString(ts time.Time) string {
	return fmt.Sprintf("%04d", ts.Year())
}

func TimeOnlyString(ts time.Time) string {
	return ts.Format("15:04:05")
}

func DateOnly(ts time.Time) time.Time {
	y, m, d := ts.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, ts.Location())
}

func TimestampKey(ts time.Time) string {
	return ts.UTC().Format(time.RFC3339Nano)
}
