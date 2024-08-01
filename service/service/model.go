package service

import (
	"time"

	"github.com/eolinker/apipark/stores/service"
)

type ServiceType string

func (s ServiceType) String() string {
	return string(s)
}

func (s ServiceType) Int() int {
	switch s {
	case "inner":
		return 1
	case "public":
		return 2
	default:
		return 0
	}
}

func ToServiceType(s int) ServiceType {

	switch s {
	case 1:
		return "inner"
	case 2:
		return "public"
	default:
		return "unknown"
	}
}

type Service struct {
	Id          string
	Name        string
	Description string
	Team        string
	Prefix      string
	Logo        string
	ServiceType ServiceType
	AsServer    bool
	AsApp       bool
	CreateTime  time.Time
	UpdateTime  time.Time
}

func FromEntity(e *service.Service) *Service {
	return &Service{
		Id:          e.UUID,
		Name:        e.Name,
		Description: e.Description,
		Team:        e.Team,
		Prefix:      e.Prefix,
		Logo:        e.Logo,
		ServiceType: ToServiceType(e.ServiceType),
		AsServer:    e.AsServer,
		AsApp:       e.AsApp,
		CreateTime:  e.CreateAt,
		UpdateTime:  e.UpdateAt,
	}
}

type Create struct {
	Id          string
	Name        string
	Description string
	Team        string
	Prefix      string
	ServiceType ServiceType
	AsServer    bool
	AsApp       bool
}

type Edit struct {
	Name        *string
	Description *string
	ServiceType *ServiceType
	Logo        *string
}

type CreateTag struct {
	Tid string
	Sid string
}
