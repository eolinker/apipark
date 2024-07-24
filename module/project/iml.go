package project

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/eolinker/apipark/service/partition"

	"github.com/eolinker/apipark/service/subscribe"
	"gorm.io/gorm"

	"github.com/eolinker/ap-account/service/member"

	"github.com/eolinker/apipark/service/api"

	"github.com/eolinker/apipark/service/service"

	project_role "github.com/eolinker/apipark/service/project-role"

	"github.com/eolinker/ap-account/service/user"

	"github.com/eolinker/go-common/auto"

	team_member "github.com/eolinker/apipark/service/team-member"

	project_member "github.com/eolinker/apipark/service/project-member"

	"github.com/eolinker/go-common/store"

	"github.com/google/uuid"

	"github.com/eolinker/go-common/utils"

	"github.com/eolinker/apipark/service/team"

	"github.com/eolinker/apipark/service/project"

	project_dto "github.com/eolinker/apipark/module/project/dto"
)

var (
	_ IProjectModule = (*imlProjectModule)(nil)
)

type imlProjectModule struct {
	partitionService     partition.IPartitionService    `autowired:""`
	projectService       project.IProjectService        `autowired:""`
	projectMemberService project_member.IMemberService  `autowired:""`
	teamService          team.ITeamService              `autowired:""`
	teamMemberService    team_member.ITeamMemberService `autowired:""`
	serviceService       service.IServiceService        `autowired:""`
	apiService           api.IAPIService                `autowired:""`
	transaction          store.ITransaction             `autowired:""`
}

func (i *imlProjectModule) searchMyProjects(ctx context.Context, teamId string, keyword string) ([]*project.Project, error) {

	userID := utils.UserId(ctx)
	condition := make(map[string]interface{})
	condition["as_server"] = true
	if teamId != "" {
		_, err := i.teamService.Get(ctx, teamId)
		if err != nil {
			return nil, err
		}
		condition["team"] = teamId
		return i.projectService.Search(ctx, keyword, condition, "update_at desc")
	} else {
		membersForUser, err := i.teamMemberService.FilterMembersForUser(ctx, userID)
		if err != nil {
			return nil, err
		}
		teamIds := membersForUser[userID]
		condition["team"] = teamIds
		return i.projectService.Search(ctx, keyword, condition, "update_at desc")
	}

}

func (i *imlProjectModule) SearchMyProjects(ctx context.Context, teamId string, keyword string) ([]*project_dto.ProjectItem, error) {
	projects, err := i.searchMyProjects(ctx, teamId, keyword)
	if err != nil {
		return nil, err
	}
	projectIDs := utils.SliceToSlice(projects, func(p *project.Project) string {
		return p.Id
	})
	apiCountMap, err := i.apiService.CountByGroup(ctx, "", map[string]interface{}{"project": projectIDs}, "project")
	if err != nil {
		return nil, err
	}
	serviceCountMap, err := i.serviceService.CountByGroup(ctx, "", map[string]interface{}{"project": projectIDs}, "project")
	if err != nil {
		return nil, err
	}

	items := make([]*project_dto.ProjectItem, 0, len(projects))
	for _, model := range projects {
		if teamId != "" && model.Team != teamId {
			continue
		}
		apiCount := apiCountMap[model.Id]
		serviceCount := serviceCountMap[model.Id]
		items = append(items, &project_dto.ProjectItem{
			Id:          model.Id,
			Name:        model.Name,
			Description: model.Description,
			Master:      auto.UUID(model.Master),
			CreateTime:  auto.TimeLabel(model.CreateTime),
			UpdateTime:  auto.TimeLabel(model.UpdateTime),
			Team:        auto.UUID(model.Team),
			ApiNum:      apiCount,
			ServiceNum:  serviceCount,
			CanDelete:   apiCount == 0 && serviceCount == 0,
		})
	}
	return items, nil
}

