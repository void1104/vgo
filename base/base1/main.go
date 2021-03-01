package main

import (
	"fmt"
	"log"
	"net/http"
)

/**
net/http库的简单使用
*/

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	/**
	- 第一个参数代表地址
	- 第二个参数代表处理所有的HTTP请求的实例（因为上面已经写了两个处理器，所以这里第二个参数为nil）
	*/
	log.Fatal(http.ListenAndServe(":9999", nil))
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}
