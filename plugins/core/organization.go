package core

import (
	"net/http"

	"github.com/eolinker/go-common/pm3"
)

func (p *plugin) organizationApi() []pm3.Api {
	return []pm3.Api{
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/manager/organizations", []string{"context", "query:keyword"}, []string{"organizations"}, p.organizationController.Search),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/manager/organization", []string{"context", "query:id"}, []string{"organization"}, p.organizationController.Get),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/manager/organization/:id", []string{"context", "path:id"}, []string{"organization"}, p.organizationController.Get),
		pm3.CreateApiWidthDoc(http.MethodPost, "/api/v1/manager/organization", []string{"context", "body"}, []string{"organization"}, p.organizationController.Create),
		pm3.CreateApiWidthDoc(http.MethodPut, "/api/v1/manager/organization", []string{"context", "query:id", "body"}, []string{"organization"}, p.organizationController.Edit),
		pm3.CreateApiWidthDoc(http.MethodDelete, "/api/v1/manager/organization", []string{"context", "query:id"}, []string{"id"}, p.organizationController.Delete),
		pm3.CreateApiWidthDoc(http.MethodDelete, "/api/v1/manager/organization/:id", []string{"context", "path:id"}, []string{"id"}, p.organizationController.Delete),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/simple/organization/partitions", []string{"context", "query:organization"}, []string{"partitions"}, p.organizationController.Partitions),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/simple/organizations", []string{"context"}, []string{"organizations"}, p.organizationController.Simple),
	}

}
