package router

// go test -c -gcflags "-N -l"

import (
	// `fmt`
	"net/http"
	"strings"
	"testing"

	"github.com/lovego/xiaomei"
)

func notFound(req *xiaomei.Request, res *xiaomei.Response) {}

type testRouteData struct {
	method string
	path   string
	name   string
}

var testRoutes = []testRouteData{
	{`get`, `/`, `index`},
	{`get`, `/new`, `new`},
	{`post`, `/`, `create`},
	{`getx`, `/(\d+)`, `show`},
	{`getx`, `/(\d+)/edit`, `edit`},
	{`postx`, `/(\d+)/update`, `update`},
	{`postx`, `/(\d+)/destroy`, `destroy`},
}

func TestNew(t *testing.T) {
	var r *Router = New()
	if r == nil || r.strRoutes == nil || r.regRoutes == nil {
		t.Error()
	}
}

func TestRouter(t *testing.T) {
	r := New()
	matched := make(map[string]bool)
	for _, route := range testRoutes {
		testAddRoute(r, route, matched, t)
		testHandleReq(r, route, matched, ``, t)
	}
}

func TestGroupRouter(t *testing.T) {
	var prefix = `/admin`

	r := New()
	g := r.Group(prefix)
	matched := make(map[string]bool)
	for _, route := range testRoutes {
		testAddRoute(g, route, matched, t)
		testHandleReq(r, route, matched, prefix, t)
	}
}

func testAddRoute(r *Router, route testRouteData, matched map[string]bool, testing *testing.T) {
	switch route.method {
	case `get`:
		r.Get(route.path, func(req *xiaomei.Request, res *xiaomei.Response) {
			matched[route.name] = true
		})
	case `post`:
		r.Post(route.path, func(req *xiaomei.Request, res *xiaomei.Response) {
			matched[route.name] = true
		})
	case `getx`:
		r.GetX(route.path, func(req *xiaomei.Request, res *xiaomei.Response, params []string) {
			matched[route.name] = true
		})
	case `postx`:
		r.PostX(route.path, func(req *xiaomei.Request, res *xiaomei.Response, params []string) {
			matched[route.name] = true
		})
	}
}

func testHandleReq(
	r *Router, route testRouteData, matched map[string]bool, prefix string, t *testing.T,
) {
	method := strings.ToUpper(strings.TrimSuffix(route.method, `x`))
	path := strings.Replace(route.path, `(\d+)`, `123`, 1)
	req, _ := http.NewRequest(method, prefix+path, nil)
	r.Handle(&xiaomei.Request{Request: req}, nil)
	if !matched[route.name] {
		t.Errorf(`%s not matched`, prefix+path)
	}
}
