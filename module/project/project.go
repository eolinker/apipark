package project

import (
	"context"
	"reflect"

	project_dto "github.com/eolinker/apipark/module/project/dto"

	"github.com/eolinker/go-common/autowire"
)

type IProjectModule interface {
	// GetProject 获取项目信息
	GetProject(ctx context.Context, id string) (*project_dto.Project, error)
	// Search 搜索项目
	Search(ctx context.Context, teamID string, keyword string) ([]*project_dto.ProjectItem, error)
	// SearchMyProjects 搜索我的项目列表
	SearchMyProjects(ctx context.Context, teamId string, keyword string) ([]*project_dto.ProjectItem, error)
	// CreateProject 创建项目
	CreateProject(ctx context.Context, teamID string, input *project_dto.CreateProject) (*project_dto.Project, error)
	// EditProject 编辑项目
	EditProject(ctx context.Context, id string, input *project_dto.EditProject) (*project_dto.Project, error)
	// DeleteProject 删除项目
	DeleteProject(ctx context.Context, id string) error
	// SimpleProjects 获取简易项目列表
	SimpleProjects(ctx context.Context, keyword string) ([]*project_dto.SimpleProjectItem, error)

	// MySimpleProjects 获取我的简易项目列表
	MySimpleProjects(ctx context.Context, keyword string) ([]*project_dto.SimpleProjectItem, error)
}

type IAppModule interface {
	CreateApp(ctx context.Context, teamID string, input *project_dto.CreateApp) (*project_dto.App, error)
	UpdateApp(ctx context.Context, appId string, input *project_dto.UpdateApp) (*project_dto.App, error)
	SearchMyApps(ctx context.Context, teamId string, keyword string) ([]*project_dto.AppItem, error)
	// SimpleApps 获取简易项目列表
	SimpleApps(ctx context.Context, keyword string) ([]*project_dto.SimpleAppItem, error)
	MySimpleApps(ctx context.Context, keyword string) ([]*project_dto.SimpleAppItem, error)
	GetApp(ctx context.Context, appId string) (*project_dto.App, error)
	DeleteApp(ctx context.Context, appId string) error
}

type IProjectMemberModule interface {
	// Members 获取项目成员列表
	Members(ctx context.Context, id string, keyword string) ([]*project_dto.MemberItem, error)
	// AddMember 添加项目成员
	AddMember(ctx context.Context, id string, userIDs []string) error
	// RemoveMember 移除项目成员
	RemoveMember(ctx context.Context, id string, userIDs []string) error
	// EditProjectMember 修改成员信息
	EditProjectMember(ctx context.Context, pid string, uid string, roles []string) error
	// SimpleMembers 简易成员列表
	SimpleMembers(ctx context.Context, pid string) ([]*project_dto.SimpleMemberItem, error)
	SimpleMembersToAdd(ctx context.Context, pid string, keyword string) ([]*project_dto.TeamMemberToAdd, error)
}

func init() {
	autowire.Auto[IProjectModule](func() reflect.Value {
		m := new(imlProjectModule)
		return reflect.ValueOf(m)
	})
	autowire.Auto[IProjectMemberModule](func() reflect.Value {
		return reflect.ValueOf(new(imlProjectMemberModule))
	})

	autowire.Auto[IAppModule](func() reflect.Value {
		return reflect.ValueOf(new(imlAppModule))
	})

}
