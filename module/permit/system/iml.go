package system

import (
	"context"
	"reflect"

	"github.com/eolinker/ap-account/service/role"
	"github.com/eolinker/go-common/autowire"
	"github.com/eolinker/go-common/permit"
	"github.com/eolinker/go-common/utils"
)

var (
	_ ISystemPermitModule = (*imlSystemPermitModule)(nil)
	_ autowire.Complete   = (*imlSystemPermitModule)(nil)
)

type imlSystemPermitModule struct {
	permitService     permit.IPermit          `autowired:""`
	roleService       role.IRoleService       `autowired:""`
	roleMemberService role.IRoleMemberService `autowired:""`
}

func (m *imlSystemPermitModule) Permissions(ctx context.Context) ([]string, error) {

	uid := utils.UserId(ctx)

	roleMembers, err := m.roleMemberService.List(ctx, role.SystemTarget(), uid)
	if err != nil {
		return nil, err
	}
	if len(roleMembers) == 0 {
		return []string{}, nil
	}
	roleIds := utils.SliceToSlice(roleMembers, func(rm *role.Member) string {
		return rm.Role
	})
	roles, err := m.roleService.List(ctx, roleIds...)
	if err != nil {
		return nil, err
	}
	permits := make(map[string]struct{})
	for _, r := range roles {
		for _, p := range r.Permit {
			permits[p] = struct{}{}
		}
	}

	return utils.MapToSlice(permits, func(k string, v struct{}) string {
		return k
	}), nil
}

func (m *imlSystemPermitModule) OnComplete() {
}

func init() {
	autowire.Auto[ISystemPermitModule](func() reflect.Value {
		return reflect.ValueOf(new(imlSystemPermitModule))
	})
}
