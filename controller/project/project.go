package project

import (
	"reflect"

	project_monitor_dto "github.com/eolinker/apipark/module/project-monitor/dto"

	project_dto "github.com/eolinker/apipark/module/project/dto"
	"github.com/gin-gonic/gin"

	"github.com/eolinker/go-common/autowire"
)

type IProjectController interface {
	// GetProject 获取项目信息
	GetProject(ctx *gin.Context, id string) (*project_dto.Project, error)
	// SearchMyProjects 搜索项目
	SearchMyProjects(ctx *gin.Context, teamID string, keyword string) ([]*project_dto.ProjectItem, error)
	Search(ctx *gin.Context, teamID string, keyword string) ([]*project_dto.ProjectItem, error)
	// CreateProject 创建项目
	CreateProject(ctx *gin.Context, teamID string, project *project_dto.CreateProject) (*project_dto.Project, error)
	// EditProject 编辑项目
	EditProject(ctx *gin.Context, id string, project *project_dto.EditProject) (*project_dto.Project, error)
	// DeleteProject 删除项目
	DeleteProject(ctx *gin.Context, id string) error
	// SimpleProjects 获取简易项目列表
	SimpleProjects(ctx *gin.Context, keyword string, partition string) ([]*project_dto.SimpleProjectItem, error)
	// MySimpleProjects 获取我的简易项目列表
	MySimpleProjects(ctx *gin.Context, keyword string) ([]*project_dto.SimpleProjectItem, error)
}

type IProjectMemberController interface {
	// Members 获取项目成员列表
	Members(ctx *gin.Context, pid string, keyword string) ([]*project_dto.MemberItem, error)
	// AddMember 添加项目成员
	AddMember(ctx *gin.Context, pid string, users *project_dto.Users) error
	// RemoveMember 移除项目成员
	RemoveMember(ctx *gin.Context, pid string, uid string) error
	// EditProjectMember 修改成员信息
	EditProjectMember(ctx *gin.Context, pid string, uid string, edit *project_dto.EditProjectMember) error
	// SimpleMembers 简易系统成员列表
	SimpleMembers(ctx *gin.Context, pid string) ([]*project_dto.SimpleMemberItem, error)
	SimpleMembersToAdd(ctx *gin.Context, pid string, keyword string) ([]*project_dto.TeamMemberToAdd, error)
}

type IProjectMonitorController interface {
	// MonitorPartitions 获取项目监控分区列表
	MonitorPartitions(ctx *gin.Context, pid string) ([]*project_monitor_dto.MonitorPartition, error)
}

type IAppController interface {
	// CreateApp 创建应用
	CreateApp(ctx *gin.Context, teamID string, project *project_dto.CreateApp) (*project_dto.App, error)

	UpdateApp(ctx *gin.Context, appId string, project *project_dto.UpdateApp) (*project_dto.App, error)
	SearchMyApps(ctx *gin.Context, teamId string, keyword string) ([]*project_dto.AppItem, error)
	// SimpleApps 获取简易项目列表
	SimpleApps(ctx *gin.Context, keyword string) ([]*project_dto.SimpleAppItem, error)
	MySimpleApps(ctx *gin.Context, keyword string) ([]*project_dto.SimpleAppItem, error)
	GetApp(ctx *gin.Context, appId string) (*project_dto.App, error)
	DeleteApp(ctx *gin.Context, appId string) error
}

func init() {
	autowire.Auto[IProjectController](func() reflect.Value {
		return reflect.ValueOf(new(imlProjectController))
	})
	autowire.Auto[IProjectMemberController](func() reflect.Value {
		return reflect.ValueOf(new(imlProjectMemberController))
	})
	autowire.Auto[IProjectMonitorController](func() reflect.Value {
		return reflect.ValueOf(new(imlProjectMonitorController))
	})

	autowire.Auto[IAppController](func() reflect.Value {
		return reflect.ValueOf(new(imlAppController))
	})
}
