package vgo

import (
	"../../day3-router/vgo"
	"fmt"
	"log"
	"runtime"
	"strings"
)

/**
目的：在vgo中添加一个非常简单的错误处理机制，即在此类错误发生时，
	 向用户返回Internal Server Error，并且在日志中打印必要
	 的错误信息，方便进行错误定位。
错误处理也可以作为一个中间件，增强vgo框架的能力
*/

// print stack trace for debug
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

func Recovery() vgo.HandlerFunc {
	return func(c *vgo.Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
			}
		}()

		//c.Next()
	}
}
