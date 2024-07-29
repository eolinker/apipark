package team

import (
	"context"
	"errors"
	"fmt"

	"github.com/eolinker/ap-account/service/account"
	"github.com/eolinker/ap-account/service/member"
	"github.com/eolinker/ap-account/service/role"
	"github.com/eolinker/ap-account/service/user"
	user_group "github.com/eolinker/ap-account/service/user-group"
	permit_identity "github.com/eolinker/apipark/middleware/permit/identity"
	permit_dto "github.com/eolinker/apipark/module/permit/dto"
	permit_type "github.com/eolinker/apipark/service/permit-type"
	"github.com/eolinker/apipark/service/team"
	team_member "github.com/eolinker/apipark/service/team-member"
	"github.com/eolinker/go-common/access"
	"github.com/eolinker/go-common/autowire"
	"github.com/eolinker/go-common/permit"
	"github.com/eolinker/go-common/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	_ ITeamPermitModule                   = (*imlTeamPermitModule)(nil)
	_ permit_identity.IdentityTeamService = (*imlTeamPermitModule)(nil)
	_ autowire.Complete                   = (*imlTeamPermitModule)(nil)
)

const (
	templateDomain = "/template/team"
)

var (
	specialRolesTeam = []*permit_dto.Option{
		permit_type.AnyOne,
		permit_type.ProjectMaster,
		permit_type.ProjectMember,
	}
)

type imlTeamPermitModule struct {
	permitService          permit.IPermit                     `autowired:""`
	teamService            team.ITeamService                  `autowired:""`
	teamMemberService      team_member.ITeamMemberService     `autowired:""`
	accountService         account.IAccountService            `autowired:""`
	userService            user.IUserService                  `autowired:""`
	userGroupService       user_group.IUserGroupService       `autowired:""`
	userGroupMemberService user_group.IUserGroupMemberService `autowired:""`
	roleService            role.IRoleService                  `autowired:""`
}

func (m *imlTeamPermitModule) Permissions(ctx context.Context, teamId string) ([]string, error) {
	uid := utils.UserId(ctx)
	identifyTeam, err := m.IdentifyTeam(ctx, teamId, uid)
	if err != nil {
		return nil, err
	}
	teamDomain := fmt.Sprintf("/%s", teamId)
	accessList, _ := access.Get(permit_identity.TeamGroup)

	// 受团队内权限设置生效的的权限
	teamGranted, _ := m.permitService.GrantForDomain(ctx, teamDomain)
	if teamGranted == nil {
		teamGranted = make(map[string][]string)
	}
	// 排除掉项目生效的权限,剩下的是模版生效的权限
	templateAccess := utils.SliceToSlice(accessList, func(s access.Access) string {
		return s.Name
	}, func(a access.Access) bool {
		_, has := teamGranted[a.Name]
		return !has
	})
	myAccess := make([]string, 0, len(accessList))
	if len(teamGranted) > 0 {
		projectAccess, _ := m.permitService.Access(ctx, teamDomain, identifyTeam...)
		myAccess = append(myAccess, projectAccess...)
	}
	if len(templateAccess) > 0 {
		teamAccess, _ := m.permitService.Access(ctx, templateDomain, identifyTeam...)
		myAccess = append(myAccess, teamAccess...)
	}
	return myAccess, nil
}

func (m *imlTeamPermitModule) OnComplete() {
	autowire.Inject[permit_identity.IdentityTeamService](m)
	permit.AddDomainHandler(permit_identity.TeamGroup, m.domainHandler)
}

func (m *imlTeamPermitModule) IdentifyTeam(ctx context.Context, team string, uid string) ([]string, error) {
	t, err := m.teamService.Get(ctx, team)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("team not found")
		}
	}
	result := make([]string, 0, 5)
	if t.Master == uid {
		result = append(result, permit_type.TeamMaster.Key)
		result = append(result, permit_type.TeamMember.Key)
	} else {
		members, err := m.teamMemberService.Members(ctx, []string{t.Id}, []string{uid})
		if err == nil && len(members) > 0 {
			result = append(result, permit_type.TeamMember.Key)
		}

	}
	members, err := m.userGroupMemberService.FilterMembersForUser(ctx, uid)
	if err != nil && len(members) > 0 {
		cs := members[uid]
		for _, c := range cs {
			result = append(result, permit_type.UserGroup.KeyOf(c))
		}
	}

	return result, nil

}

