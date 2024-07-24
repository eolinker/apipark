package team

import (
	"context"
	"fmt"

	"github.com/eolinker/apipark/service/project"

	"github.com/eolinker/go-common/store"

	"github.com/eolinker/ap-account/service/user"

	team_member "github.com/eolinker/apipark/service/team-member"

	"github.com/google/uuid"

	team_dto "github.com/eolinker/apipark/module/team/dto"
	"github.com/eolinker/apipark/service/team"
)

var (
	_ ITeamModule = (*imlTeamModule)(nil)
)

type imlTeamModule struct {
	service        team.ITeamService              `autowired:""`
	memberService  team_member.ITeamMemberService `autowired:""`
	userService    user.IUserService              `autowired:""`
	projectService project.IProjectService        `autowired:""`
	transaction    store.ITransaction             `autowired:""`
}

func (m *imlTeamModule) GetTeam(ctx context.Context, id string) (*team_dto.Team, error) {
	tv, err := m.service.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	projects, err := m.projectService.CountTeam(ctx, id, "")
	if err != nil {
		return nil, err
	}

	return team_dto.ToTeam(tv, int(projects)), nil

}

func (m *imlTeamModule) Search(ctx context.Context, keyword string) ([]*team_dto.Item, error) {
	list, err := m.service.Search(ctx, keyword, nil)
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

func (m *imlTeamModule) Create(ctx context.Context, input *team_dto.CreateTeam) (*team_dto.Team, error) {
	if input.Id == "" {
		input.Id = uuid.New().String()
	}

	err := m.transaction.Transaction(ctx, func(ctx context.Context) error {
		err := m.service.Create(ctx, &team.CreateTeam{
			Id:          input.Id,
			Name:        input.Name,
			Description: input.Description,
			Master:      input.Master,
		})
		if err != nil {
			return err
		}
		return m.memberService.AddMemberTo(ctx, input.Id, input.Master)
	})
	if err != nil {
		return nil, err
	}
	return m.GetTeam(ctx, input.Id)
}

func (m *imlTeamModule) Edit(ctx context.Context, id string, input *team_dto.EditTeam) (*team_dto.Team, error) {
	err := m.transaction.Transaction(ctx, func(ctx context.Context) error {
		if input.Master != nil {
			// 负责人是否在团队内，若不在，则新增
			members, err := m.memberService.Members(ctx, []string{id}, []string{*input.Master})
			if err != nil {
				return err
			}
			if len(members) == 0 {
				err = m.memberService.AddMemberTo(ctx, id, *input.Master)
				if err != nil {
					return err
				}
			}
		}
		return m.service.Save(ctx, id, &team.EditTeam{
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

func (m *imlTeamModule) Delete(ctx context.Context, id string) error {
	err := m.transaction.Transaction(ctx, func(ctx context.Context) error {
		count, err := m.projectService.Count(ctx, "", map[string]interface{}{
			"team": id,
		})
		if err != nil {
			return err
		}
		if count != 0 {
			return fmt.Errorf("team has projects,cannot delete")
		}
		return m.service.Delete(ctx, id)
	})
	return err
}