func (i *imlProjectModule) SimpleAPPS(ctx context.Context, keyword string) ([]*project_dto.SimpleProjectItem, error) {
	w := make(map[string]interface{})
	w["as_app"] = true
	projects, err := i.projectService.Search(ctx, keyword, w)
	if err != nil {
		return nil, err
	}
	return utils.SliceToSlice(projects, func(p *project.Project) *project_dto.SimpleProjectItem {
		return &project_dto.SimpleProjectItem{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,

			Team: auto.UUID(p.Team),
		}
	}), nil
}

func (i *imlProjectModule) SimpleProjects(ctx context.Context, keyword string) ([]*project_dto.SimpleProjectItem, error) {
	w := make(map[string]interface{})
	w["as_server"] = true
	//if partition != "" {
	//	pp, err := i.projectPartitionService.ListByPartition(ctx, partition)
	//	if err != nil {
	//		return nil, err
	//	}
	//	w["uuid"] = utils.SliceToSlice(pp, func(p *project.Partition) string {
	//		return p.Project
	//	})
	//}
	projects, err := i.projectService.Search(ctx, keyword, w)
	if err != nil {
		return nil, err
	}

	items := make([]*project_dto.SimpleProjectItem, 0, len(projects))
	for _, p := range projects {

		items = append(items, &project_dto.SimpleProjectItem{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Team:        auto.UUID(p.Team),
		})
	}
	return items, nil
}

func (i *imlProjectModule) MySimpleProjects(ctx context.Context, keyword string) ([]*project_dto.SimpleProjectItem, error) {
	projects, err := i.searchMyProjects(ctx, "", keyword)

	if err != nil {
		return nil, err
	}

	items := make([]*project_dto.SimpleProjectItem, 0, len(projects))
	for _, p := range projects {

		items = append(items, &project_dto.SimpleProjectItem{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Team:        auto.UUID(p.Team),
		})
	}
	return items, nil
}

func (i *imlProjectModule) GetProject(ctx context.Context, id string) (*project_dto.Project, error) {
	projectInfo, err := i.projectService.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return project_dto.ToProject(projectInfo), nil
}

func (i *imlProjectModule) Search(ctx context.Context, teamID string, keyword string) ([]*project_dto.ProjectItem, error) {
	var list []*project.Project
	var err error
	if teamID != "" {
		_, err = i.teamService.Get(ctx, teamID)
		if err != nil {
			return nil, err
		}
		list, err = i.projectService.Search(ctx, keyword, map[string]interface{}{"team": teamID}, "update_at desc")
	} else {
		list, err = i.projectService.Search(ctx, keyword, nil, "update_at desc")
	}
	if err != nil {
		return nil, err
	}

	projectIds := utils.SliceToSlice(list, func(s *project.Project) string {
		return s.Id
	})

	apiCountMap, err := i.apiService.CountByGroup(ctx, "", map[string]interface{}{"project": projectIds}, "project")
	if err != nil {
		return nil, err
	}
	serviceCountMap, err := i.serviceService.CountByGroup(ctx, "", map[string]interface{}{"project": projectIds}, "project")
	if err != nil {
		return nil, err
	}

	items := make([]*project_dto.ProjectItem, 0, len(list))
	for _, model := range list {
		apiCount := apiCountMap[model.Id]
		serviceCount := serviceCountMap[model.Id]
		items = append(items, &project_dto.ProjectItem{
			Id:          model.Id,
			Name:        model.Name,
			Description: model.Description,
			Master:      auto.UUID(model.Master),
			CreateTime:  auto.TimeLabel(model.CreateTime),
			UpdateTime:  auto.TimeLabel(model.UpdateTime),
			Team:        auto.UUID(model.Team),
			ApiNum:      apiCount,
			ServiceNum:  serviceCount,
			CanDelete:   apiCount == 0 && serviceCount == 0,
		})
	}
	return items, nil
}

