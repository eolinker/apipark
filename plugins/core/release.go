package core

import (
	"github.com/eolinker/go-common/pm3"
	"net/http"
)

func (p *plugin) releaseApis() []pm3.Api {
	return []pm3.Api{

		pm3.CreateApiWidthDoc(http.MethodPost, "/api/v1/project/release", []string{"context", "query:project", "body"}, []string{}, p.releaseController.Create),
		pm3.CreateApiWidthDoc(http.MethodDelete, "/api/v1/project/release", []string{"context", "query:project", "query:id"}, []string{}, p.releaseController.Delete),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/release", []string{"context", "query:project", "query:id"}, []string{"release"}, p.releaseController.Detail),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/releases", []string{"context", "query:project"}, []string{"releases"}, p.releaseController.List),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/release/preview", []string{"context", "query:project"}, []string{"running", "diff", "complete"}, p.releaseController.Preview),
	}
}
