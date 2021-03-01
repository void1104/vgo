package main

import (
	"fmt"
	"log"
	"net/http"
)

// Engine 实现ServeHTTP()方法，作为处理所有的HTTP请求的实例
type Engine struct{}

/**
实现框架的第一步，将所有的请求转向我们自己的处理逻辑
	- 在实现了Engine之后，我们拦截了所有的HTTP请求，拥有了统一的入口
	- 在这里我们可以自由定义路由映射的规则，也可以统一添加一些处理逻辑，例如日志，异常处理等。
*/
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	case "/hello":
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

func main() {
	engine := new(Engine)
	log.Fatal(http.ListenAndServe(":9999", engine))
}
