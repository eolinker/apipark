package core

import (
	"net/http"

	"github.com/eolinker/go-common/pm3"
)

func (p *plugin) ProjectApi() []pm3.Api {
	return []pm3.Api{
		// 项目
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/info", []string{"context", "query:project"}, []string{"project"}, p.projectController.GetProject),
		pm3.CreateApiWidthDoc(http.MethodPut, "/api/v1/project/info", []string{"context", "query:project", "body"}, []string{"project"}, p.projectController.EditProject),
		pm3.CreateApiWidthDoc(http.MethodDelete, "/api/v1/project/info", []string{"context", "query:project"}, nil, p.projectController.DeleteProject),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/my_projects", []string{"context", "query:team", "query:keyword"}, []string{"projects"}, p.projectController.SearchMyProjects),

		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/simple/projects/mine", []string{"context", "query:keyword"}, []string{"projects"}, p.projectController.MySimpleProjects),

		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/simple/projects", []string{"context", "query:keyword", "query:partition"}, []string{"projects"}, p.projectController.SimpleProjects),

		// 项目成员相关
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/members", []string{"context", "query:project", "query:keyword"}, []string{"members"}, p.projectMemberController.Members),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/simple/project/members", []string{"context", "query:project"}, []string{"members"}, p.projectMemberController.SimpleMembers),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/members/toadd", []string{"context", "query:project", "query:keyword"}, []string{"members"}, p.projectMemberController.SimpleMembersToAdd),

		pm3.CreateApiWidthDoc(http.MethodPost, "/api/v1/project/member", []string{"context", "query:project", "body"}, nil, p.projectMemberController.AddMember),
		pm3.CreateApiWidthDoc(http.MethodDelete, "/api/v1/project/member", []string{"context", "query:project", "query:user"}, nil, p.projectMemberController.RemoveMember),
		pm3.CreateApiWidthDoc(http.MethodPut, "/api/v1/project/member", []string{"context", "query:project", "query:user", "body"}, nil, p.projectMemberController.EditProjectMember),

		// 项目监控
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/monitor/partitions/simple", []string{"context", "query:project"}, []string{"partitions"}, p.projectMonitorController.MonitorPartitions),

		// 应用相关
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/app/info", []string{"context", "query:app"}, []string{"project"}, p.appController.GetApp),
		pm3.CreateApiWidthDoc(http.MethodDelete, "/api/v1/app", []string{"context", "query:app"}, nil, p.appController.DeleteApp),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/simple/apps", []string{"context", "query:keyword"}, []string{"projects"}, p.appController.SimpleApps),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/simple/apps/mine", []string{"context", "query:keyword"}, []string{"projects"}, p.appController.MySimpleApps),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/my_apps", []string{"context", "query:team", "query:keyword"}, []string{"projects"}, p.appController.SearchMyApps),
		pm3.CreateApiWidthDoc(http.MethodPut, "/api/v1/app/info", []string{"context", "query:app", "body"}, []string{"projects"}, p.appController.UpdateApp),
	}
}
