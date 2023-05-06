package hertz

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/stretchr/testify/assert"

	"github.com/lonelypale/gopkg/errors/status"
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
	h.GET("/test8", Bind(func(ctx *app.RequestContext) (int, int, int, error) {
		assert.NotNil(t, ctx, "*app.RequestContext cannot be nil")
		ctx.Header(consts.HeaderContentType, "text")
		return 1, 2, 3, nil
	}))

	//curl -X POST localhost:8888/test9 -d '{}'
	h.POST("/test9", Bind(func(ctx *app.RequestContext) error {
		assert.NotNil(t, ctx, "*app.RequestContext cannot be nil")
		return status.Error(errors.New("test error"))
	}))

	h.Spin()

	//更换json
	// Render
	//render.ResetJSONMarshal(json.Marshal)
	// Binding
	//binding.ResetJSONUnmarshaler(json.Unmarshal)
}

func TestCookie(t *testing.T) {
	h := server.Default()
	h.Use(Cors())

	//curl -v localhost:8888/test_cookie_set
	h.GET("/test_cookie_set", Bind(func(ctx *app.RequestContext) error {
		ctx.SetCookie("key1", "value1", 3600, "/", "127.0.0.1", protocol.CookieSameSiteNoneMode, true, true)
		ctx.SetCookie("key2", "value2", 3600, "/", "127.0.0.1:8888", protocol.CookieSameSiteNoneMode, true, true)
		ctx.SetCookie("key3", "value3", 3600, "/", "localhost", protocol.CookieSameSiteNoneMode, true, true)
		ctx.SetCookie("key4", "value4", 3600, "/", "localhost:63342", protocol.CookieSameSiteNoneMode, true, true)
		return nil
	}))

	h.POST("/test_cookie_get", Bind(func(ctx *app.RequestContext) error {
		val1 := ctx.Cookie("key1")
		fmt.Println("key1:", string(val1))
		return nil
	}))

	h.Spin()
}
