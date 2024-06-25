package project

import (
	project_monitor "github.com/eolinker/apipark/module/project-monitor"
	project_monitor_dto "github.com/eolinker/apipark/module/project-monitor/dto"

	"github.com/eolinker/apipark/module/project"
	project_dto "github.com/eolinker/apipark/module/project/dto"
	"github.com/gin-gonic/gin"
)

var (
	_ IProjectController        = (*imlProjectController)(nil)
	_ IProjectMemberController  = (*imlProjectMemberController)(nil)
	_ IProjectMonitorController = (*imlProjectMonitorController)(nil)
	_ IAppController            = (*imlAppController)(nil)
)

type imlProjectController struct {
	module project.IProjectModule `autowired:""`
}

func (i *imlProjectController) SearchMyProjects(ctx *gin.Context, teamId string, keyword string) ([]*project_dto.ProjectItem, error) {
	return i.module.SearchMyProjects(ctx, teamId, keyword)
}

func (i *imlProjectController) SimpleProjects(ctx *gin.Context, keyword string, partition string) ([]*project_dto.SimpleProjectItem, error) {
	return i.module.SimpleProjects(ctx, keyword, partition)
}

func (i *imlProjectController) MySimpleProjects(ctx *gin.Context, keyword string) ([]*project_dto.SimpleProjectItem, error) {
	return i.module.MySimpleProjects(ctx, keyword)
}

func (i *imlProjectController) GetProject(ctx *gin.Context, id string) (*project_dto.Project, error) {
	return i.module.GetProject(ctx, id)
}

func (i *imlProjectController) Search(ctx *gin.Context, teamID string, keyword string) ([]*project_dto.ProjectItem, error) {
	return i.module.Search(ctx, teamID, keyword)
}

func (i *imlProjectController) CreateProject(ctx *gin.Context, teamID string, project *project_dto.CreateProject) (*project_dto.Project, error) {
	return i.module.CreateProject(ctx, teamID, project)
}

func (i *imlProjectController) EditProject(ctx *gin.Context, id string, project *project_dto.EditProject) (*project_dto.Project, error) {
	return i.module.EditProject(ctx, id, project)
}

func (i *imlProjectController) DeleteProject(ctx *gin.Context, id string) error {
	return i.module.DeleteProject(ctx, id)
}

type imlProjectMemberController struct {
	module project.IProjectMemberModule `autowired:""`
}

func (i *imlProjectMemberController) SimpleMembersToAdd(ctx *gin.Context, pid string, keyword string) ([]*project_dto.TeamMemberToAdd, error) {
	return i.module.SimpleMembersToAdd(ctx, pid, keyword)
}

func (i *imlProjectMemberController) SimpleMembers(ctx *gin.Context, pid string) ([]*project_dto.SimpleMemberItem, error) {
	return i.module.SimpleMembers(ctx, pid)
}

func (i *imlProjectMemberController) Members(ctx *gin.Context, id string, keyword string) ([]*project_dto.MemberItem, error) {
	return i.module.Members(ctx, id, keyword)
}

func (i *imlProjectMemberController) AddMember(ctx *gin.Context, pid string, users *project_dto.Users) error {
	return i.module.AddMember(ctx, pid, users.Users)
}

func (i *imlProjectMemberController) RemoveMember(ctx *gin.Context, pid string, uid string) error {
	return i.module.RemoveMember(ctx, pid, []string{uid})
}

func (i *imlProjectMemberController) EditProjectMember(ctx *gin.Context, pid string, uid string, edit *project_dto.EditProjectMember) error {
	return i.module.EditProjectMember(ctx, pid, uid, edit.Roles)
}

type imlProjectMonitorController struct {
	module project_monitor.IProjectMonitor `autowired:""`
}

func (i *imlProjectMonitorController) MonitorPartitions(ctx *gin.Context, pid string) ([]*project_monitor_dto.MonitorPartition, error) {
	return i.module.MonitorPartitions(ctx, pid)
}

type imlAppController struct {
	module project.IAppModule `autowired:""`
}

func (i *imlAppController) CreateApp(ctx *gin.Context, teamID string, input *project_dto.CreateApp) (*project_dto.App, error) {
	return i.module.CreateApp(ctx, teamID, input)
}
func (i *imlAppController) UpdateApp(ctx *gin.Context, appId string, input *project_dto.UpdateApp) (*project_dto.App, error) {
	return i.module.UpdateApp(ctx, appId, input)
}

func (i *imlAppController) SearchMyApps(ctx *gin.Context, teamId string, keyword string) ([]*project_dto.AppItem, error) {
	return i.module.SearchMyApps(ctx, teamId, keyword)
}

func (i *imlAppController) SimpleApps(ctx *gin.Context, keyword string) ([]*project_dto.SimpleAppItem, error) {
	return i.module.SimpleApps(ctx, keyword)
}

func (i *imlAppController) MySimpleApps(ctx *gin.Context, keyword string) ([]*project_dto.SimpleAppItem, error) {
	return i.module.MySimpleApps(ctx, keyword)
}

func (i *imlAppController) GetApp(ctx *gin.Context, appId string) (*project_dto.App, error) {
	return i.module.GetApp(ctx, appId)
}

func (i *imlAppController) DeleteApp(ctx *gin.Context, appId string) error {
	return i.module.DeleteApp(ctx, appId)
}
