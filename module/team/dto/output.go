package team_dto

import (
	"github.com/eolinker/apipark/service/team"
	"github.com/eolinker/go-common/auto"
)

type Item struct {
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	CreateTime  auto.TimeLabel `json:"create_time"`
	UpdateTime  auto.TimeLabel `json:"update_time"`
	CanDelete   bool           `json:"can_delete"`
	ProjectNum  int64          `json:"system_num"`
}

func ToItem(model *team.Team, projectNum int64) *Item {
	return &Item{
		Id:          model.Id,
		Name:        model.Name,
		Description: model.Description,
		CreateTime:  auto.TimeLabel(model.CreateTime),
		UpdateTime:  auto.TimeLabel(model.UpdateTime),
		ProjectNum:  projectNum,
		CanDelete:   projectNum == 0,
	}
}

type Team struct {
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	CreateTime  auto.TimeLabel `json:"create_time"`
	UpdateTime  auto.TimeLabel `json:"update_time"`
	Creator     auto.Label     `json:"creator" aolabel:"user"`
	Updater     auto.Label     `json:"updater" aolabel:"user"`
	CanDelete   bool           `json:"can_delete"`
}

func ToTeam(model *team.Team, projectNum int) *Team {
	return &Team{
		Id:          model.Id,
		Name:        model.Name,
		Description: model.Description,
		CreateTime:  auto.TimeLabel(model.CreateTime),
		UpdateTime:  auto.TimeLabel(model.UpdateTime),
		Creator:     auto.UUID(model.Creator),
		Updater:     auto.UUID(model.Updater),
		CanDelete:   projectNum == 0,
	}
}
