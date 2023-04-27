package http

import "github.com/lonelypale/gopkg/errors/status"

type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Details []any  `json:"details,omitempty"`
}

type Message struct {
	Success bool   `json:"success"`         //是否成功
	Data    any    `json:"data,omitempty"`  //结果数据
	Error   *Error `json:"error,omitempty"` //错误信息
}

func NewMessage(vals ...any) *Message {
	var success bool
	var data any
	var err *Error

	switch len(vals) {
	case 1:
		switch v := vals[0].(type) {
		case Error:
			err = &v
		case *Error:
			err = v
		case status.Status:
			err = &Error{
				Code:    v.Code(),
				Message: v.Message(),
				Details: v.Details(),
			}
		case error:
			err = &Error{
				Message: v.Error(),
			}
		default:
			success = true
			data = v
		}
	default:
		success = true
		if vals != nil {
			data = vals
		}
	}

	return &Message{
		Success: success,
		Data:    data,
		Error:   err,
	}
}
