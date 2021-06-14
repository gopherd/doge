package erron

import (
	"errors"
	"strconv"
)

const (
	EUnknown = -1
	EOK      = 0

	// User-defined errno should be greater than zero
)

type errno struct {
	code int
	err  error
}

func (err *errno) Errno() int {
	return err.code
}

func (err *errno) Error() string {
	return "(" + strconv.Itoa(err.code) + ") " + err.err.Error()
}

func (err *errno) Unwrap() error {
	return err.err
}

// WithErrno wraps the error with code
func WithErrno(code int, err error) error {
	if err == nil {
		return nil
	}
	return &errno{
		code: code,
		err:  err,
	}
}

// WithErrnof returns an error that formats as the given text with code.
func WithErrnof(code int, format string, args ...interface{}) error {
	return &errno{
		code: code,
		err:  New(format, args...),
	}
}

// Errno finds the first error in err's chain that contains errno.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
func Errno(err error) int {
	if err == nil {
		return EOK
	}
	for {
		if e, ok := err.(interface{ Errno() int }); ok {
			return e.Errno()
		}
		if err = errors.Unwrap(err); err == nil {
			break
		}
	}
	return EUnknown
}
