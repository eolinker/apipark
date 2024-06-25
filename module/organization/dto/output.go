package organization_dto

import (
	"github.com/eolinker/apipark/service/organization"
	"github.com/eolinker/go-common/auto"
)

type Item struct {
	Id          string         `json:"id,omitempty"`                            // Id is the UUID of the organization
	Name        string         `json:"name,omitempty"`                          // Name is the name of the organization
	Description string         `json:"description,omitempty"`                   // Description is the description of the organization
	Master      auto.Label     `json:"master,omitempty" aolabel:"user"`         // Master is the UUID of the organization's master partition
	Partition   []auto.Label   `json:"partition,omitempty" aolabel:"partition"` // Partition is the list of the organization's partition UUID
	Prefix      string         `json:"prefix,omitempty"`                        // Prefix is the prefix of the organization's UUID
	CreateTime  auto.TimeLabel `json:"create_time,omitempty"`                   // CreateTime is the time when the organization was created
	UpdateTime  auto.TimeLabel `json:"update_time,omitempty"`                   // UpdateTime is the time when the organization was updated
	Updater     auto.Label     `json:"updater,omitempty" aolabel:"user"`        //
	Creator     auto.Label     `json:"creator,omitempty" aolabel:"user"`
	CanDelete   bool           `json:"can_delete"`
}

type Detail struct {
	Id          string         `json:"id"`                             //
	Name        string         `json:"name"`                           //
	Description string         `json:"description"`                    //
	Master      auto.Label     `json:"master" aolabel:"user"`          //
	Prefix      string         `json:"prefix"`                         //
	Updater     auto.Label     `json:"updater" aolabel:"user"`         //
	Creator     auto.Label     `json:"creator" aolabel:"user"`         //
	CreateTime  auto.TimeLabel `json:"create_time"`                    //
	UpdateTime  auto.TimeLabel `json:"update_time"`                    //
	Partitions  []auto.Label   `json:"partitions" aolabel:"partition"` //
}

func NewDetail(m *organization.Organization) *Detail {
	return &Detail{
		Id:          m.UUID,
		Name:        m.Name,
		Description: m.Description,
		Master:      auto.UUID(m.Master),
		Prefix:      m.Prefix,
		Updater:     auto.UUID(m.Updater),
		Creator:     auto.UUID(m.Creator),
		CreateTime:  auto.TimeLabel(m.CreateTime),
		UpdateTime:  auto.TimeLabel(m.UpdateTime),
		Partitions:  auto.List(m.Partitions),
	}
}

type Partition struct {
	Id   string `json:"id,omitempty"`   //
	Name string `json:"name,omitempty"` //
}

type Simple struct {
	Id   string `json:"id,omitempty"`   //
	Name string `json:"name,omitempty"` //
}
