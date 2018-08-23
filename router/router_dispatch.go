package router

import (
	"strings"

	"github.com/lovego/xiaomei"
)

// 处理请求
func (r *Router) Handle(req *xiaomei.Request, res *xiaomei.Response) bool {
	method := strings.ToUpper(req.Method)
	if method == `HEAD` {
		method = `GET`
	}
	path := cleanPath(req.URL.Path)
	if r.strRoutesMatch(method, path, req, res) || r.regRoutesMatch(method, path, req, res) {
		return true
	}
	return false
}

func (r *Router) strRoutesMatch(method string, path string, req *xiaomei.Request, res *xiaomei.Response) bool {
	routes := r.strRoutes[method]
	if routes == nil {
		return false
	}
	handler := routes[path]
	if handler == nil {
		return false
	}
	handler(req, res)
	return true
}

// 按斜线作为分隔符从深到浅依次匹配
func (r *Router) regRoutesMatch(method string, path string, req *xiaomei.Request, res *xiaomei.Response) bool {
	routes := r.regRoutes[method]
	if routes == nil {
		return false
	}
	prefix := path
	for {
		if slice := routes[prefix]; slice != nil {
			var pathToMatch = path
			if prefix != `/` {
				pathToMatch = path[len(prefix):]
			}
			for _, route := range slice {
				if m := route.reg.FindStringSubmatch(pathToMatch); m != nil {
					route.handler(req, res, m)
					return true
				}
			}
		}
		// 上一层路径
		if i := strings.LastIndexByte(prefix, '/'); i > 0 {
			prefix = prefix[:i]
		} else if i == 0 && prefix != `/` {
			prefix = `/`
		} else {
			return false
		}
	}
}
