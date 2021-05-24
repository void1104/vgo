package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"vgo/core"
	"vgo/log"
)

func Logger() core.HandlerFunc {
	return func(c *core.Context) {
		c.Next()
		log.Info("test log write local")
	}
}

func Cors() core.HandlerFunc {
	return func(c *core.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Request-Method", "OPTION,GET,POST")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Next()
	}
}

const TestLogPath = "D:\\log.txt" // 测试日志路径

func main() {
	// 1.构建框架环境
	r := core.New()
	//log.SetLogPath(TestLogPath) // 自定义设置日志输出路径
	// 2. 注册中间件
	r.Use(Logger())

	// 3. 注册路由，跨域报错 -> 无法访问。
	r.GET("/log/list", logList)
	r.POST("/login", login)
	r.GET("/ping", func(c *core.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 4. 注册分组路由
	gr := r.Group("/cors")
	gr.Use(Cors()) // cors分组下使用cors中间件设置跨域
	{
		gr.GET("/log/list", logList)
		gr.POST("/login", login)

		// 5. 动态路由 - 参数匹配
		gr.GET("/hello/:name/space", func(ctx *core.Context) {
			ctx.String(http.StatusOK, "The dynamic routing passes in parameters: %s", ctx.Params["name"])
		})

		// 6. 动态路由 - 模糊匹配
		gr.GET("/static/*filepath", func(ctx *core.Context) {
			ctx.String(http.StatusOK, "The dynamic routing passes in parameters: /%s", ctx.Params["filepath"])
		})
	}

	_ = r.Run(":9999")
}
