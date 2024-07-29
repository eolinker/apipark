package permit

import (
	"github.com/eolinker/go-common/pm3"
)

func (p *pluginPermit) getSTeamPermitApis() []pm3.Api {
	return []pm3.Api{
		//pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/team/setting/permissions", []string{"context", "query:team"}, []string{"permissions"}, p.teamPermitController.List),
		//pm3.CreateApiWidthDoc(http.MethodPost, "/api/v1/team/setting/permission", []string{"context", "query:team", "body"}, []string{}, p.teamPermitController.Grant),
		//pm3.CreateApiWidthDoc(http.MethodDelete, "/api/v1/team/setting/permission", []string{"context", "query:team", "query:access", "query:key"}, []string{}, p.teamPermitController.Remove),
		//pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/team/setting/permission/options", []string{"context", "query:team", "query:keyword"}, []string{"options"}, p.teamPermitController.Options),
		//pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/profile/permission/team", []string{"context", "query:team"}, []string{"access"}, p.teamPermitController.Permissions),
	}
}
