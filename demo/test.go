package main

import (
	"net/http"
	"vgo/core"
)

// demo 路由请求方式
//func main() {
//	r := vgo.New()
//
//	r.GET("/", func(ctx *vgo.Context) {
//		ctx.HTML(http.StatusOK, "<h1>Hello Vgo</h1>")
//	})
//
//	r.GET("/hello", func(ctx *vgo.Context) {
//		ctx.String(http.StatusOK, "hello %s, you're at %s\n")
//	})
//
//	r.POST("/login", func(ctx *vgo.Context) {
//		ctx.JSON(http.StatusOK, vgo.H{
//			"username": ctx.PostForm("username"),
//			"password": ctx.PostForm("password"),
//		})
//	})
//	r.Run(":9999")
//}

// demo 路由冲突
func main() {
	r := core.New()
	r.GET("/hello", func(ctx *core.Context) {
		ctx.String(http.StatusOK, "test router conflict")
	})

	r.GET("/hello", func(ctx *core.Context) {
		ctx.String(http.StatusOK, "test router conflict")
	})

	r.Run(":9999")
}
