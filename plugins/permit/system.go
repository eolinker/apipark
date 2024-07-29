package permit

import (
	"github.com/eolinker/go-common/pm3"
)

func (p *pluginPermit) getSystemApis() []pm3.Api {
	return []pm3.Api{
		//pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/system/permissions", []string{"context"}, []string{"system", "team", "project"}, p.systemPermitController.List),
		//pm3.CreateApiWidthDoc(http.MethodPost, "/api/v1/system/permission", []string{"context", "body"}, []string{}, p.systemPermitController.Grant),
		//pm3.CreateApiWidthDoc(http.MethodPost, "/api/v1/system/permission/team", []string{"context", "body"}, []string{}, p.systemPermitController.GrantTemplateForTeam),
		//pm3.CreateApiWidthDoc(http.MethodPost, "/api/v1/system/permission/project", []string{"context", "body"}, []string{}, p.systemPermitController.GrantTemplateForProject),
		//pm3.CreateApiWidthDoc(http.MethodDelete, "/api/v1/system/permission", []string{"context", "query:access", "query:key"}, []string{}, p.systemPermitController.Remove),
		//pm3.CreateApiWidthDoc(http.MethodDelete, "/api/v1/system/permission/team", []string{"context", "query:access", "query:key"}, []string{}, p.systemPermitController.Remove),
		//pm3.CreateApiWidthDoc(http.MethodDelete, "/api/v1/system/permission/project", []string{"context", "query:access", "query:key"}, []string{}, p.systemPermitController.Remove),
		//pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/system/permission/options", []string{"context", "query:keyword"}, []string{"options"}, p.systemPermitController.Options),
		//pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/system/permission/options/team", []string{"context", "query:keyword"}, []string{"options"}, p.systemPermitController.OptionsForTeamTemplate),
		//pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/system/permission/options/project", []string{"context", "query:keyword"}, []string{"options"}, p.systemPermitController.OptionsForProjectTemplate),
		//pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/profile/permission/system", []string{"context"}, []string{"access"}, p.systemPermitController.Permissions),
	}
}
