package router

import (
	"path"
	"regexp"
	"strings"
)

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

// 获取路由的根
func (r *Router) Root() *Router {
	return &Router{
		strRoutes: r.strRoutes,
		regRoutes: r.regRoutes,
	}
}

func (r *Router) Get(p string, handler StrRouteHandler) *Router {
	return r.Add(`GET`, p, handler)
}

func (r *Router) Post(p string, handler StrRouteHandler) *Router {
	return r.Add(`POST`, p, handler)
}

func (r *Router) GetPost(p string, handler StrRouteHandler) *Router {
	return r.Add(`GET`, p, handler).Add(`POST`, p, handler)
}

func (r *Router) Put(p string, handler StrRouteHandler) *Router {
	return r.Add(`PUT`, p, handler)
}

func (r *Router) Options(p string, handler StrRouteHandler) *Router {
	return r.Add(`OPTIONS`, p, handler)
}

func (r *Router) Delete(p string, handler StrRouteHandler) *Router {
	return r.Add(`DELETE`, p, handler)
}

func (r *Router) GetX(reg string, handler RegRouteHandler) *Router {
	return r.AddX(`GET`, reg, handler)
}

func (r *Router) PostX(reg string, handler RegRouteHandler) *Router {
	return r.AddX(`POST`, reg, handler)
}

func (r *Router) GetPostX(reg string, handler RegRouteHandler) *Router {
	return r.AddX(`GET`, reg, handler).AddX(`POST`, reg, handler)
}

func (r *Router) PutX(reg string, handler RegRouteHandler) *Router {
	return r.AddX(`PUT`, reg, handler)
}

func (r *Router) DeleteX(reg string, handler RegRouteHandler) *Router {
	return r.AddX(`DELETE`, reg, handler)
}

func (r *Router) OptionsX(reg string, handler RegRouteHandler) *Router {
	return r.AddX(`OPTIONS`, reg, handler)
}

// 增加字符串路由
func (r *Router) Add(method string, p string, handler StrRouteHandler) *Router {
	method = strings.ToUpper(method)
	p = cleanPath(p)
	if r.basePath != `` {
		p = path.Join(r.basePath, p)
	}
	if r.strRoutes[method] == nil {
		r.strRoutes[method] = make(map[string]StrRouteHandler)
	}
	if _, ok := r.strRoutes[method][p]; ok {
		panic(`string route conflict: ` + method + ` ` + p)
	}
	r.strRoutes[method][p] = handler
	return r
}

// 增加正则路由
func (r *Router) AddX(method string, reg string, handler RegRouteHandler) *Router {
	if r.regRoutes[method] == nil {
		r.regRoutes[method] = make(map[string][]RegRoute)
	}
	regex := `^` + reg + `$`
	basePath := cleanPath(r.basePath)
	for _, regRoute := range r.regRoutes[method][basePath] {
		if regRoute.reg.String() == regex {
			panic(`regexp route conflict: ` + method + ` ` + basePath + reg)
		}
	}
	r.regRoutes[method][basePath] = append(r.regRoutes[method][basePath],
		RegRoute{regexp.MustCompile(regex), handler},
	)
	return r
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
