package team_dto

import (
	"github.com/eolinker/apipark/service/team"
	team_member "github.com/eolinker/apipark/service/team-member"
	"github.com/eolinker/go-common/auto"
)

type Item struct {
	Id           string         `json:"id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Master       auto.Label     `json:"master" aolabel:"user"`
	CreateTime   auto.TimeLabel `json:"create_time"`
	UpdateTime   auto.TimeLabel `json:"update_time"`
	Organization auto.Label     `json:"organization" aolabel:"organization"`
	ProjectNum   int64          `json:"system_num"`
	CanDelete    bool           `json:"can_delete"`
}

//func ToItem(model *team.Team) *Item {
//	return &Item{
//		Id:           model.Id,
//		Name:         model.Name,
//		Description:  model.Description,
//		Master:       auto.UUID(model.Master),
//		CreateTime:   auto.TimeLabel(model.CreateTime),
//		UpdateTime:   auto.TimeLabel(model.UpdateTime),
//		Organization: auto.UUID(model.Organization),
//		ProjectNum:   model.ProjectNum,
//	}
//}

func ToItem(model *team.Team, projectNum int64) *Item {
	return &Item{
		Id:           model.Id,
		Name:         model.Name,
		Description:  model.Description,
		Master:       auto.UUID(model.Master),
		CreateTime:   auto.TimeLabel(model.CreateTime),
		UpdateTime:   auto.TimeLabel(model.UpdateTime),
		Organization: auto.UUID(model.Organization),
		ProjectNum:   projectNum,
		CanDelete:    projectNum == 0,
	}
}

type SimpleTeam struct {
	Id                  string       `json:"id"`
	Name                string       `json:"name"`
	Description         string       `json:"description"`
	Organization        auto.Label   `json:"organization" aolabel:"organization"`
	AvailablePartitions []auto.Label `json:"available_partitions" aolabel:"partition"`
	DisablePartitions   []auto.Label `json:"disable_partitions" aolabel:"partition"`
	AppNum              int64        `json:"app_num"`
}

type Team struct {
	Id                  string         `json:"id"`
	Name                string         `json:"name"`
	Description         string         `json:"description"`
	Master              auto.Label     `json:"master" aolabel:"user"`
	CreateTime          auto.TimeLabel `json:"create_time"`
	UpdateTime          auto.TimeLabel `json:"update_time"`
	Organization        auto.Label     `json:"organization" aolabel:"organization"`
	Creator             auto.Label     `json:"creator" aolabel:"user"`
	Updater             auto.Label     `json:"updater" aolabel:"user"`
	AvailablePartitions []auto.Label   `json:"available_partitions" aolabel:"partition"`
	DisablePartitions   []auto.Label   `json:"disable_partitions" aolabel:"partition"`
}

func ToTeam(model *team.Team) *Team {
	return &Team{
		Id:           model.Id,
		Name:         model.Name,
		Description:  model.Description,
		Master:       auto.UUID(model.Master),
		CreateTime:   auto.TimeLabel(model.CreateTime),
		UpdateTime:   auto.TimeLabel(model.UpdateTime),
		Organization: auto.UUID(model.Organization),
		Creator:      auto.UUID(model.Creator),
		Updater:      auto.UUID(model.Updater),
	}
}

type Member struct {
	UserId     string         `json:"user_id"`
	Name       auto.Label     `json:"name" aolabel:"user"`
	Role       string         `json:"role"`
	UserGroup  []auto.Label   `json:"user_group" aolabel:"user_group"`
	AttachTime auto.TimeLabel `json:"attach_time"`
	CanDelete  bool           `json:"can_delete"`
}

func ToMember(model *team_member.Member, masterID string, groupIDs ...string) *Member {
	role := "团队管理员"
	canDelete := false
	if model.UID != masterID {
		role = "团队成员"
		canDelete = true
	}
	return &Member{
		UserId:     model.UID,
		Name:       auto.UUID(model.UID),
		Role:       role,
		UserGroup:  auto.List(groupIDs),
		AttachTime: auto.TimeLabel(model.CreateTime),
		CanDelete:  canDelete,
	}
}

type SimpleMember struct {
	User       auto.Label   `json:"user" aolabel:"user"`
	Mail       string       `json:"mail"`
	Department []auto.Label `json:"department"`
}
