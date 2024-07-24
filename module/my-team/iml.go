package my_team

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/eolinker/apipark/service/partition"

	"github.com/eolinker/apipark/service/project"

	department_member "github.com/eolinker/ap-account/service/department-member"
	"github.com/eolinker/go-common/auto"

	"github.com/eolinker/ap-account/service/user"

	"github.com/eolinker/go-common/store"

	user_group "github.com/eolinker/ap-account/service/user-group"

	team_member "github.com/eolinker/apipark/service/team-member"

	team_dto "github.com/eolinker/apipark/module/my-team/dto"
	"github.com/eolinker/apipark/service/team"
	"github.com/eolinker/go-common/utils"
)

var (
	_ ITeamModule = (*imlTeamModule)(nil)
)

type imlTeamModule struct {
	teamService             team.ITeamService                  `autowired:""`
	teamMemberService       team_member.ITeamMemberService     `autowired:""`
	userGroupMemberService  user_group.IUserGroupMemberService `autowired:""`
	userService             user.IUserService                  `autowired:""`
	departmentMemberService department_member.IMemberService   `autowired:""`
	projectService          project.IProjectService            `autowired:""`
	partitionService        partition.IPartitionService        `autowired:""`
	transaction             store.ITransaction                 `autowired:""`
}

func (m *imlTeamModule) GetTeam(ctx context.Context, id string) (*team_dto.Team, error) {
	tv, err := m.teamService.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	//availablePartitions, err := m.organizationService.Partitions(ctx, tv.Organization)
	//if err != nil {
	//	return nil, err
	//}
	globalPartitions, err := m.partitionService.List(ctx)
	if err != nil {
		return nil, err
	}
	globalPartitionMap := utils.SliceToMapO(globalPartitions, func(p *partition.Partition) (string, struct{}) {
		return p.UUID, struct{}{}
	})
	//for _, p := range availablePartitions {
	//	delete(globalPartitionMap, p)
	//}
	return &team_dto.Team{
		Id:           tv.Id,
		Name:         tv.Name,
		Description:  tv.Description,
		Master:       auto.UUID(tv.Master),
		CreateTime:   auto.TimeLabel(tv.CreateTime),
		UpdateTime:   auto.TimeLabel(tv.UpdateTime),
		Organization: auto.UUID(tv.Organization),
		Creator:      auto.UUID(tv.Creator),
		Updater:      auto.UUID(tv.Updater),
		//AvailablePartitions: auto.List(availablePartitions),
		DisablePartitions: auto.List(utils.MapToSlice(globalPartitionMap, func(k string, v struct{}) string {
			return k
		})),
	}, nil
}

func (m *imlTeamModule) Search(ctx context.Context, keyword string) ([]*team_dto.Item, error) {
	userID := utils.UserId(ctx)
	memberMap, err := m.teamMemberService.FilterMembersForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	teamIDs, ok := memberMap[userID]
	if !ok || len(teamIDs) == 0 {
		return make([]*team_dto.Item, 0), nil
	}
	list, err := m.teamService.Search(ctx, keyword, map[string]interface{}{
		"uuid": teamIDs,
	})
	if err != nil {
		return nil, err
	}
	projectNumMap, err := m.projectService.CountByTeam(ctx, keyword)
	if err != nil {
		return nil, err
	}

	outList := make([]*team_dto.Item, 0, len(list))
	for _, v := range list {
		outList = append(outList, team_dto.ToItem(v, projectNumMap[v.Id]))
	}
	return outList, nil
}

func (m *imlTeamModule) Edit(ctx context.Context, id string, input *team_dto.EditTeam) (*team_dto.Team, error) {
	err := m.transaction.Transaction(ctx, func(ctx context.Context) error {
		if input.Master != nil {
			// 负责人是否在团队内，若不在，则新增
			members, err := m.teamMemberService.Members(ctx, []string{id}, []string{*input.Master})
			if err != nil {
				return err
			}
			if len(members) == 0 {
				err = m.teamMemberService.AddMemberTo(ctx, id, *input.Master)
				if err != nil {
					return err
				}
			}
		}
		return m.teamService.Save(ctx, id, &team.EditTeam{
			Name:        input.Name,
			Description: input.Description,
			Master:      input.Master,
		})
	})

	if err != nil {
		return nil, err
	}
	return m.GetTeam(ctx, id)
}

