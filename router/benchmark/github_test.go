package benchmark

import (
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/lovego/xiaomei"
	xiaomeiRouter "github.com/lovego/xiaomei/router"
)

var xiaomeiGithubRouter = xiaomeiRouter.New()
var httpGithubRouter = httprouter.New()
var xiaomeiGithubHits, httpGithubHits int

func init() {
	var paramRegexp = regexp.MustCompile(`:\w+`)
	for _, path := range staticRoutes {
		if strings.IndexByte(path, ':') > 0 {
			regPath := paramRegexp.ReplaceAllString(path, `(:\w+)`)
			xiaomeiGithubRouter.GetX(regPath, func(*xiaomei.Request, *xiaomei.Response, []string) {
				xiaomeiGithubHits++
			})
		} else {
			xiaomeiGithubRouter.Get(path, func(*xiaomei.Request, *xiaomei.Response) {
				xiaomeiGithubHits++
			})
		}
		httpGithubRouter.GET(path, func(http.ResponseWriter, *http.Request, httprouter.Params) {
			httpGithubHits++
		})
	}
}

func BenchmarkXiaomeiRouter_Github(b *testing.B) {
	b.ReportAllocs()

	xiaomeiGithubHits = 0
	request, err := http.NewRequest("GET", "http://localhost/users", nil)
	if err != nil {
		panic(err)
	}
	req := &xiaomei.Request{Request: request}
	for i := 0; i < b.N; i++ {
		for _, path := range staticRoutes {
			req.Request.URL.Path = path
			xiaomeiGithubRouter.Handle(req, nil)
		}
	}
	if xiaomeiGithubHits != b.N*len(staticRoutes) {
		b.Errorf("hits want: %d, got: %d\n", b.N*len(staticRoutes), xiaomeiGithubHits)
	}
}

func BenchmarkHttpRouter_Github(b *testing.B) {
	b.ReportAllocs()

	httpGithubHits = 0
	request, err := http.NewRequest("GET", "http://localhost/users", nil)
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		for _, path := range staticRoutes {
			request.URL.Path = path
			httpGithubRouter.ServeHTTP(nil, request)
		}
	}
	if httpGithubHits != b.N*len(staticRoutes) {
		b.Errorf("hits want: %d, got: %d\n", b.N*len(staticRoutes), httpGithubHits)
	}
}

