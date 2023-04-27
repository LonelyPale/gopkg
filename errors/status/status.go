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

	Code() string
	Message() string
	Details() []any

	SetCode(code string) Status
	SetMessage(message string) Status
	SetDetails(details ...any) Status
	SetError(err error) Status
}

type status struct {
	code    string // error code
	message string // error message
	details []any  // error details
	err     error  // error interface
	stack   *stack // error call stack
}

func New(code string, message any, details ...any) Status {
	if code == "" && message == nil {
		return nil
	} else if message == nil {
		message = ""
	}

	var msg string
	var err error
	switch val := message.(type) {
	case Status:
		return &status{
			code:    code,
			details: details,
			err:     val,
		}
	case string:
		msg = val
	case error:
		err = val
	default:
		msg = fmt.Sprintf("%v", val)
	}

	return &status{
		code:    code,
		message: msg,
		details: details,
		err:     err,
		stack:   callers(),
	}
}

// Error new status with message and details
func Error(message any, details ...any) Status {
	if message == nil {
		return nil
	}

	var msg string
	var err error
	switch val := message.(type) {
	case Status:
		if details == nil {
			return val
		}
		return &status{
			details: details,
			err:     val,
		}
	case string:
		msg = val
	case error:
		err = val
	default:
		msg = fmt.Sprintf("%v", val)
	}

	return &status{
		message: msg,
		details: details,
		err:     err,
		stack:   callers(),
	}
}

// Errorf new status with format message
func Errorf(format string, args ...interface{}) Status {
	return &status{
		message: fmt.Sprintf(format, args...),
		stack:   callers(),
	}
}

func (s *status) SetCode(code string) Status {
	s.code = code
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

func (s *status) SetError(err error) Status {
	s.err = err
	return s
}

func (s *status) Unwrap() error {
	return s.err
}

// Error implement error
func (s *status) Error() string {
	var pre bool
	var out strings.Builder
	if s.code != "" {
		out.WriteString(s.code)
		out.WriteString(": ")
		pre = true
	}
	if s.message != "" {
		out.WriteString(s.message)
		pre = true
	}
	if s.err != nil {
		if pre {
			out.WriteString("\n")
		}
		out.WriteString(s.err.Error())
		pre = true
	}
	if s.stack != nil {
		if pre {
			out.WriteString("\n")
		}
		out.WriteString(s.stack.String())
	}
	return out.String()
}

// Code return error code
func (s *status) Code() string {
	return s.code
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
