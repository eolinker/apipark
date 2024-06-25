package organization

import (
	"time"

	"github.com/eolinker/apipark/stores/organization"
)

type Organization struct {
	UUID        string
	Name        string
	Description string
	Master      string
	Prefix      string
	Creator     string
	Partitions  []string
	CreateTime  time.Time
	Updater     string
	UpdateTime  time.Time
}

func fromEntity(ov *organization.Organization, pl []string) *Organization {
	return &Organization{
		UUID:        ov.UUID,
		Name:        ov.Name,
		Description: ov.Description,
		Master:      ov.Master,
		Creator:     ov.Creator,
		Updater:     ov.Updater,
		CreateTime:  ov.CreateAt,
		UpdateTime:  ov.UpdateAt,
		Partitions:  pl,
		Prefix:      ov.Prefix,
	}
}
