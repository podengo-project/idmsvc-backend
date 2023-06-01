package errors

import (
	"fmt"
	"path/filepath"
	"runtime"
)

type LocationError struct {
	err  error
	file string
	line int
}

func (e LocationError) Error() string {
	return fmt.Sprintf("%s:%d - %v", e.file, e.line, e.err)
}

func NewLocationError(err error) LocationError {
	return NewLocationErrorWithLevel(err, 1)
}

func NewLocationErrorWithLevel(err error, level uint) LocationError {
	_, file, line, _ := runtime.Caller(1 + int(level))
	return LocationError{
		err:  err,
		file: filepath.Base(file),
		line: line,
	}
}
