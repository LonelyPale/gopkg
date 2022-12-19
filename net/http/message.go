package http

import "github.com/lonelypale/gopkg/errors/status"

type Message struct {
	Code int    `json:"code"`           //状态码
	Msg  string `json:"msg,omitempty"`  //消息
	Data any    `json:"data,omitempty"` //结果数据
}

func NewMessage(code int, msg string, datas ...any) *Message {
	var data any
	if len(datas) == 1 {
		data = datas[0]
	} else if len(datas) > 1 {
		data = datas
	}
	return &Message{Code: code, Msg: msg, Data: data}
}

func NewSuccessMessage(datas ...any) *Message {
	return NewMessage(status.Success, "success", datas...)
}

func NewErrorMessage(err error) *Message {
	switch e := err.(type) {
	case status.ErrorCode:
		return &Message{Code: e.Code(), Msg: e.Error()}
	default:
		return &Message{Code: status.Fail, Msg: e.Error()}
	}
}
