package project

import (
	"context"
	"errors"
	"fmt"

	"github.com/eolinker/ap-account/service/member"
	"github.com/eolinker/ap-account/service/role"
	"github.com/eolinker/ap-account/service/user"
	user_group "github.com/eolinker/ap-account/service/user-group"
	permit_identity "github.com/eolinker/apipark/middleware/permit/identity"
	permit_dto "github.com/eolinker/apipark/module/permit/dto"
	permit_type "github.com/eolinker/apipark/service/permit-type"
	"github.com/eolinker/apipark/service/service"
	"github.com/eolinker/apipark/service/team"
	team_member "github.com/eolinker/apipark/service/team-member"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/go-common/access"
	"github.com/eolinker/go-common/autowire"
	"github.com/eolinker/go-common/permit"
	"github.com/eolinker/go-common/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	templateDomain = "/template/project"
)

var (
	_ IProjectPermitModule                   = (*imlProjectPermitModule)(nil)
	_ permit_identity.IdentityProjectService = (*imlProjectPermitModule)(nil)
	_ autowire.Complete                      = (*imlProjectPermitModule)(nil)
)

var (
	specialRolesProject = []*permit_dto.Option{
		permit_type.AnyOne,
		permit_type.TeamMaster,
		permit_type.ProjectMaster,
		permit_type.TeamMember,
		permit_type.ProjectMember,
	}
)

type imlProjectPermitModule struct {
	projectService service.IServiceService `autowired:""`
	teamProject    team.ITeamService       `autowired:""`
	permitService  permit.IPermit          `autowired:""`
	//projectMemberService   project_member.IMemberService       `autowired:""`
	teamMemberService   team_member.ITeamMemberService      `autowired:""`
	identityTeamService permit_identity.IdentityTeamService `autowired:""`
	//userGroupMemberService user_group.IUserGroupMemberService  `autowired:""`
	userGroupService user_group.IUserGroupService `autowired:""`
	userService      user.IUserService            `autowired:""`
	roleService      role.IRoleService            `autowired:""`
	//projectRoleService     project_role.IProjectRoleService    `autowired:""`
}

func (m *imlProjectPermitModule) Permissions(ctx *gin.Context, projectId string) ([]string, error) {
	pro, err := m.projectService.Get(ctx, projectId)
	if err != nil {
		return nil, err
	}
	uid := utils.UserId(ctx)
	accessList, _ := access.Get(permit_identity.ProjectGroup)
	if uid == "admin" {
		return utils.SliceToSlice(accessList, func(s access.Access) string {
			return s.Name
		}), nil
	}
	projectIdentity, err := m.getIdentity(ctx, pro, uid)
	if err != nil {
		return nil, err
	}

	projectDomain := fmt.Sprintf("/%s/%s", pro.Team, projectId)

	// 受项目内权限设置生效的的权限
	projectGranted, _ := m.permitService.GrantForDomain(ctx, projectDomain)
	if projectGranted == nil {
		projectGranted = make(map[string][]string)
	}
	// 排除掉项目生效的权限,剩下的是模版生效的权限
	templateAccess := utils.SliceToSlice(accessList, func(s access.Access) string {
		return s.Name
	}, func(a access.Access) bool {
		_, has := projectGranted[a.Name]
		return !has
	})
	myAccess := make([]string, 0, len(accessList))
	if len(projectGranted) > 0 {
		projectAccess, _ := m.permitService.Access(ctx, projectDomain, projectIdentity...)
		myAccess = append(myAccess, projectAccess...)
	}
	if len(templateAccess) > 0 {
		teamAccess, _ := m.permitService.Access(ctx, templateDomain, projectIdentity...)
		myAccess = append(myAccess, teamAccess...)
	}
	return myAccess, nil
}

func (m *imlProjectPermitModule) OnComplete() {
	autowire.Inject[permit_identity.IdentityProjectService](m)
	permit.AddDomainHandler(permit_identity.ProjectGroup, m.domainHandler)
}

func (m *imlProjectPermitModule) domainHandler(ctx *gin.Context) ([]string, []string, bool) {
	projectId := ctx.Query("project")
	p, err := m.projectService.Get(ctx, projectId)
	if err != nil {
		return nil, nil, false
	}

	uid := utils.UserId(ctx)

	projectIdentity, err := m.getIdentity(ctx, p, uid)
	if err != nil {
		return nil, nil, false
	}

	return []string{
		fmt.Sprint("/", p.Team, "/", p.Name),
		fmt.Sprintf("/%s", p.Team),
		templateDomain,
	}, projectIdentity, true
}

