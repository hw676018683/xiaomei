package benchmark

import (
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/lovego/xiaomei"
	xiaomeirouter "github.com/lovego/xiaomei/router"
)

type route struct {
	method string
	path   string
}

type xiaomeiRouterTestCase struct {
	router *xiaomeirouter.Router
	routes []route
	hits   int
}

type httpRouterTestCase struct {
	router *httprouter.Router
	routes []route
	hits   int
}

func loadXiaomeiRouterTestCase(routes []route) *xiaomeiRouterTestCase {
	var router = xiaomeirouter.New()
	var tc = xiaomeiRouterTestCase{router: router, routes: routes}
	var paramRegexp = regexp.MustCompile(`:\w+`)
	for _, route := range routes {
		if strings.IndexByte(route.path, ':') > 0 {
			regPath := paramRegexp.ReplaceAllString(route.path, `(:\w+)`)
			router.AddX(route.method, regPath, func(*xiaomei.Request, *xiaomei.Response, []string) {
				tc.hits++
			})
		} else {
			router.Add(route.method, route.path, func(*xiaomei.Request, *xiaomei.Response) {
				tc.hits++
			})
		}
	}
	return &tc
}

func loadHttpRouterTestCase(routes []route) *httpRouterTestCase {
	var router = httprouter.New()
	var tc = httpRouterTestCase{router: router, routes: routes}
	for _, route := range routes {
		router.Handle(route.method, route.path, func(http.ResponseWriter, *http.Request, httprouter.Params) {
			tc.hits++
		})
	}
	return &tc
}

func runXiaomeiRouterTestCase(b *testing.B, tc *xiaomeiRouterTestCase) {
	b.ReportAllocs()
	tc.hits = 0

	request, err := http.NewRequest("GET", "http://localhost/", nil)
	if err != nil {
		panic(err)
	}
	req := &xiaomei.Request{Request: request}
	for i := 0; i < b.N; i++ {
		for _, route := range tc.routes {
			request.Method = route.method
			request.URL.Path = route.path
			tc.router.Handle(req, nil)
		}
	}
	if tc.hits != b.N*len(tc.routes) {
		b.Errorf("hits want: %d, got: %d\n", b.N*len(tc.routes), tc.hits)
	}
}

func runHttpRouterTestCase(b *testing.B, tc *httpRouterTestCase) {
	b.ReportAllocs()
	tc.hits = 0

	request, err := http.NewRequest("GET", "http://localhost/", nil)
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		for _, route := range tc.routes {
			request.Method = route.method
			request.URL.Path = route.path
			tc.router.ServeHTTP(nil, request)
		}
	}
	if tc.hits != b.N*len(tc.routes) {
		b.Errorf("hits want: %d, got: %d\n", b.N*len(tc.routes), tc.hits)
	}
}
