package system

import (
	"context"
	"errors"
	"github.com/eolinker/ap-account/service/role"
	"github.com/eolinker/ap-account/service/user"
	user_group "github.com/eolinker/ap-account/service/user-group"
	permit_identity "github.com/eolinker/apipark/middleware/permit/identity"
	"github.com/eolinker/apipark/module/permit/dto"
	permit_type "github.com/eolinker/apipark/service/permit-type"
	"github.com/eolinker/go-common/access"
	"github.com/eolinker/go-common/autowire"
	"github.com/eolinker/go-common/permit"
	"github.com/eolinker/go-common/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"reflect"
)

var (
	_ ISystemPermitModule                   = (*imlSystemPermitModule)(nil)
	_ permit_identity.IdentitySystemService = (*imlSystemPermitModule)(nil)
	_ autowire.Complete                     = (*imlSystemPermitModule)(nil)
)

var (
	specialRolesForSystem = []*permit_dto.Option{permit_type.AnyOne}
	specialRolesTeam      = []*permit_dto.Option{
		permit_type.AnyOne,
		permit_type.TeamMember,
		permit_type.TeamMaster,
	}
	specialRolesProject = []*permit_dto.Option{
		permit_type.AnyOne,
		permit_type.ProjectMember,
		permit_type.ProjectMaster,
		permit_type.TeamMaster,
		permit_type.TeamMember,
	}
)

type imlSystemPermitModule struct {
	permitService          permit.IPermit                     `autowired:""`
	userGroupService       user_group.IUserGroupService       `autowired:""`
	userGroupMemberService user_group.IUserGroupMemberService `autowired:""`
	userService            user.IUserService                  `autowired:""`
	roleService            role.IRoleService                  `autowired:""`
}

func (m *imlSystemPermitModule) Permissions(ctx context.Context) ([]string, error) {

	uid := utils.UserId(ctx)
	if uid == "admin" {
		accessList, _ := access.Get(permit_identity.SystemGroup)

		return utils.SliceToSlice(accessList, func(s access.Access) string {
			return s.Name
		}), nil
	}
	identitySystem, err := m.IdentifySystem(ctx, uid)
	if err != nil {
		return nil, err
	}
	return m.permitService.Access(ctx, systemDomain, identitySystem...)
}

func (m *imlSystemPermitModule) OnComplete() {
	autowire.Inject[permit_identity.IdentitySystemService](m)
	permit.AddDomainHandler(permit_identity.SystemGroup, m.domain)
}

func (m *imlSystemPermitModule) RemoveTeamTemplateAccess(ctx context.Context, access, key string) error {
	return m.permitService.Remove(ctx, access, teamDomain, key)

}

func (m *imlSystemPermitModule) RemoveProjectTemplateAccess(ctx context.Context, access, key string) error {
	return m.permitService.Remove(ctx, access, projectDomain, key)
}
func (m *imlSystemPermitModule) roleOptions(ctx context.Context, keyword string) ([]*permit_dto.Option, error) {
	roles, err := m.roleService.Search(ctx, keyword)
	if err == nil && len(roles) > 0 {
		// 角色
		rolesO := utils.SliceToSlice(roles, func(r *role.Role) *permit_dto.Option {
			return permit_type.Role.Target(r.Id, r.Name)
		})
		return permit_dto.SearchOptions(rolesO, keyword), nil

	}
	return nil, nil
}
func (m *imlSystemPermitModule) userGroupOptions(ctx context.Context, keyword string) ([]*permit_dto.Option, error) {
	ugs, err := m.userGroupService.Search(ctx, keyword)
	if err != nil {
		return nil, err
	}
	ugsO := utils.SliceToSlice(ugs, func(r *user_group.UserGroup) *permit_dto.Option {
		return permit_type.UserGroup.Target(r.Id, r.Name)
	})
	return permit_dto.SearchOptions(ugsO, keyword), nil
}
func (m *imlSystemPermitModule) OptionsForSystem(ctx context.Context, keyword string) ([]*permit_dto.Option, error) {
	result := make([]*permit_dto.Option, 0, 10)

	specialRoles := permit_dto.SearchOptions(specialRolesForSystem, keyword)
	if len(specialRoles) > 0 {
		result = append(result, specialRoles...)
	}
	userGroupOptions, err := m.userGroupOptions(ctx, keyword)
	if err != nil {
		return nil, err
	}
	if len(userGroupOptions) > 0 {
		result = append(result, userGroupOptions...)
	}
	//roleOptions, err := m.roleOptions(ctx, keyword)
	//if err != nil {
	//	return nil, err
	//}
	//result = append(result, roleOptions...)
	if keyword != "" {
		ul, err := m.userService.Search(ctx, keyword, -1)
		if err == nil && len(ul) > 0 {
			result = append(result, utils.SliceToSlice(ul, func(r *user.User) *permit_dto.Option {
				return permit_type.User.Target(r.UID, r.Username)

			})...)
		}
	}
	return result, nil

}

