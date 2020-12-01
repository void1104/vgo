package base3

import (
	"fmt"
	"net/http"
)

/**
vgo框架的设计以及API均参考了gin
*/

func main() {
	r := New()
	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})

	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q]=%q\n", k, v)
		}
	})

	r.Run(":9999")
}
