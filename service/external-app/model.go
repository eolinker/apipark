package external_app

import (
	"time"

	external_app "github.com/eolinker/apipark/stores/external-app"
)

type ExternalApp struct {
	ID       string
	Name     string
	Token    string
	Desc     string
	Tags     []string
	Enable   bool
	Creator  string
	Updater  string
	CreateAt time.Time
	UpdateAt time.Time
}

type CreateExternalApp struct {
	Id     string
	Name   string
	Token  string
	Desc   string
	Tags   []string
	Enable bool
}

type EditExternalApp struct {
	Name   *string
	Token  *string
	Desc   *string
	Tags   []string
	Enable *bool
}

func FromEntity(ov *external_app.ExternalApp) *ExternalApp {
	return &ExternalApp{
		ID:       ov.UUID,
		Name:     ov.Name,
		Token:    ov.Token,
		Desc:     ov.Desc,
		Tags:     ov.Tags,
		Enable:   ov.Enable,
		Creator:  ov.Creator,
		Updater:  ov.Updater,
		CreateAt: ov.CreateAt,
		UpdateAt: ov.UpdateAt,
	}
}
