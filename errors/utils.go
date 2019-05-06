package errors

import "github.com/erichnascimento/golib/runtime"

// Catcher is a type which catches errors silently
type Catcher func(error)

type CloseableWithError interface{
	Close() error
}

// Closer is a helper function which receives a closeable and close it sending eventual error to default
// registered error handler
func Closer(c CloseableWithError) {
	runtime.DefaultErrorHandler(c.Close())
}
