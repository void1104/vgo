package demo

import (
	"net/http"
	"vgo/context/vgo"
)

func Demo(c *vgo.Context) {
	c.HTML(http.StatusOK, "<h1>Hello vgo!</h1>")
}