func (m *imlTeamModule) SimpleTeams(ctx context.Context, keyword string) ([]*team_dto.SimpleTeam, error) {
	userID := utils.UserId(ctx)
	memberMap, err := m.teamMemberService.FilterMembersForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	teamIDs, ok := memberMap[userID]
	if !ok || len(teamIDs) == 0 {
		return make([]*team_dto.SimpleTeam, 0), nil
	}
	list, err := m.teamService.Search(ctx, keyword, map[string]interface{}{
		"uuid": teamIDs,
	})
	if err != nil {
		return nil, err
	}
	//partitionMap, err := m.organizationService.PartitionsByOrganization(ctx)
	//if err != nil {
	//	return nil, err
	//}
	globalPartitions, err := m.partitionService.List(ctx)
	if err != nil {
		return nil, err
	}
	apps, err := m.projectService.Search(ctx, "", map[string]interface{}{
		"team":   teamIDs,
		"as_app": true,
		"master": utils.UserId(ctx),
	})
	appCount := make(map[string]int64)
	for _, app := range apps {
		if _, ok := appCount[app.Team]; !ok {
			appCount[app.Team] = 0
		}
		appCount[app.Team]++
	}

	outList := utils.SliceToSlice(list, func(s *team.Team) *team_dto.SimpleTeam {
		globalPartitionMap := utils.SliceToMapO(globalPartitions, func(p *partition.Partition) (string, struct{}) {
			return p.UUID, struct{}{}
		})
		//availablePartitions := partitionMap[s.Organization]
		//for _, p := range availablePartitions {
		//	delete(globalPartitionMap, p)
		//}
		return &team_dto.SimpleTeam{
			Id:           s.Id,
			Name:         s.Name,
			Description:  s.Description,
			Organization: auto.UUID(s.Organization),
			//AvailablePartitions: auto.List(availablePartitions),
			DisablePartitions: auto.List(utils.MapToSlice(globalPartitionMap, func(k string, v struct{}) string {
				return k
			})),
			AppNum: appCount[s.Id],
		}
	})
	return outList, nil
}

func (m *imlTeamModule) AddMember(ctx context.Context, id string, uuids ...string) error {
	_, err := m.teamService.Get(ctx, id)
	if err != nil {
		return err
	}

	return m.teamMemberService.AddMemberTo(ctx, id, uuids...)
}

func (m *imlTeamModule) RemoveMember(ctx context.Context, id string, uuids ...string) error {
	teamInfo, err := m.teamService.Get(ctx, id)
	if err != nil {
		return err
	}
	newUuids := make([]string, 0, len(uuids))
	for _, uuid := range uuids {
		if teamInfo.Master == uuid {
			continue
		}
		newUuids = append(newUuids, uuid)
	}
	return m.teamMemberService.RemoveMemberFrom(ctx, id, newUuids...)
}

func (m *imlTeamModule) Members(ctx context.Context, id string, keyword string) ([]*team_dto.Member, error) {
	teamInfo, err := m.teamService.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	users, err := m.userService.Search(ctx, keyword, -1)
	if err != nil {
		return nil, err
	}
	userIds := utils.SliceToSlice(users, func(s *user.User) string {
		return s.UID
	})
	members, err := m.teamMemberService.Members(ctx, []string{id}, userIds)
	if err != nil {
		return nil, err
	}

	groupMemberMap, err := m.userGroupMemberService.FilterMembersForUser(ctx, userIds...)
	out := make([]*team_dto.Member, 0, len(members))
	for _, member := range members {
		gIDs, _ := groupMemberMap[member.UID]
		out = append(out, team_dto.ToMember(member, teamInfo.Master, gIDs...))
	}

	return out, nil
}

func (m *imlTeamModule) SimpleMembers(ctx context.Context, id string, keyword string) ([]*team_dto.SimpleMember, error) {
	if id == "" {
		return nil, fmt.Errorf("team id is empty")
	}
	teamInfo, err := m.teamService.Get(ctx, id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if teamInfo == nil {
		return nil, fmt.Errorf("team %s not extist", id)
	}
	users, err := m.userService.Search(ctx, keyword, -1)
	if err != nil {
		return nil, err
	}
	userMap := make(map[string]*user.User)
	userIds := make([]string, 0, len(users))
	for _, u := range users {
		userIds = append(userIds, u.UID)
		userMap[u.UID] = u
	}
	teamMembers, err := m.teamMemberService.Members(ctx, []string{id}, userIds)
	if err != nil {
		return nil, err
	}
	departmentMembers, err := m.departmentMemberService.Members(ctx, nil, userIds)
	if err != nil {
		return nil, err
	}
	departmentMemberMap := make(map[string][]string)
	for _, member := range departmentMembers {
		if _, ok := departmentMemberMap[member.UID]; !ok {
			departmentMemberMap[member.UID] = make([]string, 0)
		}
		departmentMemberMap[member.UID] = append(departmentMemberMap[member.UID], member.Come)
	}

	out := make([]*team_dto.SimpleMember, 0, len(teamMembers))
	for _, member := range teamMembers {
		u, ok := userMap[member.UID]
		if !ok {
			continue
		}

		out = append(out, &team_dto.SimpleMember{
			User:       auto.UUID(u.UID),
			Mail:       u.Email,
			Department: auto.List(departmentMemberMap[member.UID]),
		})
	}

	return out, nil
}
