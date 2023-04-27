package hertz

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	log "github.com/sirupsen/logrus"

	"github.com/lonelypale/gopkg/errors/status"
	"github.com/lonelypale/gopkg/net/http"
	"github.com/lonelypale/gopkg/ref"
)

// Bind 转换成 app.HandlerFunc 形式的 Web 处理接口
// fn 必须是函数
// 入参(0-n): context.Context、*app.RequestContext、struct{}、*struct{}、map、array
// struct-tag: json、form、uri、query、header、any
// 出参(0-n): any 无要求
// 注意: hertz 绑定 json 请求的 body 时, 内容必须是 content-type -> application/json
func Bind(fn any) app.HandlerFunc {
	if fnType := reflect.TypeOf(fn); ref.IsFuncType(fnType) {
		bindType := make([]reflect.Type, fnType.NumIn())
		for n := 0; n < fnType.NumIn(); n++ {
			bindType[n] = fnType.In(n)
		}

		binder := &bindHandler{
			fn:       fn,
			fnType:   fnType,
			fnValue:  reflect.ValueOf(fn),
			bindType: bindType,
		}
		return binder.Invoke
	}

	panic(errors.New("fn should be func()"))
}

// appRequestContextType *app.RequestContext 的反射类型
var appRequestContextType = reflect.TypeOf((*app.RequestContext)(nil))

// bindHandler Bind 形式的 Web 处理接口
type bindHandler struct {
	fn       interface{}
	fnType   reflect.Type
	fnValue  reflect.Value
	bindType []reflect.Type
}

func (b *bindHandler) Invoke(c context.Context, ctx *app.RequestContext) {
	WebInvoke(c, ctx, b.call)
}

func (b *bindHandler) call(c context.Context, ctx *app.RequestContext) ([]any, error) {
	bindNum := len(b.bindType)
	in := make([]reflect.Value, bindNum)

	// 反射创建需要绑定的请求参数
	for i := 0; i < bindNum; i++ {
		typ := b.bindType[i]
		switch {
		case ref.IsContextType(typ):
			in[i] = reflect.ValueOf(c)
		case typ == appRequestContextType:
			in[i] = reflect.ValueOf(ctx)
		default:
			isPtr := typ.Kind() == reflect.Pointer
			var val reflect.Value
			if isPtr {
				val = reflect.New(typ.Elem())
			} else {
				val = reflect.New(typ)
			}

			if err := ctx.BindAndValidate(val.Interface()); err != nil {
				return nil, status.Error(err)
			}

			if isPtr {
				in[i] = val
			} else {
				in[i] = val.Elem()
			}
		}
	}

	// 执行处理函数，并返回结果
	out := b.fnValue.Call(in)
	result := make([]any, len(out))
	for i, o := range out {
		result[i] = o.Interface()
	}

	return result, nil
}

// WebInvoke 可自定义的 web 执行函数
var WebInvoke = defaultWebInvoke

// defaultWebInvoke 默认的 web 执行函数
func defaultWebInvoke(c context.Context, ctx *app.RequestContext, fn func(c context.Context, ctx *app.RequestContext) ([]any, error)) {
	var result *http.Message
	var err error
	defer func() {
		if r := recover(); r != nil {
			switch e := r.(type) {
			case error:
				err = e
			default:
				err = fmt.Errorf("%v", e)
			}
		}

		if err != nil {
			result = http.NewMessage(err)
			request := "body is stream"
			if !ctx.Request.IsBodyStream() {
				request = string(ctx.Request.Body())
			}
			if request != "" {
				log.Errorf("%s\n\nrequest: %s", err, request)
			} else {
				log.Error(err)
			}
		}

		ctx.JSON(consts.StatusOK, result)
	}()

	out, err := fn(c, ctx)
	if err != nil {
		return
	}

	switch len(out) {
	case 0:
		ctx.Response.Header.SetNoDefaultContentType(true)
		contentType := ctx.Response.Header.Get(consts.HeaderContentType)
		if contentType == "" {
			result = http.NewMessage()
		}
	case 1:
		switch v := out[0].(type) {
		case http.Message:
			result = &v
		case *http.Message:
			result = v
		case error:
			err = v
			return
		default:
			result = http.NewMessage(v)
		}
	default:
		// 当返回值为多个时，最后一个固定是error
		lastIndex := len(out) - 1
		last := out[lastIndex]
		if v, ok := last.(error); ok {
			err = v
		} else {
			result = http.NewMessage(out[:lastIndex]...)
		}
	}
}
