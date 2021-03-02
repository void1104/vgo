package main

import (
	"net/http"
	"vgo/context/demo"
	"vgo/context/vgo"
)

func main() {
	r := vgo.New()
	r.GET("/", demo.Demo)

	r.GET("/hello", func(c *vgo.Context) {
		c.String(http.StatusOK, "hello %s,you are at %s \n", c.Query("name"), c.Req.URL.Path)
	})

	r.POST("/login", func(c *vgo.Context) {
		c.JSON(http.StatusOK, vgo.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":9999")
}
