package utils

import (
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func NewLogger() log.Logger {
	var logger log.Logger
	{
		logger = log.NewJSONLogger(os.Stdout)
		logger = log.With(logger, "timestamp", log.DefaultTimestampUTC, "caller", log.Caller(4))
		logger = level.NewFilter(logger, level.AllowAll())
	}

	return logger
}
