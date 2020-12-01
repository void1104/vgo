package main

import (
	"./vgo"
	"net/http"
)

func main() {
	r := vgo.New()
	r.GET("/", func(ctx *vgo.Context) {
		ctx.HTML(http.StatusOK, "<h1>Hello dnw</h1>")
	})

	r.GET("/hello", func(ctx *vgo.Context) {
		ctx.String(http.StatusOK, "hello %s, you're at %s\n", ctx.Query("name"), ctx.Path)
	})

	r.GET("/hello/:name", func(ctx *vgo.Context) {
		ctx.String(http.StatusOK, "hello %s, you're at %s\"", ctx.Param("name"), ctx.Path)
	})

	r.GET("/assets/*filepath", func(ctx *vgo.Context) {
		ctx.JSON(http.StatusOK, vgo.H{"filepath": ctx.Param("filepath")})
	})

	r.Run(":9999")
}
