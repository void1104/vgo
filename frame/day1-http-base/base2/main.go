package main

import (
	"fmt"
	"log"
	"net/http"
)

/*
Engine is the uni handler for all requests
通过查看net/http的代码可以发现，Handler是一个接口，需要实现方法ServerHTTP。
也就是说，只要任何实现了ServerHTTP接口的实例，都可以作为ListenAndServer的第二个参数
*/
type Engine struct {
}

/**
在实现Engine之前，我们调用http.HandlerFunc实现了路由和Handler的映射，也就是只能针对具体的
路由写处理逻辑。比如/hello。但是在实现Engine之后，我们拦截了所有的HTTP请求，拥有了统一的控制入口。
在这里我们可以自由定义路由映射的规则，也可以统一添加一些处理逻辑，例如日志，异常处理等。
*/
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path=%q\n", req.URL.Path)
	case "/hello":
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q]=%q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND:%s\n", req.URL)
	}
}

func main() {
	engine := new(Engine)
	log.Fatal(http.ListenAndServe(":9999", engine))
}
