package core

import (
	"fmt"
	"reflect"
	"testing"
)

/**
自动化测试： go test .
*/

func newTestRouter() *Router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

// TestParsePattern 测试解析路由
func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("单元测试 parsePattern 未通过")
	}
}

// TestGetRoute 测试获取路由
func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, ps := r.getRoute("GET", "/hello/vgo")

	if n == nil {
		t.Fatal("nil 不应该被返回")
	}

	if n.pattern != "/hello/:name" {
		t.Fatal("无法匹配 /hello/:name")
	}

	if ps["name"] != "vgo" {
		t.Fatal("动态路由变量不为 'vgo'")
	}

	fmt.Printf("路由匹配 : %s, params['name']: %s\n", n.pattern, ps["name"])
}

// TestGetRoute2 测试获取路由
func TestGetRoute2(t *testing.T) {
	r := newTestRouter()
	n1, ps1 := r.getRoute("GET", "/assets/file1.txt/")
	ok1 := n1.pattern == "/assets/*filepath" && ps1["filepath"] == "file1.txt"
	if !ok1 {
		t.Fatal("路由应该为 /assets/*filepath & filepath should be css/test.css")
	}
}

// TestRouteConflict 测试路由冲突
func TestRouteConflict(t *testing.T) {
	r := newTestRouter()
	r.addRoute("GET", "/conflict", nil)
	r.addRoute("GET", "/conflict", nil)
}