var githubRoutes = []string{
	// auth
	"/authorizations",
	"/authorizations/:id",
	"/authorizations",
	"/authorizations/clients/:client_id",
	"/authorizations/:id",
	"/authorizations/:id",
	"/applications/:client_id/tokens/:access_token",
	"/applications/:client_id/tokens",
	"/applications/:client_id/tokens/:access_token",

	// activity
	"/events",
	"/repos/:owner/:repo/events",
	"/networks/:owner/:repo/events",
	"/orgs/:org/events",
	"/users/:user/received_events",
	"/users/:user/received_events/public",
	"/users/:user/events",
	"/users/:user/events/public",
	"/users/:user/events/orgs/:org",
	"/feeds",
	"/notifications",
	"/repos/:owner/:repo/notifications",
	"/notifications",
	"/repos/:owner/:repo/notifications",
	"/notifications/threads/:id",
	"/notifications/threads/:id",
	"/notifications/threads/:id/subscription",
	"/notifications/threads/:id/subscription",
	"/notifications/threads/:id/subscription",
	"/repos/:owner/:repo/stargazers",
	"/users/:user/starred",
	"/user/starred",
	"/user/starred/:owner/:repo",
	"/user/starred/:owner/:repo",
	"/user/starred/:owner/:repo",
	"/repos/:owner/:repo/subscribers",
	"/users/:user/subscriptions",
	"/user/subscriptions",
	"/repos/:owner/:repo/subscription",
	"/repos/:owner/:repo/subscription",
	"/repos/:owner/:repo/subscription",
	"/user/subscriptions/:owner/:repo",
	"/user/subscriptions/:owner/:repo",
	"/user/subscriptions/:owner/:repo",

	// Gists
	"/users/:user/gists",
	"/gists",
	"/gists/public",
	"/gists/starred",
	"/gists/:id",
	"/gists",
	"/gists/:id",
	"/gists/:id/star",
	"/gists/:id/star",
	"/gists/:id/star",
	"/gists/:id/forks",
	"/gists/:id",

	// Git Data
	"/repos/:owner/:repo/git/blobs/:sha",
	"/repos/:owner/:repo/git/blobs",
	"/repos/:owner/:repo/git/commits/:sha",
	"/repos/:owner/:repo/git/commits",
	"/repos/:owner/:repo/git/refs/*ref",
	"/repos/:owner/:repo/git/refs",
	"/repos/:owner/:repo/git/refs",
	"/repos/:owner/:repo/git/refs/*ref",
	"/repos/:owner/:repo/git/refs/*ref",
	"/repos/:owner/:repo/git/tags/:sha",
	"/repos/:owner/:repo/git/tags",
	"/repos/:owner/:repo/git/trees/:sha",
	"/repos/:owner/:repo/git/trees",

	// Issues
	"/issues",
	"/user/issues",
	"/orgs/:org/issues",
	"/repos/:owner/:repo/issues",
	"/repos/:owner/:repo/issues/:number",
	"/repos/:owner/:repo/issues",
	"/repos/:owner/:repo/issues/:number",
	"/repos/:owner/:repo/assignees",
	"/repos/:owner/:repo/assignees/:assignee",
	"/repos/:owner/:repo/issues/:number/comments",
	"/repos/:owner/:repo/issues/comments",
	"/repos/:owner/:repo/issues/comments/:id",
	"/repos/:owner/:repo/issues/:number/comments",
	"/repos/:owner/:repo/issues/comments/:id",
	"/repos/:owner/:repo/issues/comments/:id",
	"/repos/:owner/:repo/issues/:number/events",
	"/repos/:owner/:repo/issues/events",
	"/repos/:owner/:repo/issues/events/:id",
	"/repos/:owner/:repo/labels",
	"/repos/:owner/:repo/labels/:name",
	"/repos/:owner/:repo/labels",
	"/repos/:owner/:repo/labels/:name",
	"/repos/:owner/:repo/labels/:name",
	"/repos/:owner/:repo/issues/:number/labels",
	"/repos/:owner/:repo/issues/:number/labels",
	"/repos/:owner/:repo/issues/:number/labels/:name",
	"/repos/:owner/:repo/issues/:number/labels",
	"/repos/:owner/:repo/issues/:number/labels",
	"/repos/:owner/:repo/milestones/:number/labels",
	"/repos/:owner/:repo/milestones",
	"/repos/:owner/:repo/milestones/:number",
	"/repos/:owner/:repo/milestones",
	"/repos/:owner/:repo/milestones/:number",
	"/repos/:owner/:repo/milestones/:number",

	// Miscellaneous
	"/emojis",
	"/gitignore/templates",
	"/gitignore/templates/:name",
	"/markdown",
	"/markdown/raw",
	"/meta",
	"/rate_limit",

	// Organizations
	"/users/:user/orgs",
	"/user/orgs",
	"/orgs/:org",
	"/orgs/:org",
	"/orgs/:org/members",
	"/orgs/:org/members/:user",
	"/orgs/:org/members/:user",
	"/orgs/:org/public_members",
	"/orgs/:org/public_members/:user",
	"/orgs/:org/public_members/:user",
	"/orgs/:org/public_members/:user",
	"/orgs/:org/teams",
	"/teams/:id",
	"/orgs/:org/teams",
	"/teams/:id",
	"/teams/:id",
	"/teams/:id/members",
	"/teams/:id/members/:user",
	"/teams/:id/members/:user",
	"/teams/:id/members/:user",
	"/teams/:id/repos",
	"/teams/:id/repos/:owner/:repo",
	"/teams/:id/repos/:owner/:repo",
	"/teams/:id/repos/:owner/:repo",
	"/user/teams",

	// Pull Requests
	"/repos/:owner/:repo/pulls",
	"/repos/:owner/:repo/pulls/:number",
	"/repos/:owner/:repo/pulls",
	"/repos/:owner/:repo/pulls/:number",
	"/repos/:owner/:repo/pulls/:number/commits",
	"/repos/:owner/:repo/pulls/:number/files",
	"/repos/:owner/:repo/pulls/:number/merge",
	"/repos/:owner/:repo/pulls/:number/merge",
	"/repos/:owner/:repo/pulls/:number/comments",
	"/repos/:owner/:repo/pulls/comments",
	"/repos/:owner/:repo/pulls/comments/:number",
	"/repos/:owner/:repo/pulls/:number/comments",
	"/repos/:owner/:repo/pulls/comments/:number",
	"/repos/:owner/:repo/pulls/comments/:number",

	// Repositories
	"/user/repos",
	"/users/:user/repos",
	"/orgs/:org/repos",
	"/repositories",
	"/user/repos",
	"/orgs/:org/repos",
	"/repos/:owner/:repo",
	"/repos/:owner/:repo",
	"/repos/:owner/:repo/contributors",
	"/repos/:owner/:repo/languages",
	"/repos/:owner/:repo/teams",
	"/repos/:owner/:repo/tags",
	"/repos/:owner/:repo/branches",
	"/repos/:owner/:repo/branches/:branch",
	"/repos/:owner/:repo",
	"/repos/:owner/:repo/collaborators",
	"/repos/:owner/:repo/collaborators/:user",
	"/repos/:owner/:repo/collaborators/:user",
	"/repos/:owner/:repo/collaborators/:user",
	"/repos/:owner/:repo/comments",
	"/repos/:owner/:repo/commits/:sha/comments",
	"/repos/:owner/:repo/commits/:sha/comments",
	"/repos/:owner/:repo/comments/:id",
	"/repos/:owner/:repo/comments/:id",
	"/repos/:owner/:repo/comments/:id",
	"/repos/:owner/:repo/commits",
	"/repos/:owner/:repo/commits/:sha",
	"/repos/:owner/:repo/readme",
	"/repos/:owner/:repo/contents/*path",
	"/repos/:owner/:repo/contents/*path",
	"/repos/:owner/:repo/contents/*path",
	"/repos/:owner/:repo/:archive_format/:ref",
	"/repos/:owner/:repo/keys",
	"/repos/:owner/:repo/keys/:id",
	"/repos/:owner/:repo/keys",
	"/repos/:owner/:repo/keys/:id",
	"/repos/:owner/:repo/keys/:id",
	"/repos/:owner/:repo/downloads",
	"/repos/:owner/:repo/downloads/:id",
	"/repos/:owner/:repo/downloads/:id",
	"/repos/:owner/:repo/forks",
	"/repos/:owner/:repo/forks",
	"/repos/:owner/:repo/hooks",
	"/repos/:owner/:repo/hooks/:id",
	"/repos/:owner/:repo/hooks",
	"/repos/:owner/:repo/hooks/:id",
	"/repos/:owner/:repo/hooks/:id/tests",
	"/repos/:owner/:repo/hooks/:id",
	"/repos/:owner/:repo/merges",
	"/repos/:owner/:repo/releases",
	"/repos/:owner/:repo/releases/:id",
	"/repos/:owner/:repo/releases",
	"/repos/:owner/:repo/releases/:id",
	"/repos/:owner/:repo/releases/:id",
	"/repos/:owner/:repo/releases/:id/assets",
	"/repos/:owner/:repo/stats/contributors",
	"/repos/:owner/:repo/stats/commit_activity",
	"/repos/:owner/:repo/stats/code_frequency",
	"/repos/:owner/:repo/stats/participation",
	"/repos/:owner/:repo/stats/punch_card",
	"/repos/:owner/:repo/statuses/:ref",
	"/repos/:owner/:repo/statuses/:ref",

	// Search
	"/search/repositories",
	"/search/code",
	"/search/issues",
	"/search/users",
	"/legacy/issues/search/:owner/:repository/:state/:keyword",
	"/legacy/repos/search/:keyword",
	"/legacy/user/search/:keyword",
	"/legacy/user/email/:email",

	// Users
	"/users/:user",
	"/user",
	"/user",
	"/users",
	"/user/emails",
	"/user/emails",
	"/user/emails",
	"/users/:user/followers",
	"/user/followers",
	"/users/:user/following",
	"/user/following",
	"/user/following/:user",
	"/users/:user/following/:target_user",
	"/user/following/:user",
	"/user/following/:user",
	"/users/:user/keys",
	"/user/keys",
	"/user/keys/:id",
	"/user/keys",
	"/user/keys/:id",
	"/user/keys/:id",
}
