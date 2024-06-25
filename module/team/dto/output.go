package team_dto

import (
	"github.com/eolinker/apipark/service/team"
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
	CanDelete    bool           `json:"can_delete"`
	ProjectNum   int64          `json:"system_num"`
}

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

type Team struct {
	Id           string         `json:"id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Master       auto.Label     `json:"master" aolabel:"user"`
	CreateTime   auto.TimeLabel `json:"create_time"`
	UpdateTime   auto.TimeLabel `json:"update_time"`
	Organization auto.Label     `json:"organization" aolabel:"organization"`
	Creator      auto.Label     `json:"creator" aolabel:"user"`
	Updater      auto.Label     `json:"updater" aolabel:"user"`
	CanDelete    bool           `json:"can_delete"`
}

func ToTeam(model *team.Team, projectNum int) *Team {
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
		CanDelete:    projectNum == 0,
	}
}
