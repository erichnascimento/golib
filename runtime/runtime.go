package runtime

import "github.com/apex/log"

var logger = log.Log

func Logger() log.Interface {
	return logger
}

// SetLogger set the default logger
// This method does not prevent race condition
func SetLogger(newLogger log.Interface) {
	if newLogger == nil {
		panic("logger can not be nil")
	}
	logger = newLogger
}

var DefaultErrorHandler = func(err error) {
	if err != nil {
		Logger().Errorf("error caught by default error handler: %v", err)
	}
}
