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

func main() {
	r := core.New()
	r.GET("/log/list", logList)

	r.Run(":9999")
}

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
	for i := 0; i < 10; i++ {
		log.Info("test log write local")
	}
	ctx.JSON(http.StatusOK, list)
}
