package subscribe

import (
	"time"

	"github.com/eolinker/apipark/stores/subscribe"
)

type Subscribe struct {
	Id string
	// 被订阅服务相关
	Project   string
	Service   string
	Partition string

	// 订阅方相关
	Application string
	From        int
	ApplyStatus int
	CreateAt    time.Time
	//Applier     string
	//Approver    string
}

type CreateSubscribe struct {
	Uuid      string
	Service   string
	Project   string
	Partition string

	Application string
	//Applier     string
	//Approver    string
	ApplyStatus int
	From        int
}

type UpdateSubscribe struct {
	ApplyStatus *int
	//Approver    *string
	//Applier     *string
}

func FromEntity(e *subscribe.Subscribe) *Subscribe {
	return &Subscribe{
		Id:          e.UUID,
		Project:     e.Project,
		Service:     e.Service,
		Partition:   e.Partition,
		ApplyStatus: e.ApplyStatus,
		Application: e.Application,
		From:        e.From,
		CreateAt:    e.CreateAt,
		//Applier:     e.Applier,
		//Approver:    e.Approver,
	}
}

type CreateApply struct {
	Uuid              string
	Service           string
	Project           string
	Team              string
	Organization      string
	ApplyPartitions   []string
	Application       string
	ApplyTeam         string
	ApplyOrganization string
	Reason            string
	Status            int
	Applier           string
}

type EditApply struct {
	ApplyPartitions []string
	Opinion         *string
	Status          *int
	Approver        *string
}

type Apply struct {
	Id                string
	Service           string
	Project           string
	Team              string
	Application       string
	ApplyTeam         string
	ApplyOrganization string
	ApplyPartitions   []string
	Partitions        []string
	Applier           string
	ApplyAt           time.Time
	Approver          string
	ApproveAt         time.Time
	Status            int
	Opinion           string
	Reason            string
}

func FromApplyEntity(e *subscribe.Apply) *Apply {
	return &Apply{
		Id:                e.Uuid,
		Service:           e.Service,
		Project:           e.Project,
		Team:              e.Team,
		Application:       e.Application,
		ApplyTeam:         e.ApplyTeam,
		ApplyOrganization: e.ApplyOrganization,
		ApplyPartitions:   e.ApplyPartitions,
		Partitions:        e.Partitions,
		Applier:           e.Applier,
		ApplyAt:           e.ApplyAt,
		Approver:          e.Approver,
		ApproveAt:         e.ApproveAt,
		Status:            e.Status,
		Opinion:           e.Opinion,
		Reason:            e.Reason,
	}
}
