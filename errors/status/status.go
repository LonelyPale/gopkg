package status

import (
	"fmt"
	"runtime"
	"strings"
)

var _ error = (*status)(nil)
var _ Status = (*status)(nil)

type Status interface {
	Error() string
	Unwrap() error

	Code() int
	Type() string
	Message() string
	Details() []any

	SetCode(code int) Status
	SetType(typ string) Status
	SetError(err error) Status
	SetMessage(message string) Status
	SetDetails(details ...any) Status
}

type status struct {
	code    int
	typ     string
	err     error
	message string
	details []any
	stack   *stack
}

func New(code int, typ string, message any, details ...any) Status {
	var msg string
	var err error
	switch val := message.(type) {
	case string:
		msg = val
	case error:
		err = val
	default:
		msg = fmt.Sprintf("%v", val)
	}
	return &status{
		code:    code,
		typ:     typ,
		err:     err,
		message: msg,
		details: details,
		stack:   callers(),
	}
}

// Error new status with message and details
func Error(message string, details ...any) Status {
	return &status{
		message: message,
		details: details,
		stack:   callers(),
	}
}

// Errorf new status with format message
func Errorf(format string, args ...interface{}) Status {
	return &status{
		message: fmt.Sprintf(format, args...),
		details: nil,
		stack:   callers(),
	}
}

func (s *status) SetCode(code int) Status {
	s.code = code
	return s
}

func (s *status) SetType(typ string) Status {
	s.typ = typ
	return s
}

func (s *status) SetError(err error) Status {
	s.err = err
	return s
}

func (s *status) SetMessage(message string) Status {
	s.message = message
	return s
}

func (s *status) SetDetails(details ...any) Status {
	s.details = details
	return s
}

func (s *status) Unwrap() error {
	return s.err
}

// Error implement error
func (s *status) Error() string {
	if s.message != "" && s.err != nil {
		return fmt.Sprintf("%s: %s\n%s", s.message, s.err, s.stack.String())
	} else if s.err != nil {
		return fmt.Sprintf("%s\n%s", s.err, s.stack.String())
	}
	return fmt.Sprintf("%s\n%s", s.message, s.stack.String())
}

// Code return error code
func (s *status) Code() int {
	return s.code
}

// Type return error typ
func (s *status) Type() string {
	return s.typ
}

// Message return error message for developer
func (s *status) Message() string {
	return s.message
}

// Details return error details
func (s *status) Details() []any {
	return s.details
}

// stack represents a stack of program counters.
type stack []uintptr

func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

func (s *stack) String() string {
	builder := new(strings.Builder)
	frames := runtime.CallersFrames(*s)
	for {
		frame, more := frames.Next()
		_, _ = fmt.Fprintf(builder, "%s:%d %s\n", frame.File, frame.Line, frame.Function)
		if !more {
			break
		}
	}
	return builder.String()
}
