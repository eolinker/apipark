package service

import (
	"time"

	"github.com/eolinker/apipark/stores/service"
)

const (
	StatusOn  = "on"
	StatusOff = "off"
)

type Service struct {
	Id           string
	Name         string
	Description  string
	Logo         string
	ServiceType  string
	Project      string
	Team         string
	Organization string
	Catalogue    string
	Status       string
	CreateTime   time.Time
	UpdateTime   time.Time
}

func FromEntity(e *service.Service) *Service {
	return &Service{
		Id:           e.UUID,
		Name:         e.Name,
		Description:  e.Description,
		Logo:         e.Logo,
		ServiceType:  e.ServiceType,
		Project:      e.Project,
		Team:         e.Team,
		Organization: e.Organization,
		Catalogue:    e.Catalogue,
		Status:       e.Status,
		CreateTime:   e.CreateAt,
		UpdateTime:   e.UpdateAt,
	}
}

type CreateService struct {
	Uuid         string
	Name         string
	Description  string
	Logo         string
	ServiceType  string
	Project      string
	Team         string
	Organization string
	Catalogue    string
	Status       string
	Tag          string
}

type EditService struct {
	Uuid        string
	Name        *string
	Description *string
	Logo        *string
	ServiceType *string
	Catalogue   *string
	Status      *string
	Tag         *string
}

type SearchServicePage struct {
	Keyword   string
	Page      int
	Size      int
	Catalogue []string
	Uuids     []string
}

type CreateTag struct {
	Tid string
	Sid string
}

type CreatePartition struct {
	Pid string
	Sid string
}

type Partition struct {
	Pid string
	Sid string
}

type Doc struct {
	ID         int64
	DocID      string
	Name       string
	Creator    string
	Updater    string
	Doc        string
	UpdateTime time.Time
	CreateTime time.Time
}

type SaveDoc struct {
	Sid string
	Doc string
}

type Api struct {
	Sid  string
	Aid  string
	Sort int
}

type Tag struct {
	Tid string
	Sid string
}