func (m *imlTeamPermitModule) Options(ctx context.Context, teamId string, keyword string) ([]*permit_dto.Option, error) {
	t, err := m.teamService.Get(ctx, teamId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("team not found")
		}
		return nil, err
	}
	result := make([]*permit_dto.Option, 0, 10)
	specialRoles := permit_dto.SearchOptions(specialRolesTeam, keyword)
	if len(specialRoles) > 0 {
		result = append(result, specialRoles...)
	}
	ugs, err := m.userGroupService.Search(ctx, keyword)
	if err == nil && len(ugs) > 0 {
		// 用户组
		ugsO := utils.SliceToSlice(ugs, func(r *user_group.UserGroup) *permit_dto.Option {
			return permit_type.UserGroup.Target(r.Id, r.Name)
		})
		ugsO = permit_dto.SearchOptions(ugsO, keyword)
		if len(ugsO) > 0 {
			result = append(result, ugsO...)
		}
	}

	roles, err := m.roleService.Search(ctx, keyword)
	if err == nil && len(roles) > 0 {
		// 角色
		rolesO := utils.SliceToSlice(roles, func(r *role.Role) *permit_dto.Option {
			return permit_type.Role.Target(r.Id, r.Name)
		})
		rolesO = permit_dto.SearchOptions(rolesO, keyword)
		if len(rolesO) > 0 {
			result = append(result, rolesO...)
		}
	}
	// keyword 不为空时,搜索团队成员
	if keyword != "" {
		users, _ := m.userService.Search(ctx, keyword, -1)
		if len(users) > 0 {
			// 团队成员
			members, err := m.teamMemberService.Members(ctx, []string{t.Id}, utils.SliceToSlice(users, func(r *user.User) string {
				return r.UID
			}))
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			if len(members) > 0 {
				userMap := utils.SliceToMap(users, func(r *user.User) string {
					return r.UID
				})
				ugsO := utils.SliceToSlice(members, func(r *member.Member) *permit_dto.Option {
					us := userMap[r.UID]
					return permit_type.User.Target(us.UID, us.GetLabel())

				}, func(m *member.Member) bool {
					_, has := userMap[m.UID]
					return has
				})
				ugsO = permit_dto.SearchOptions(ugsO, keyword)
				if len(ugsO) > 0 {
					result = append(result, ugsO...)
				}
			}

		}
	}

	return result, nil
}

func (m *imlTeamPermitModule) Grant(ctx *gin.Context, teamId, access string, key string) error {

	domain, err := m.domain(ctx, teamId)
	if err != nil {
		return fmt.Errorf("team:%w", err)
	}

	granted, err := m.permitService.Granted(ctx, access, domain)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if len(granted) == 0 {
		tmpGranted, err := m.permitService.Granted(ctx, access, templateDomain)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		err = m.permitService.Reset(ctx, access, domain, tmpGranted...)
		if err != nil {
			return err
		}
	}

	return m.permitService.Add(ctx, access, domain, key)
}

func (m *imlTeamPermitModule) Remove(ctx *gin.Context, teamId string, access string, key string) error {
	domain, err := m.domain(ctx, teamId)
	if err != nil {
		return fmt.Errorf("team:%w", err)
	}
	return m.permitService.Remove(ctx, access, domain, key)
}

func (m *imlTeamPermitModule) TeamAccess(ctx *gin.Context, teamId string) ([]*permit_dto.Permission, error) {
	domain, err := m.domain(ctx, teamId)
	if err != nil {
		return nil, fmt.Errorf("team:%w", err)
	}

	accesses, has := access.Get(accessGroup)
	if !has {
		return nil, errors.New("no access for team")
	}

	grants, err := m.permitService.GrantForDomain(ctx, domain)
	if err != nil {
		return nil, err
	}
	templateGrants, err := m.permitService.GrantForDomain(ctx, templateDomain)
	if err != nil {
		return nil, err
	}
	targets := make([]*permit_type.Target, 0, len(grants))

	result := utils.SliceToSlice(accesses, func(s access.Access) *permit_dto.Permission {
		r := &permit_dto.Permission{
			Access: s.Name,
			Name:   s.CName,
			//Description: s.Desc,
			Grant: nil,
		}
		if gs, has := grants[r.Access]; has {
			r.Grant = permit_type.TargetsOf(gs...)
		} else {
			r.Grant = permit_type.TargetsOf(templateGrants[r.Access]...)
		}
		targets = append(targets, r.Grant...)
		return r
	})
	permit_type.CompleteLabels(ctx, targets...)
	return result, nil
}

func (m *imlTeamPermitModule) domain(ctx context.Context, teamId string) (string, error) {
	_, err := m.teamService.Get(ctx, teamId)
	if err != nil {
		return "", err
	}
	return fmt.Sprint("/", teamId), nil
}

func (m *imlTeamPermitModule) domainHandler(ctx *gin.Context) ([]string, []string, bool) {
	teamId := ctx.Query("team")

	t, err := m.teamService.Get(ctx, teamId)
	if err != nil {
		return nil, nil, false
	}
	uid := utils.UserId(ctx)

	teamIdentity, err := m.IdentifyTeam(ctx, teamId, uid)
	if err != nil {
		return nil, nil, false
	}
	return []string{
		fmt.Sprintf("/%s", t.Id),
		templateDomain,
	}, teamIdentity, true

}
