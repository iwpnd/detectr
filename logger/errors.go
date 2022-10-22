package logger

import "fmt"

type ErrInvalidLogLevel struct {
	Level string
}

func (err ErrInvalidLogLevel) Error() string {
	return fmt.Sprintf("%s is an invalid log level", err.Level)
}
