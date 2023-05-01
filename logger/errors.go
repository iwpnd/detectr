package logger

import "fmt"

// ErrInvalidLogLevel ...
type ErrInvalidLogLevel struct {
	Level string
}

// Error ...
func (err ErrInvalidLogLevel) Error() string {
	return fmt.Sprintf("%s is an invalid log level", err.Level)
}
