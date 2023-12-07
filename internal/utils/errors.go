package utils

import (
	"errors"
)

var (
	ErrDateInvalidFormat = errors.New("invalid date format. Should be in format: YYYY/MM/DD")
	ErrInvalidDateRange  = errors.New("invalid date range")
	ErrNotFound          = errors.New("not found")
	ErrPeriodicity       = errors.New("periodicity not supported")
	ErrDecodeRequest     = errors.New("unable to decode the request")
)