func (i *imlProjectModule) CreateProject(ctx context.Context, teamID string, input *project_dto.CreateProject) (*project_dto.Project, error) {

	if input.Id == "" {
		input.Id = uuid.New().String()
	}
	mo := &project.CreateProject{
		Id:          input.Id,
		Name:        input.Name,
		Description: input.Description,
		Master:      input.Master,
		Team:        teamID,
		Prefix:      input.Prefix,
	}
	if input.AsApp == nil {
		// 默认值为false
		mo.AsApp = false
	} else {
		mo.AsApp = *input.AsApp
	}
	if input.AsServer == nil {
		// 默认值为true
		mo.AsServer = true
	} else {
		mo.AsServer = *input.AsServer
	}
	input.Prefix = strings.Trim(strings.Trim(input.Prefix, " "), "/")
	err := i.transaction.Transaction(ctx, func(ctx context.Context) error {

		// 判断用户是否在团队内
		members, err := i.teamMemberService.Members(ctx, []string{teamID}, []string{input.Master})
		if err != nil {
			return err
		}
		if len(members) == 0 {
			return fmt.Errorf("master is not in team")
		}

		err = i.projectService.Create(ctx, mo)
		if err != nil {
			return err
		}

		return i.projectMemberService.AddMemberTo(ctx, input.Id, input.Master)
	})
	if err != nil {
		return nil, err
	}
	return i.GetProject(ctx, input.Id)
}

func (i *imlProjectModule) EditProject(ctx context.Context, id string, input *project_dto.EditProject) (*project_dto.Project, error) {
	_, err := i.projectService.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	err = i.transaction.Transaction(ctx, func(ctx context.Context) error {

		if input.Master != nil {
			projectInfo, err := i.projectService.Get(ctx, id)
			if err != nil {
				return err
			}

			// 判断用户是否在团队内
			members, err := i.teamMemberService.Members(ctx, []string{projectInfo.Team}, []string{*input.Master})
			if err != nil {
				return err
			}
			if len(members) == 0 {
				return fmt.Errorf("master is not in team")
			}
			// 负责人是否在项目内，若不在，则新增
			projectMembers, err := i.projectMemberService.Members(ctx, []string{id}, []string{*input.Master})
			if err != nil {
				return err
			}
			if len(projectMembers) == 0 {
				err = i.projectMemberService.AddMemberTo(ctx, id, *input.Master)
				if err != nil {
					return err
				}
			}
		}
		return i.projectService.Save(ctx, id, &project.EditProject{
			Name:        input.Name,
			Description: input.Description,
			Master:      input.Master,
		})
	})

	if err != nil {
		return nil, err
	}
	return i.GetProject(ctx, id)
}

func (i *imlProjectModule) DeleteProject(ctx context.Context, id string) error {

	err := i.transaction.Transaction(ctx, func(ctx context.Context) error {
		count, err := i.apiService.CountByProject(ctx, id)
		if err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("project has apis, can not delete")
		}
		// 删除项目成员
		err = i.projectMemberService.Delete(ctx, id)
		if err != nil {
			return err
		}

		return i.projectService.Delete(ctx, id)
	})
	return err
}

var _ IProjectMemberModule = (*imlProjectMemberModule)(nil)

type imlProjectMemberModule struct {
	projectService       project.IProjectService          `autowired:""`
	projectMemberService project_member.IMemberService    `autowired:""`
	projectRoleService   project_role.IProjectRoleService `autowired:""`
	userService          user.IUserService                `autowired:""`
	transaction          store.ITransaction               `autowired:""`
	teamMemberService    team_member.ITeamMemberService   `autowired:""`
}

