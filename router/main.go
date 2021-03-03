package main

import (
	"net/http"
	"vgo/router/vgo"
)

func main() {
	r := vgo.New()
	r.GET("/", func(c *vgo.Context) {
		c.HTML(http.StatusOK, "<h1>Hello v-go</h1>")
	})

	r.GET("/hello/h1/h2", func(c *vgo.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Req.URL.Path)
	})

	r.GET("/hello/h1/h2", func(c *vgo.Context) {
		c.JSON(http.StatusOK, vgo.H{"filepath": c.Param("lang")})
	})



	r.Run(":9999")
}
