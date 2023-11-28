package utils

import "time"

const DATE_FORMAT = "2006/01/02"
const MINIMUM_YEAR = 1900
const DEFAULT_DAYS_TO_GO_BACK = 30

func GetDateFromDateToFromToday(days int) (time.Time, time.Time) {
	dateTo := time.Now()
	dateFrom := dateTo.AddDate(0, 0, -days)

	return dateFrom, dateTo
}

func ConvertStringDate(date string) (time.Time, error) {
	return time.Parse(DATE_FORMAT, date)
}

func IsDatesValid(dateFrom time.Time, dateTo time.Time) bool {
	return dateFrom.Year() >= MINIMUM_YEAR &&
		dateTo.Year() >= MINIMUM_YEAR &&
		dateFrom.Before(dateTo) &&
		dateTo.Before(time.Now())
}
