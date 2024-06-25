package core

import (
	"net/http"

	"github.com/eolinker/go-common/pm3"
)

func (p *plugin) apiApis() []pm3.Api {
	return []pm3.Api{
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/apis", []string{"context", "query:keyword", "query:project"}, []string{"apis"}, p.apiController.Search),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/apis/simple", []string{"context", "query:keyword", "query:project"}, []string{"apis"}, p.apiController.SimpleSearch),
		pm3.CreateApiWidthDoc(http.MethodPost, "/api/v1/simple/project/apis", []string{"context", "query:partition", "body"}, []string{"apis"}, p.apiController.SimpleList),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/api/detail", []string{"context", "query:project", "query:api"}, []string{"api"}, p.apiController.Detail),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/api/detail/simple", []string{"context", "query:project", "query:api"}, []string{"api"}, p.apiController.SimpleDetail),
		pm3.CreateApiWidthDoc(http.MethodPost, "/api/v1/project/api", []string{"context", "query:project", "body"}, []string{"api"}, p.apiController.Create),
		pm3.CreateApiWidthDoc(http.MethodPut, "/api/v1/project/api", []string{"context", "query:project", "query:api", "body"}, []string{"api"}, p.apiController.Edit),
		pm3.CreateApiWidthDoc(http.MethodDelete, "/api/v1/project/api", []string{"context", "query:project", "query:api"}, nil, p.apiController.Delete),
		pm3.CreateApiWidthDoc(http.MethodPost, "/api/v1/project/api/copy", []string{"context", "query:project", "query:api", "body"}, []string{"api"}, p.apiController.Copy),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/api/doc", []string{"context", "query:project", "query:api"}, []string{"api"}, p.apiController.ApiDocDetail),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/api/proxy", []string{"context", "query:project", "query:api"}, []string{"api"}, p.apiController.ApiProxyDetail),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/api/define", []string{"context", "query:project"}, []string{"prefix", "force"}, p.apiController.Prefix),
	}
}