func (i *imlProjectMemberModule) SimpleMembersToAdd(ctx context.Context, pid string, keyword string) ([]*project_dto.TeamMemberToAdd, error) {

	pro, err := i.projectService.Get(ctx, pid)
	if err != nil {
		return nil, err
	}

	users, err := i.userService.Search(ctx, keyword, -1)
	if err != nil {
		return nil, err
	}
	userIds := utils.SliceToSlice(users, func(s *user.User) string {
		return s.UID
	})

	members, err := i.teamMemberService.Members(ctx, []string{pro.Team}, userIds)
	if err != nil {
		return nil, err
	}
	userMaps := utils.SliceToMap(users, func(s *user.User) string {
		return s.UID
	})

	return utils.SliceToSlice(members, func(s *member.Member) *project_dto.TeamMemberToAdd {
		uf := userMaps[s.UID]
		return &project_dto.TeamMemberToAdd{
			Id:    uf.UID,
			Name:  uf.Username,
			Email: uf.Email,
			//Department: uf.,
		}
	}, func(m *member.Member) bool {
		_, has := userMaps[m.UID]
		return has
	}), nil

}

func (i *imlProjectMemberModule) SimpleMembers(ctx context.Context, pid string) ([]*project_dto.SimpleMemberItem, error) {
	members, err := i.projectMemberService.Members(ctx, []string{pid}, nil)
	if err != nil {
		return nil, err
	}
	userIds := utils.SliceToSlice(members, func(s *member.Member) string {
		return s.UID
	})
	users, err := i.userService.Get(ctx, userIds...)
	if err != nil {
		return nil, err
	}
	out := utils.SliceToSlice(users, func(s *user.User) *project_dto.SimpleMemberItem {
		return &project_dto.SimpleMemberItem{
			Id:   s.UID,
			Name: s.Username,
		}
	})
	return out, nil
}

func (i *imlProjectMemberModule) Members(ctx context.Context, id string, keyword string) ([]*project_dto.MemberItem, error) {
	pInfo, err := i.projectService.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	users, err := i.userService.Search(ctx, keyword, -1)
	if err != nil {
		return nil, err
	}
	userMap := make(map[string]*user.User)
	userIds := make([]string, 0, len(users))
	for _, u := range users {
		userIds = append(userIds, u.UID)
		userMap[u.UID] = u
	}
	members, err := i.projectMemberService.Members(ctx, []string{id}, userIds)
	if err != nil {
		return nil, err
	}
	roleMap, err := i.projectRoleService.RoleMap(ctx, id)
	if err != nil {
		return nil, err
	}

	out := utils.SliceToSlice(members, func(info *project_member.Member) *project_dto.MemberItem {
		roleIDs := make([]string, 0, len(roleMap[info.UID]))
		for _, r := range roleMap[info.UID] {
			roleIDs = append(roleIDs, r.Rid)
		}
		item := &project_dto.MemberItem{
			User:  auto.UUID(info.UID),
			Email: "",
			Roles: auto.List(roleIDs),
		}

		u, ok := userMap[info.UID]
		if ok {
			item.Email = u.Email
		}
		if pInfo.Master == info.UID {
			item.CanDelete = false
		} else {
			item.CanDelete = true
		}
		return item

	})
	return out, nil
}

func (i *imlProjectMemberModule) AddMember(ctx context.Context, pid string, userIDs []string) error {
	_, err := i.projectService.Get(ctx, pid)
	if err != nil {
		return err
	}
	return i.projectMemberService.AddMemberTo(ctx, pid, userIDs...)
}

func (i *imlProjectMemberModule) RemoveMember(ctx context.Context, pid string, userIDs []string) error {
	// 删除的成员是否有项目负责人
	info, err := i.projectService.Get(ctx, pid)
	if err != nil {
		return err
	}
	uids := make([]string, 0, len(userIDs))
	for _, id := range userIDs {
		if id != info.Master {
			uids = append(uids, id)
		}
	}
	if len(uids) == 0 {
		return nil
	}
	return i.projectMemberService.RemoveMemberFrom(ctx, pid, userIDs...)
}

