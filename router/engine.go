package router

import "vgo/context"

// HandlerFunc 定义vgo对于请求的handler
type HandlerFunc func(ctx *context.Context)
