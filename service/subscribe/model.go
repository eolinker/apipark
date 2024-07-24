package subscribe

import (
	"time"

	"github.com/eolinker/apipark/stores/subscribe"
)

type Subscribe struct {
	Id string
	// 被订阅服务相关
	Project string
	Service string

	// 订阅方相关
	Application string
	From        int
	ApplyStatus int
	CreateAt    time.Time
	//Applier     string
	//Approver    string
}

type CreateSubscribe struct {
	Uuid    string
	Service string
	Project string

	Application string
	ApplyStatus int
	From        int
}

type UpdateSubscribe struct {
	ApplyStatus *int
}

func FromEntity(e *subscribe.Subscribe) *Subscribe {
	return &Subscribe{
		Id:          e.UUID,
		Project:     e.Project,
		Service:     e.Service,
		ApplyStatus: e.ApplyStatus,
		Application: e.Application,
		From:        e.From,
		CreateAt:    e.CreateAt,
	}
}

type CreateApply struct {
	Uuid        string
	Service     string
	Project     string
	Team        string
	Application string
	ApplyTeam   string
	Reason      string
	Status      int
	Applier     string
}

type EditApply struct {
	Opinion  *string
	Status   *int
	Approver *string
}

type Apply struct {
	Id          string
	Service     string
	Project     string
	Team        string
	Application string
	ApplyTeam   string
	Applier     string
	ApplyAt     time.Time
	Approver    string
	ApproveAt   time.Time
	Status      int
	Opinion     string
	Reason      string
}

func FromApplyEntity(e *subscribe.Apply) *Apply {
	return &Apply{
		Id:          e.Uuid,
		Service:     e.Service,
		Project:     e.Project,
		Team:        e.Team,
		Application: e.Application,
		ApplyTeam:   e.ApplyTeam,
		Applier:     e.Applier,
		ApplyAt:     e.ApplyAt,
		Approver:    e.Approver,
		ApproveAt:   e.ApproveAt,
		Status:      e.Status,
		Opinion:     e.Opinion,
		Reason:      e.Reason,
	}
}