func (m *imlSystemPermitModule) OptionsForTeamTemplate(ctx context.Context, keyword string) ([]*permit_dto.Option, error) {
	result := make([]*permit_dto.Option, 0, 10)

	specialRoles := permit_dto.SearchOptions(specialRolesTeam, keyword)
	if len(specialRoles) > 0 {
		result = append(result, specialRoles...)
	}
	userGroupOptions, err := m.userGroupOptions(ctx, keyword)
	if err != nil {
		return nil, err
	}
	if len(userGroupOptions) > 0 {
		result = append(result, userGroupOptions...)
	}
	roleOptions, err := m.roleOptions(ctx, keyword)
	if err != nil {
		return nil, err
	}
	result = append(result, roleOptions...)
	return result, nil
}

func (m *imlSystemPermitModule) OptionsForProjectTemplate(ctx context.Context, keyword string) ([]*permit_dto.Option, error) {
	result := make([]*permit_dto.Option, 0, 10)
	result = append(result, permit_dto.SearchOptions(specialRolesProject, keyword)...)

	userGroupOptions, err := m.userGroupOptions(ctx, keyword)
	if err == nil && len(userGroupOptions) > 0 {
		result = append(result, userGroupOptions...)
	}

	roleOptions, err := m.roleOptions(ctx, keyword)
	if err == nil && len(roleOptions) > 0 {
		result = append(result, roleOptions...)
	}
	return result, nil
}

func (m *imlSystemPermitModule) GrantTemplateForTeam(ctx context.Context, access, key string) error {
	return m.permitService.Add(ctx, access, teamDomain, key)
}

func (m *imlSystemPermitModule) GrantTemplateForProject(ctx context.Context, access, key string) error {
	return m.permitService.Add(ctx, access, projectDomain, key)
}

func (m *imlSystemPermitModule) TeamAccess(ctx context.Context) ([]*permit_dto.Permission, error) {

	accesses, has := access.Get(teamAccessGroup)
	if !has {
		return nil, errors.New("no access for team")
	}

	grants, err := m.permitService.GrantForDomain(ctx, teamDomain)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	targets := make([]*permit_type.Target, 0, len(grants))

	result := utils.SliceToSlice(accesses, func(s access.Access) *permit_dto.Permission {

		r := &permit_dto.Permission{
			Access:      s.Name,
			Name:        s.CName,
			Description: s.Desc,
			Grant:       nil,
		}
		r.Grant = permit_type.TargetsOf(grants[r.Access]...)
		targets = append(targets, r.Grant...)
		return r
	})

	permit_type.CompleteLabels(ctx, targets...)
	return result, nil
}

