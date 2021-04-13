package main

import (
	"fmt"
	"log"
	"net/http"
)

// Engine is the uni handler for all requests
type Engine struct{}

/**
服务器入口处理函数
	- http.ResponseWriter: 用来针对该请求的响应
	- http.Request: 该对象包含了该HTTP请求的所有的信息, 比如请求地址,Header和Body等信息
*/
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		_, _ = fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	case "/hello":
		for k, v := range req.Header {
			_, _ = fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		_, _ = fmt.Fprintf(w, "404 NOT FOUND")
	}
}

/**
在实现了Engine之后,我们拦截了所有的HTTP请求,拥有了统一的控制入口.
在这里我们可以自由定义路由映射的规则,也可以统一添加一些处理逻辑,例如日志,异常处理等.
*/
func main() {
	engine := new(Engine)
	log.Fatal(http.ListenAndServe(":9999", engine))
}