func (i *imlProjectMemberModule) EditProjectMember(ctx context.Context, pid string, uid string, roles []string) error {
	return i.transaction.Transaction(ctx, func(ctx context.Context) error {
		_, err := i.projectService.Get(ctx, pid)
		if err != nil {
			return err
		}
		_, err = i.userService.Get(ctx, uid)
		if err != nil {
			return err
		}

		err = i.projectRoleService.DeleteRole(ctx, pid, uid)
		if err != nil {
			return err
		}
		for _, rid := range roles {
			err = i.projectRoleService.AddRole(ctx, pid, uid, rid)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

var _ IAppModule = &imlAppModule{}

type imlAppModule struct {
	teamService          team.ITeamService              `autowired:""`
	projectService       project.IProjectService        `autowired:""`
	projectMemberService project_member.IMemberService  `autowired:""`
	teamMemberService    team_member.ITeamMemberService `autowired:""`
	subscribeService     subscribe.ISubscribeService    `autowired:""`
	transaction          store.ITransaction             `autowired:""`
}

func (i *imlAppModule) CreateApp(ctx context.Context, teamID string, input *project_dto.CreateApp) (*project_dto.App, error) {
	//teamInfo, err := i.teamService.Get(ctx, teamID)
	//if err != nil {
	//	return nil, err
	//}
	if input.Id == "" {
		input.Id = uuid.New().String()
	}
	userId := utils.UserId(ctx)
	mo := &project.CreateProject{
		Id:          input.Id,
		Name:        input.Name,
		Description: input.Description,
		Master:      userId,
		Team:        teamID,
		AsApp:       true,
	}
	// 判断用户是否在团队内
	members, err := i.teamMemberService.Members(ctx, []string{teamID}, []string{userId})
	if err != nil {
		return nil, err
	}
	if len(members) == 0 {
		return nil, fmt.Errorf("master is not in team")
	}

	err = i.transaction.Transaction(ctx, func(ctx context.Context) error {

		err = i.projectService.Create(ctx, mo)
		if err != nil {
			return err
		}

		return i.projectMemberService.AddMemberTo(ctx, input.Id, userId)
	})
	if err != nil {
		return nil, err
	}
	return i.GetApp(ctx, input.Id)
}

func (i *imlAppModule) UpdateApp(ctx context.Context, appId string, input *project_dto.UpdateApp) (*project_dto.App, error) {
	userId := utils.UserId(ctx)
	info, err := i.projectService.Get(ctx, appId)
	if err != nil {
		return nil, err
	}
	if !info.AsApp {
		return nil, fmt.Errorf("not app")
	}
	if info.Master != userId {
		return nil, fmt.Errorf("user is not app master, can not update")
	}

	err = i.projectService.Save(ctx, appId, &project.EditProject{
		Name:        input.Name,
		Description: input.Description,
	})
	if err != nil {
		return nil, err
	}
	return i.GetApp(ctx, info.Id)
}

func (i *imlAppModule) searchMyApps(ctx context.Context, teamId string, keyword string) ([]*project.Project, error) {
	userID := utils.UserId(ctx)
	condition := make(map[string]interface{})
	condition["as_app"] = true
	condition["master"] = userID
	if teamId != "" {
		_, err := i.teamService.Get(ctx, teamId)
		if err != nil {
			return nil, err
		}
		condition["team"] = teamId
		return i.projectService.Search(ctx, keyword, condition, "update_at desc")
	} else {
		membersForUser, err := i.teamMemberService.FilterMembersForUser(ctx, userID)
		if err != nil {
			return nil, err
		}
		teamIds := membersForUser[userID]
		condition["team"] = teamIds

		return i.projectService.Search(ctx, keyword, condition, "update_at desc")
	}
}

func (i *imlAppModule) SearchMyApps(ctx context.Context, teamId string, keyword string) ([]*project_dto.AppItem, error) {
	projects, err := i.searchMyApps(ctx, teamId, keyword)
	if err != nil {
		return nil, err
	}
	projectIds := utils.SliceToSlice(projects, func(p *project.Project) string {
		return p.Id
	})

	subscribers, err := i.subscribeService.SubscriptionsByApplication(ctx, projectIds...)
	if err != nil {
		return nil, err
	}

	subscribeCount := map[string]int64{}
	subscribeVerifyCount := map[string]int64{}
	verifyTmp := map[string]struct{}{}
	subscribeTmp := map[string]struct{}{}
	for _, s := range subscribers {
		key := fmt.Sprintf("%s-%s", s.Service, s.Application)
		switch s.ApplyStatus {
		case subscribe.ApplyStatusSubscribe:
			if _, ok := subscribeTmp[key]; !ok {
				subscribeTmp[key] = struct{}{}
				subscribeCount[s.Application]++
			}
		case subscribe.ApplyStatusReview:
			if _, ok := verifyTmp[key]; !ok {
				verifyTmp[key] = struct{}{}
				subscribeVerifyCount[s.Application]++
			}
		default:

		}
	}
	items := make([]*project_dto.AppItem, 0, len(projects))
	for _, model := range projects {
		subscribeNum := subscribeCount[model.Id]
		verifyNum := subscribeVerifyCount[model.Id]
		items = append(items, &project_dto.AppItem{
			Id:                 model.Id,
			Name:               model.Name,
			Description:        model.Description,
			CreateTime:         auto.TimeLabel(model.CreateTime),
			UpdateTime:         auto.TimeLabel(model.UpdateTime),
			Team:               auto.UUID(model.Team),
			SubscribeNum:       subscribeNum,
			SubscribeVerifyNum: verifyNum,
			CanDelete:          subscribeNum == 0,
		})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].SubscribeNum != items[j].SubscribeNum {
			return items[i].SubscribeNum > items[j].SubscribeNum
		}
		if items[i].SubscribeVerifyNum != items[j].SubscribeVerifyNum {
			return items[i].SubscribeVerifyNum > items[j].SubscribeVerifyNum
		}
		return items[i].Name < items[j].Name
	})
	return items, nil
}

func (i *imlAppModule) SimpleApps(ctx context.Context, keyword string) ([]*project_dto.SimpleAppItem, error) {
	w := make(map[string]interface{})
	w["as_app"] = true
	projects, err := i.projectService.Search(ctx, keyword, w)
	if err != nil {
		return nil, err
	}
	return utils.SliceToSlice(projects, func(p *project.Project) *project_dto.SimpleAppItem {
		return &project_dto.SimpleAppItem{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Team:        auto.UUID(p.Team),
		}
	}), nil
}

func (i *imlAppModule) MySimpleApps(ctx context.Context, keyword string) ([]*project_dto.SimpleAppItem, error) {
	projects, err := i.searchMyApps(ctx, "", keyword)
	if err != nil {
		return nil, err
	}
	items := make([]*project_dto.SimpleAppItem, 0, len(projects))
	for _, p := range projects {

		items = append(items, &project_dto.SimpleAppItem{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Team:        auto.UUID(p.Team),
		})
	}
	return items, nil
}

func (i *imlAppModule) GetApp(ctx context.Context, appId string) (*project_dto.App, error) {
	info, err := i.projectService.Get(ctx, appId)
	if err != nil {
		return nil, err
	}
	if !info.AsApp {
		return nil, errors.New("not app")
	}
	return &project_dto.App{
		Id:          info.Id,
		Name:        info.Name,
		Description: info.Description,
		Team:        auto.UUID(info.Team),
		CreateTime:  auto.TimeLabel(info.CreateTime),
		UpdateTime:  auto.TimeLabel(info.UpdateTime),
		AsApp:       info.AsApp,
	}, nil
}

func (i *imlAppModule) DeleteApp(ctx context.Context, appId string) error {
	info, err := i.projectService.Get(ctx, appId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return nil
	}
	if !info.AsApp {
		return errors.New("not app, can not delete")
	}
	if info.Master != utils.UserId(ctx) {
		return errors.New("not master, can not delete")
	}
	return i.projectService.Delete(ctx, appId)
}
