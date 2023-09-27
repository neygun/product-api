package log

import "log"

type Logger interface {
	Printf(msg string, args ...interface{})
}

type loggerImpl struct {
	*log.Logger
}

func New() Logger {
	l := log.Default()
	return loggerImpl{
		Logger: l,
	}
}

var loggerSingleton Logger

func GetLogger() Logger {
	if loggerSingleton == nil {
		loggerSingleton = New()
	}

	return loggerSingleton
}
