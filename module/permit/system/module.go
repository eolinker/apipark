package system

import (
	"context"
	"github.com/eolinker/apipark/module/permit/dto"
	"github.com/eolinker/go-common/autowire"
	"reflect"
)

const (
	systemDomain       = "/"
	teamDomain         = "/template/team"
	projectDomain      = "/template/project"
	accessGroup        = "system"
	teamAccessGroup    = "team"
	projectAccessGroup = "project"
)

type ISystemPermitModule interface {
	// GrantSystem 分配系统级权限
	GrantSystem(ctx context.Context, access, key string) error
	GrantTemplateForTeam(ctx context.Context, access, key string) error
	GrantTemplateForProject(ctx context.Context, access, key string) error
	RemoveSystemAccess(ctx context.Context, access, key string) error
	RemoveTeamTemplateAccess(ctx context.Context, access, key string) error
	RemoveProjectTemplateAccess(ctx context.Context, access, key string) error

	// Access 查询权限授权情况
	SystemAccess(ctx context.Context) ([]*permit_dto.Permission, error)
	TeamAccess(ctx context.Context) ([]*permit_dto.Permission, error)
	ProjectAccess(ctx context.Context) ([]*permit_dto.Permission, error)
	Get(ctx context.Context, access string) (*permit_dto.Permission, error)
	OptionsForSystem(ctx context.Context, keyword string) ([]*permit_dto.Option, error)
	OptionsForTeamTemplate(ctx context.Context, keyword string) ([]*permit_dto.Option, error)
	OptionsForProjectTemplate(ctx context.Context, keyword string) ([]*permit_dto.Option, error)

	Permissions(ctx context.Context) ([]string, error)
}

func init() {
	autowire.Auto[ISystemPermitModule](func() reflect.Value {
		return reflect.ValueOf(new(imlSystemPermitModule))
	})
}
