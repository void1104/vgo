package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"vgo/core"
	"vgo/log"
)

// 登录接口 - 测试路由POST方法
func login(ctx *core.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	if username == "pjx@cq.com" && password == "1104" {
		log.Info("用户登录成功")
		ctx.JSON(http.StatusOK, core.H{
			"success": "success",
		})
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

func testOutOfArrayIndex(ctx *core.Context) {
	arr := []int{1, 2, 3}
	fmt.Println(arr[4])
}
