package main

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"strings"
	"vgo/core"
	"vgo/log"
)

func Logger() core.HandlerFunc {
	return func(c *core.Context) {
		c.Next()
		log.Info("test log write local")
	}
}

func AuthCheck() core.HandlerFunc {
	return func(ctx *core.Context) {
		if ctx.PostForm("username") != "pjx@cq.com" {
			log.Info("访问接口失败，非权限人员访问")
			ctx.AuthFail()
		}
		log.Info("鉴权成功")
		ctx.Next()
	}
}

func main() {
	r := core.New()
	r.Use(Logger())
	r.GET("/log/list", logList)
	r.POST("/login", login)

	gr := r.Group("/admin")
	gr.Use(AuthCheck())
	{
		gr.GET("/check", nil)
	}

	r.Run(":9999")
}

// 鉴权接口 - 测试路由分组
func admin(ctx *core.Context) {

}

// 登录接口 - 测试路由POST方法
func login(ctx *core.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	if username == "pjx@cq.com" && password == "1104" {
		log.Info("用户登录成功")
		ctx.JSON(500, "success")
	} else {
		log.Info("用户登录失败，账号或密码错误")
		ctx.Fail()
	}
}

// 日志列表 - 测试路由GET方法
func logList(ctx *core.Context) {
	file, err := os.Open("./log.txt")
	if err != nil {
		ctx.Fail()
		return
	}
	var list []string
	buf := bufio.NewReader(file)
	for {
		line, inErr := buf.ReadString('\n')
		if inErr != nil {
			if inErr == io.EOF {
				break
			}
			ctx.Fail()
			return
		}
		line = strings.TrimSpace(line)
		list = append(list, line)
	}
	ctx.JSON(http.StatusOK, list)
}
