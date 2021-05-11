package main

import (
	"time"
	"vgo/core"
	"vgo/log"
)

func Logger() core.HandlerFunc {
	return func(c *core.Context) {
		t := time.Now()
		// Process request
		c.Next()
		log.Info("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
