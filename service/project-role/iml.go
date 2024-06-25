package project_role

import (
	"context"
	"github.com/eolinker/go-common/utils"
	"time"

	"github.com/eolinker/apipark/stores/project"
)

var (
	_ IProjectRoleService = (*imlProjectRoleService)(nil)
)

type imlProjectRoleService struct {
	projectRoleStore project.IMemberRoleStore `autowired:""`
}

func (i *imlProjectRoleService) AddRole(ctx context.Context, pid string, uid string, rid string) error {
	return i.projectRoleStore.Insert(ctx, &project.MemberRole{
		Pid:        pid,
		Uid:        uid,
		Rid:        rid,
		CreateTime: time.Now(),
	})
}

func (i *imlProjectRoleService) DeleteRole(ctx context.Context, pid string, uid string) error {
	_, err := i.projectRoleStore.DeleteWhere(ctx, map[string]interface{}{
		"pid": pid,
		"uid": uid,
	})
	return err
}

func (i *imlProjectRoleService) Count(ctx context.Context, condition map[string]interface{}) (int64, error) {
	return i.projectRoleStore.CountWhere(ctx, condition)
}

func (i *imlProjectRoleService) Roles(ctx context.Context, pid string, uid ...string) ([]*ProjectRole, error) {
	condition := map[string]interface{}{
		"pid": pid,
	}
	if len(uid) > 0 {
		condition["uid"] = uid
	}
	roles, err := i.projectRoleStore.List(ctx, condition)
	if err != nil {
		return nil, err
	}
	return utils.SliceToSlice(roles, func(role *project.MemberRole) *ProjectRole {
		return &ProjectRole{
			Pid:        role.Pid,
			Uid:        role.Uid,
			Rid:        role.Rid,
			CreateTime: role.CreateTime,
		}
	}), nil
}

func (i *imlProjectRoleService) RoleMap(ctx context.Context, pid string, uid ...string) (map[string][]*ProjectRole, error) {
	roles, err := i.Roles(ctx, pid, uid...)
	if err != nil {
		return nil, err
	}
	return utils.SliceToMapArray(roles, func(t *ProjectRole) string {
		return t.Uid
	}), nil

}

func (i *imlProjectRoleService) OnComplete() {

}
