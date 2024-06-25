package core

import (
	"net/http"

	"github.com/eolinker/go-common/pm3"
)

func (p *plugin) projectAuthorizationApis() []pm3.Api {
	return []pm3.Api{
		pm3.CreateApiWidthDoc(http.MethodPost, "/api/v1/project/authorization", []string{"context", "query:project", "body"}, []string{"authorization"}, p.projectAuthorizationController.AddAuthorization),
		pm3.CreateApiWidthDoc(http.MethodPut, "/api/v1/project/authorization", []string{"context", "query:project", "query:authorization", "body"}, []string{"authorization"}, p.projectAuthorizationController.EditAuthorization),
		pm3.CreateApiWidthDoc(http.MethodDelete, "/api/v1/project/authorization", []string{"context", "query:project", "query:authorization"}, nil, p.projectAuthorizationController.DeleteAuthorization),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/authorization", []string{"context", "query:project", "query:authorization"}, []string{"authorization"}, p.projectAuthorizationController.Info),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/authorizations", []string{"context", "query:project"}, []string{"authorizations"}, p.projectAuthorizationController.Authorizations),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/authorization/details", []string{"context", "query:project", "query:authorization"}, []string{"details"}, p.projectAuthorizationController.Detail),
	}
}
