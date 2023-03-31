package httpx

import (
	"context"
	"github.com/smartystreets/goconvey/convey"
	"net/url"
	"testing"
	"time"
)

func TestHttpContext(t *testing.T) {
	ctx := context.Background()
	values := url.Values{
		"name": {"donyhuang"},
	}
	httpUrl := "http://localhost:8080/test"
	convey.Convey("http", t, func() {
		convey.Convey("get", func() {
			convey.Convey("get map", func() {
				ctxWithTimeout, _ := context.WithTimeout(ctx, time.Second)

				body, err := HttpContextGet(ctxWithTimeout, httpUrl, map[string]string{
					"name": "donyhuang",
				})
				convey.So(err, convey.ShouldBeNil)
				t.Log(string(body))
			})
			convey.Convey("get values", func() {
				ctxWithTimeout, _ := context.WithTimeout(ctx, time.Second)

				body, err := HttpContextGet(ctxWithTimeout, httpUrl, url.Values{
					"name": {"donyhuang"},
				})
				convey.So(err, convey.ShouldBeNil)
				t.Log(string(body))
			})
			convey.Convey("get str", func() {
				ctxWithTimeout, _ := context.WithTimeout(ctx, time.Second)

				body, err := HttpContextGet(ctxWithTimeout, httpUrl, "name=donyhuang")
				convey.So(err, convey.ShouldBeNil)
				t.Log(string(body))
			})
			convey.Convey("get nil", func() {
				ctxWithTimeout, _ := context.WithTimeout(ctx, time.Second)
				body, err := HttpContextGet(ctxWithTimeout, httpUrl, "")
				convey.So(err, convey.ShouldBeNil)
				t.Log(string(body))
			})
		})
		convey.Convey("post", func() {
			convey.Convey("post urls", func() {
				ctxWithTimeout, _ := context.WithTimeout(ctx, time.Second)

				body, err := HttpContextPost(ctxWithTimeout, httpUrl, values)
				convey.So(err, convey.ShouldBeNil)
				t.Log(string(body))
			})
			convey.Convey("post map", func() {
				ctxWithTimeout, _ := context.WithTimeout(ctx, time.Second)

				body, err := HttpContextPost(ctxWithTimeout, httpUrl, map[string]string{
					"name": "hyx",
				})
				convey.So(err, convey.ShouldBeNil)
				t.Log(string(body))
			})
			convey.Convey("post string", func() {
				ctxWithTimeout, _ := context.WithTimeout(ctx, time.Second)
				body, err := HttpContextPost(ctxWithTimeout, httpUrl, "name=pzfyp")
				convey.So(err, convey.ShouldBeNil)
				t.Log(string(body))
			})

		})
		convey.Convey("json", func() {
			convey.Convey("json bytes", func() {
				ctxWithTimeout, _ := context.WithTimeout(ctx, time.Second)
				body, err := HttpContextJson(ctxWithTimeout, httpUrl, []byte(`{"name":"donyhuang"}`))
				convey.So(err, convey.ShouldBeNil)
				t.Log(string(body))
			})
			convey.Convey("json string", func() {
				ctxWithTimeout, _ := context.WithTimeout(ctx, time.Second)
				body, err := HttpContextJson(ctxWithTimeout, httpUrl, `{"name":"donyhuang"}`)
				convey.So(err, convey.ShouldBeNil)
				t.Log(string(body))
			})
			convey.Convey("json map", func() {
				ctxWithTimeout, _ := context.WithTimeout(ctx, time.Second)
				body, err := HttpContextJson(ctxWithTimeout, httpUrl, map[string]interface{}{
					"name": "donyhuang",
					"age":  12,
				})
				convey.So(err, convey.ShouldBeNil)
				t.Log(string(body))
			})
			convey.Convey("json struct", func() {
				var sj = struct {
					Name string
					Age  int
				}{Name: "donyhuang", Age: 21}
				ctxWithTimeout, _ := context.WithTimeout(ctx, time.Second)
				body, err := HttpContextJson(ctxWithTimeout, httpUrl, sj)
				convey.So(err, convey.ShouldBeNil)
				t.Log(string(body))
				ctxWithTimeout, _ = context.WithTimeout(ctx, time.Second)
				body, err = HttpContextJson(ctxWithTimeout, httpUrl, &sj)
				convey.So(err, convey.ShouldBeNil)
				t.Log(string(body))
			})
			convey.Convey("json slice", func() {
				var sj = []string{
					"a", "b", "c",
				}
				ctxWithTimeout, _ := context.WithTimeout(ctx, time.Second)
				body, err := HttpContextJson(ctxWithTimeout, httpUrl, sj)
				convey.So(err, convey.ShouldBeNil)
				t.Log(string(body))
				ctxWithTimeout, _ = context.WithTimeout(ctx, time.Second)
				body, err = HttpContextJson(ctxWithTimeout, httpUrl, &sj)
				convey.So(err, convey.ShouldBeNil)
				t.Log(string(body))
			})
		})
	})
}
