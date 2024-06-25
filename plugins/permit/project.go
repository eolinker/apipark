package permit

import (
	"github.com/eolinker/go-common/pm3"
	"net/http"
)

func (p *pluginPermit) getSProjectPermitApis() []pm3.Api {
	return []pm3.Api{
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/setting/permissions", []string{"context", "query:project"}, []string{"permissions"}, p.projectPermitController.List),
		pm3.CreateApiWidthDoc(http.MethodPost, "/api/v1/project/setting/permission", []string{"context", "query:project", "body"}, []string{}, p.projectPermitController.Grant),
		pm3.CreateApiWidthDoc(http.MethodDelete, "/api/v1/project/setting/permission", []string{"context", "query:project", "query:access", "query:key"}, []string{}, p.projectPermitController.Remove),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/setting/permission/options", []string{"context", "query:project", "query:keyword"}, []string{"options"}, p.projectPermitController.Options),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/profile/permission/project", []string{"context", "query:project"}, []string{"access"}, p.projectPermitController.Permissions),
	}
}