func (m *imlProjectPermitModule) IdentifyProject(ctx context.Context, project string, uid string) ([]string, error) {

	p, err := m.projectService.Get(ctx, project)
	if err != nil {
		return nil, err
	}
	identities, err := m.getIdentity(ctx, p, uid)
	if err != nil {
		return nil, nil
	}

	return identities, nil
}
func (m *imlProjectPermitModule) getIdentity(ctx context.Context, p *service.Service, uid string) ([]string, error) {
	members, err := m.projectMemberService.Members(ctx, []string{p.Id}, []string{uid})
	if err != nil {
		log.Info("get project member error", err)
	}

	targets := make([]string, 0)
	targets = append(targets, permit_type.AnyOne.Key)
	//if p.Master == uid {
	//	targets = append(targets, permit_type.ProjectMaster.Key)
	//	targets = append(targets, permit_type.ProjectMember.Key)
	//}
	if len(members) != 0 {
		// 用户属于该项目的成员,则是用户的全局用户组身份生效
		targets = append(targets, permit_type.ProjectMember.Key)
		userGroupMembers, err := m.userGroupMemberService.Members(ctx, nil, []string{uid})
		if err != nil {
			log.Info("get user group member error", err)
		}

		for _, mb := range userGroupMembers {
			targets = append(targets, permit_type.UserGroup.KeyOf(mb.Come))
		}
	}
	//roles, err := m.projectRoleService.Roles(ctx, p.Id, uid)
	//if err == nil && roles != nil {
	//	targets = append(targets, utils.SliceToSlice(roles, func(s *project_role.ProjectRole) string {
	//		return permit_type.Role.KeyOf(s.Rid)
	//	})...)
	//}

	teamIdentity, err := m.identityTeamService.IdentifyTeam(ctx, p.Team, uid)
	if err != nil {
		return nil, err
	}
	allIdentity := make([]string, 0, len(targets)+len(teamIdentity))
	allIdentity = append(allIdentity, targets...)
	allIdentity = append(allIdentity, teamIdentity...)

	return allIdentity, nil
}
func (m *imlProjectPermitModule) Options(ctx context.Context, projectId string, keyword string) ([]*permit_dto.Option, error) {
	p, err := m.projectService.Get(ctx, projectId)
	if err != nil {
		return nil, err
	}

	result := make([]*permit_dto.Option, 0)
	result = append(result, permit_dto.SearchOptions(specialRolesProject, keyword)...)
	// 用户组
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

	if keyword != "" {
		users, _ := m.userService.Search(ctx, keyword, -1)
		if len(users) > 0 {
			// 团队成员
			userIds := utils.SliceToSlice(users, func(r *user.User) string {
				return r.UID
			})
			teamMembers, _ := m.teamMemberService.Members(ctx, []string{p.Team}, userIds)
			projectMember, _ := m.projectMemberService.Members(ctx, []string{p.Id}, userIds)
			userMap := utils.SliceToMap(users, func(r *user.User) string {
				return r.UID
			})
			if len(teamMembers)+len(projectMember) > 0 {

				userOptions := make([]*permit_dto.Option, 0, len(teamMembers)+len(projectMember))
				userOptions = append(userOptions, utils.SliceToSlice(projectMember, func(r *member.Member) *permit_dto.Option {
					us := userMap[r.UID]
					return permit_type.User.Target(us.UID, us.GetLabel())
				}, func(m *member.Member) bool {
					_, has := userMap[m.UID]
					delete(userMap, m.UID)
					return has
				})...)
				userOptions = append(userOptions, utils.SliceToSlice(teamMembers, func(r *member.Member) *permit_dto.Option {
					us := userMap[r.UID]
					return permit_type.User.Target(us.UID, us.GetLabel())
				}, func(m *member.Member) bool {
					_, has := userMap[m.UID]
					delete(userMap, m.UID)
					return has
				})...)
				result = append(result, userOptions...)
			}

		}
	}
	return result, nil
}

func (m *imlProjectPermitModule) Grant(ctx *gin.Context, projectId, access string, key string) error {

	domain, err := m.domain(ctx, projectId)
	if err != nil {
		return fmt.Errorf("project:%w", err)
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

func (m *imlProjectPermitModule) Remove(ctx *gin.Context, projectId string, access string, key string) error {

	domain, err := m.domain(ctx, projectId)
	if err != nil {
		return fmt.Errorf("project:%w", err)
	}
	return m.permitService.Remove(ctx, access, domain, key)
}

func (m *imlProjectPermitModule) ProjectAccess(ctx *gin.Context, projectId string) ([]*permit_dto.Permission, error) {

	domain, err := m.domain(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("project:%w", err)
	}

	accesses, has := access.Get(accessGroup)
	if !has {
		return nil, errors.New("no access for project")
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
func (m *imlProjectPermitModule) domain(ctx context.Context, projectId string) (string, error) {

	p, err := m.projectService.Get(ctx, projectId)
	if err != nil {
		return "", err
	}
	return fmt.Sprint("/", p.Team, "/", p.Name), nil
}