func (m *imlSystemPermitModule) ProjectAccess(ctx context.Context) ([]*permit_dto.Permission, error) {

	accesses, has := access.Get(projectAccessGroup)
	if !has {
		return nil, errors.New("no access for team")
	}

	grants, err := m.permitService.GrantForDomain(ctx, projectDomain)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	targets := make([]*permit_type.Target, 0, len(grants))

	result := utils.SliceToSlice(accesses, func(s access.Access) *permit_dto.Permission {

		r := &permit_dto.Permission{
			Access:      s.Name,
			Name:        s.CName,
			Description: s.Desc,
			Grant:       nil,
		}
		r.Grant = permit_type.TargetsOf(grants[r.Access]...)
		targets = append(targets, r.Grant...)
		return r
	})

	permit_type.CompleteLabels(ctx, targets...)
	return result, nil
}

func (m *imlSystemPermitModule) domain(ctx *gin.Context) ([]string, []string, bool) {

	system, err := m.IdentifySystem(ctx, utils.UserId(ctx))
	if err != nil {
		return nil, nil, false
	}
	return []string{systemDomain}, system, true
}
func (m *imlSystemPermitModule) IdentifySystem(ctx context.Context, uid string) ([]string, error) {
	if uid == "" {
		return nil, errors.New("not login")
	}
	targets := make([]string, 0)
	targets = append(targets, permit_type.User.KeyOf(uid))
	targets = append(targets, permit_type.Special.KeyOf("all"))
	members, err := m.userGroupMemberService.FilterMembersForUser(ctx, uid)
	if err != nil && len(members) > 0 {
		cs := members[uid]
		for _, c := range cs {
			targets = append(targets, permit_type.UserGroup.KeyOf(c))
		}
	}

	return targets, nil
}

func (m *imlSystemPermitModule) GrantSystem(ctx context.Context, access, key string) error {
	return m.permitService.Add(ctx, access, systemDomain, key)
}

func (m *imlSystemPermitModule) RemoveSystemAccess(ctx context.Context, access, key string) error {
	return m.permitService.Remove(ctx, access, systemDomain, key)
}

func (m *imlSystemPermitModule) SystemAccess(ctx context.Context) ([]*permit_dto.Permission, error) {

	accesses, has := access.Get(accessGroup)
	if !has {
		return nil, errors.New("no access for system")
	}

	result := utils.SliceToSlice(accesses, func(s access.Access) *permit_dto.Permission {
		return &permit_dto.Permission{
			Access:      s.Name,
			Name:        s.CName,
			Description: s.Desc,
			Grant:       nil,
		}
	})
	grants, err := m.permitService.GrantForDomain(ctx, systemDomain)
	if err != nil {
		return nil, err
	}
	targets := make([]*permit_type.Target, 0, len(grants))
	for _, r := range result {
		r.Grant = permit_type.TargetsOf(grants[r.Access]...)
		targets = append(targets, r.Grant...)
	}

	permit_type.CompleteLabels(ctx, targets...)
	return result, nil
}

func (m *imlSystemPermitModule) Get(ctx context.Context, ac string) (*permit_dto.Permission, error) {

	accesses, has := access.Get(accessGroup)
	if !has {
		return nil, errors.New("no access for system")
	}
	var out *permit_dto.Permission

	for _, a := range accesses {
		if a.Name == ac {
			out = &permit_dto.Permission{
				Access:      a.Name,
				Name:        a.CName,
				Description: a.Desc,
				Grant:       nil,
			}

			break
		}
	}

	if out == nil {
		return nil, errors.New("no access for system")
	}
	granted, err := m.permitService.Granted(ctx, ac, systemDomain)
	if err != nil {
		return nil, err
	}

	out.Grant = permit_type.TargetsOf(granted...)
	permit_type.CompleteLabels(ctx, out.Grant...)
	return out, nil
}

func init() {
	autowire.Auto[ISystemPermitModule](func() reflect.Value {
		return reflect.ValueOf(new(imlSystemPermitModule))
	})
}
