package main

import (
	"net/http"
	"vgo/core"
	"vgo/log"
)

func onlyForV2() core.HandlerFunc {
	return func(ctx *core.Context) {
		if ctx.Params["name"] != "pjx" {
			ctx.AuthFail()
		}
		log.Error("v2 fatal")
		ctx.Next()
	}
}

func Logger() core.HandlerFunc {
	return func(c *core.Context) {
		c.Next()
		log.Info("logger 日志打印")
	}
}

func main() {
	r := core.New()
	r.Use(Logger()) // 全局中间件
	r.GET("/", func(ctx *core.Context) {
		ctx.HTML(http.StatusOK, "<h1>Hello vgo<h1/>")
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2())
	{
		v2.GET("/hello/:name", func(ctx *core.Context) {
			ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Params["name"], ctx.Path)
		})
	}
	r.Run(":9999")
}
