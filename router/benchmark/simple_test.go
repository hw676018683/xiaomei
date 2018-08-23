package benchmark

import (
	"net/http"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/lovego/xiaomei"
	xiaomeiRouter "github.com/lovego/xiaomei/router"
)

func BenchmarkXiaomeiRouter_Simple(b *testing.B) {
	b.ReportAllocs()

	var hits int
	var router = xiaomeiRouter.New()
	router.Get("/users", func(*xiaomei.Request, *xiaomei.Response) {
		hits++
	})
	request, err := http.NewRequest("GET", "http://localhost/users", nil)
	if err != nil {
		panic(err)
	}

	req := &xiaomei.Request{Request: request}
	for i := 0; i < b.N; i++ {
		router.Handle(req, nil)
	}
	if hits != b.N {
		b.Errorf("hits want: %d, got: %d\n", b.N, hits)
	}
}

func BenchmarkHttpRouter_Simple(b *testing.B) {
	b.ReportAllocs()

	var hits int
	var router = httprouter.New()
	router.GET("/users", func(http.ResponseWriter, *http.Request, httprouter.Params) {
		hits++
	})
	request, err := http.NewRequest("GET", "http://localhost/users", nil)
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(nil, request)
	}
	if hits != b.N {
		b.Errorf("hits want: %d, got: %d\n", b.N, hits)
	}
}
