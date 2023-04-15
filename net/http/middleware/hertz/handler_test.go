package hertz

import (
	"context"
	"errors"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/stretchr/testify/assert"
)

type testUser struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var tuser = testUser{
	Name: "test user",
	Age:  99,
}

func TestBind(t *testing.T) {
	h := server.Default()

	//curl -v localhost:8888/test1
	h.GET("/test1", Bind(func() {
	}))

	//curl localhost:8888/test2
	h.GET("/test2", Bind(func(c context.Context) {
		assert.NotNil(t, c, "context.Context cannot be nil")
	}))

	//curl localhost:8888/test3
	h.GET("/test3", Bind(func(ctx *app.RequestContext) {
		assert.NotNil(t, ctx, "*app.RequestContext cannot be nil")
	}))

	//curl localhost:8888/test4
	h.GET("/test4", Bind(func(c context.Context, ctx *app.RequestContext) {
		assert.NotNil(t, c, "context.Context cannot be nil")
		assert.NotNil(t, ctx, "*app.RequestContext cannot be nil")
	}))

	//curl localhost:8888/test5
	h.GET("/test5", Bind(func(ctx *app.RequestContext, c context.Context) {
		assert.NotNil(t, c, "context.Context cannot be nil")
		assert.NotNil(t, ctx, "*app.RequestContext cannot be nil")
	}))

	//curl -X POST -H "Content-Type: application/json" localhost:8888/test6 -d '{"name":"test user","age":99}'
	h.POST("/test6", Bind(func(user testUser, ctx *app.RequestContext) {
		assert.NotNil(t, ctx, "*app.RequestContext cannot be nil")
		assert.NotNil(t, user, "testUser cannot be nil")
		assert.Equal(t, user, tuser, "they should be equal")
	}))

	//curl -X POST -H "Content-Type: application/json" localhost:8888/test7 -d '{"name":"test user","age":99}'
	h.POST("/test7", Bind(func(ctx *app.RequestContext, user *testUser) {
		assert.NotNil(t, ctx, "*app.RequestContext cannot be nil")
		assert.NotNil(t, user, "u *testUser cannot be nil")
		assert.Equal(t, *user, tuser, "they should be equal")
	}))

	//curl localhost:8888/test8
	h.GET("/test8", Bind(func(ctx *app.RequestContext) error {
		assert.NotNil(t, ctx, "*app.RequestContext cannot be nil")
		ctx.Header(consts.HeaderContentType, "text")
		return nil
	}))

	//curl localhost:8888/test9
	h.GET("/test9", Bind(func(ctx *app.RequestContext) error {
		assert.NotNil(t, ctx, "*app.RequestContext cannot be nil")
		return errors.New("test error")
	}))

	h.Spin()

	//更换json
	// Render
	//render.ResetJSONMarshal(json.Marshal)
	// Binding
	//binding.ResetJSONUnmarshaler(json.Unmarshal)
}
