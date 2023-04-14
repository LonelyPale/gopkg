package http

import "github.com/lonelypale/gopkg/errors/status"

type Message struct {
	Code int    `json:"code"`           //状态码
	Msg  string `json:"msg,omitempty"`  //消息
	Data any    `json:"data,omitempty"` //结果数据
}

func NewMessage(code int, msg string, datas ...any) *Message {
	var data any
	switch len(datas) {
	case 0:
		data = nil
	case 1:
		data = datas[0]
	default:
		data = datas
	}
	return &Message{Code: code, Msg: msg, Data: data}
}

func NewSuccessMessage(datas ...any) *Message {
	return NewMessage(status.SuccessCode, "success", datas...)
}

func NewErrorMessage(err error) *Message {
	switch e := err.(type) {
	case status.ErrorCode:
		return &Message{Code: e.Code(), Msg: e.Error()}
	default:
		return &Message{Code: status.FailCode, Msg: e.Error()}
	}
}
