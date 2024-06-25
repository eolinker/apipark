package project_role

import (
	"context"
	"reflect"

	"github.com/eolinker/go-common/autowire"
)

type IProjectRoleService interface {
	RoleMap(ctx context.Context, pid string, uid ...string) (map[string][]*ProjectRole, error)
	Roles(ctx context.Context, pid string, uid ...string) ([]*ProjectRole, error)
	AddRole(ctx context.Context, pid string, uid string, rid string) error
	DeleteRole(ctx context.Context, pid string, uid string) error
	Count(ctx context.Context, condition map[string]interface{}) (int64, error)
}

func init() {
	autowire.Auto[IProjectRoleService](func() reflect.Value {
		return reflect.ValueOf(new(imlProjectRoleService))
	})
}
