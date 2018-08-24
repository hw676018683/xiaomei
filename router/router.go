package router

import (
	"path"
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

// 获取路由的根
func (r *Router) Root() *Router {
	return &Router{
		strRoutes: r.strRoutes,
		regRoutes: r.regRoutes,
	}
}

// Group 提供带basePath的路由，代码更简洁，正则匹配更高效。
// p只能是字符串路径，不能是正则表达式。
func (r *Router) Group(p string) *Router {
	basePath := cleanPath(p)
	if r.basePath != `` {
		basePath = path.Join(r.basePath, basePath)
	}
	return &Router{
		basePath:  basePath,
		strRoutes: r.strRoutes,
		regRoutes: r.regRoutes,
	}
}

func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	return path.Clean(p)
}
