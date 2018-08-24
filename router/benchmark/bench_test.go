package benchmark

import (
	"testing"
)

var xiaomeiRouterStatic = loadXiaomeiRouterTestCase(staticRoutes)
var httpRouterStatic = loadHttpRouterTestCase(staticRoutes)

var xiaomeiRouterGithub = loadXiaomeiRouterTestCase(githubAPI)
var httpRouterGithub = loadHttpRouterTestCase(githubAPI)

var xiaomeiRouterGooglePlus = loadXiaomeiRouterTestCase(googlePlusAPI)
var httpRouterGooglePlus = loadHttpRouterTestCase(googlePlusAPI)

var xiaomeiRouterParseCom = loadXiaomeiRouterTestCase(parseComAPI)
var httpRouterParseCom = loadHttpRouterTestCase(parseComAPI)

func BenchmarkXiaomeiRouter_Static157(b *testing.B) {
	runXiaomeiRouterTestCase(b, xiaomeiRouterStatic)
}

func BenchmarkHttpRouter_Static157(b *testing.B) {
	runHttpRouterTestCase(b, httpRouterStatic)
}

func BenchmarkXiaomeiRouter_Github203(b *testing.B) {
	runXiaomeiRouterTestCase(b, xiaomeiRouterGithub)
}

func BenchmarkHttpRouter_Github203(b *testing.B) {
	runHttpRouterTestCase(b, httpRouterGithub)
}

func BenchmarkXiaomeiRouter_GooglePlus13(b *testing.B) {
	runXiaomeiRouterTestCase(b, xiaomeiRouterGooglePlus)
}

func BenchmarkHttpRouter_GooglePlus13(b *testing.B) {
	runHttpRouterTestCase(b, httpRouterGooglePlus)
}

func BenchmarkXiaomeiRouter_ParseCom26(b *testing.B) {
	runXiaomeiRouterTestCase(b, xiaomeiRouterParseCom)
}

func BenchmarkHttpRouter_ParseCom26(b *testing.B) {
	runHttpRouterTestCase(b, httpRouterParseCom)
}
