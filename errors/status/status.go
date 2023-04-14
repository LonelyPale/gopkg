package status

import (
	"fmt"
	"strconv"
)

var (
	Success = Status{SuccessCode, "success", nil}
	Fail    = Status{FailCode, "fail", nil}
)

const (
	SuccessCode = 0
	FailCode    = 1
)

// ErrorCode Codes ecode error interface which has a code & message.
type ErrorCode interface {
	// Error sometimes Error return Code in string form
	// NOTE: don't use Error in monitor report even it also work for now
	Error() string
	// Code get error code.
	Code() int
	// Message get code message.
	Message() string
	// Details Detail get error detail,it may be nil.
	Details() []any
}

type Status struct {
	code    int
	message string
	details []any
}

func New(err error, details ...any) ErrorCode {
	switch v := err.(type) {
	case ErrorCode:
		return v
	default:
		return &Status{FailCode, v.Error(), details}
	}
}

// Error new status with code and message
func Error(code int, message string, details ...any) *Status {
	return &Status{code, message, details}
}

// Errorf new status with code and message
func Errorf(code int, format string, args ...interface{}) *Status {
	return &Status{code, fmt.Sprintf(format, args...), nil}
}

// Error implement error
func (s *Status) Error() string {
	return s.Message()
}

// Code return error code
func (s *Status) Code() int {
	return s.code
}

// Message return error message for developer
func (s *Status) Message() string {
	if s.message == "" {
		return strconv.Itoa(s.code)
	}
	return s.message
}

// Details return error details
func (s *Status) Details() []any {
	return s.details
}

// WithDetails WithDetails
func (s *Status) WithDetails(details ...any) *Status {
	s.details = append(s.details, details...)
	return s
}
