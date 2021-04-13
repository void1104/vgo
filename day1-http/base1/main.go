package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	/**
		用来启动Web服务
			- addr: 代表监听的端口
			- handler: 代表处理所有的HTTP请求的实例, nil表示用标准库中的实例处理
	 */
	log.Fatal(http.ListenAndServe(":9999", nil))
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	_, _ = fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		_, _ = fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}
