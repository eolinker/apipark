package project_dto

import (
	"github.com/eolinker/apipark/service/project"
	"github.com/eolinker/go-common/auto"
)

type ProjectItem struct {
	Id           string     `json:"id"`
	Name         string     `json:"name"`
	Organization auto.Label `json:"organization" aolabel:"organization"`
	Team         auto.Label `json:"team" aolabel:"team"`
	ApiNum       int64      `json:"api_num"`
	ServiceNum   int64      `json:"service_num"`
	//SubscribeNum int64          `json:"subscribe_num"`
	Description string         `json:"description"`
	Master      auto.Label     `json:"master" aolabel:"user"`
	Partition   []auto.Label   `json:"partition,omitempty" aolabel:"partition"`
	CreateTime  auto.TimeLabel `json:"create_time"`
	UpdateTime  auto.TimeLabel `json:"update_time"`
	CanDelete   bool           `json:"can_delete"`
}

type AppItem struct {
	Id                 string         `json:"id"`
	Name               string         `json:"name"`
	Team               auto.Label     `json:"team" aolabel:"team"`
	SubscribeNum       int64          `json:"subscribe_num"`
	SubscribeVerifyNum int64          `json:"subscribe_verify_num"`
	Description        string         `json:"description"`
	CreateTime         auto.TimeLabel `json:"create_time"`
	UpdateTime         auto.TimeLabel `json:"update_time"`
	CanDelete          bool           `json:"can_delete"`
}

type SimpleProjectItem struct {
	Id           string       `json:"id"`
	Name         string       `json:"name"`
	Organization auto.Label   `json:"organization" aolabel:"organization"`
	Team         auto.Label   `json:"team" aolabel:"team"`
	Partition    []auto.Label `json:"partition,omitempty" aolabel:"partition"`
	Description  string       `json:"description"`
}

type SimpleAppItem struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Team        auto.Label `json:"team" aolabel:"team"`
	Description string     `json:"description"`
}

type Project struct {
	Id                 string         `json:"id"`
	Name               string         `json:"name"`
	Prefix             string         `json:"prefix,omitempty"`
	Description        string         `json:"description"`
	Organization       auto.Label     `json:"organization" aolabel:"organization"`
	Team               auto.Label     `json:"team" aolabel:"team"`
	Master             auto.Label     `json:"master" aolabel:"user"`
	Partition          []auto.Label   `json:"partition,omitempty" aolabel:"partition"`
	OrganizationPrefix string         `json:"organization_prefix,omitempty"`
	CreateTime         auto.TimeLabel `json:"create_time"`
	UpdateTime         auto.TimeLabel `json:"update_time"`
	AsServer           bool           `json:"as_server"`
	AsApp              bool           `json:"as_app"`
}

type App struct {
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Team        auto.Label     `json:"team" aolabel:"team"`
	CreateTime  auto.TimeLabel `json:"create_time"`
	UpdateTime  auto.TimeLabel `json:"update_time"`
	AsApp       bool           `json:"as_app"`
}

func ToProject(model *project.Project, organization string, organizationPrefix string, partitions []string) *Project {
	return &Project{
		Id:                 model.Id,
		Name:               model.Name,
		Prefix:             model.Prefix,
		Description:        model.Description,
		Organization:       auto.UUID(organization),
		Team:               auto.UUID(model.Team),
		Master:             auto.UUID(model.Master),
		Partition:          auto.List(partitions),
		OrganizationPrefix: organizationPrefix,
		CreateTime:         auto.TimeLabel(model.CreateTime),
		UpdateTime:         auto.TimeLabel(model.UpdateTime),
		AsServer:           model.AsServer,
		AsApp:              model.AsApp,
	}
}

type MemberItem struct {
	User      auto.Label   `json:"user" aolabel:"user"`
	Email     string       `json:"email"`
	Roles     []auto.Label `json:"roles" aolabel:"role"`
	CanDelete bool         `json:"can_delete"`
}

type SimpleMemberItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type TeamMemberToAdd struct {
	Id         string     `json:"id,omitempty"`
	Name       string     `json:"name,omitempty"`
	Email      string     `json:"email,omitempty"`
	Department auto.Label `json:"department" aolabel:"department"`
}
