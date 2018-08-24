package router

import (
	"regexp"

	"github.com/lovego/xiaomei"
)

// 字符串路由处理函数
type StrRouteHandler func(*xiaomei.Request, *xiaomei.Response)

// 正则路由处理函数，第三个参数是正则匹配的结果
type RegRouteHandler func(*xiaomei.Request, *xiaomei.Response, []string)

type RegRoute struct {
	reg     *regexp.Regexp
	handler RegRouteHandler
}

type Router struct {
	// 基础路径
	basePath string
	// 字符串路由 method     path
	strRoutes map[string]map[string]StrRouteHandler
	// 正则路由   method     base_path
	regRoutes map[string]map[string][]RegRoute
}

func New() *Router {
	return &Router{
		strRoutes: make(map[string]map[string]StrRouteHandler),
		regRoutes: make(map[string]map[string][]RegRoute),
	}
}
