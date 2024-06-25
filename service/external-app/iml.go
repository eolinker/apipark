package external_app

import (
	"time"

	"github.com/eolinker/apipark/service/universally"
	external_app "github.com/eolinker/apipark/stores/external-app"
)

type imlExternalAppService struct {
	store external_app.IExternalAppStore `autowired:""`
	universally.IServiceGet[ExternalApp]
	universally.IServiceDelete
	universally.IServiceCreate[CreateExternalApp]
	universally.IServiceEdit[EditExternalApp]
}

func (i *imlExternalAppService) OnComplete() {
	i.IServiceGet = universally.NewGetSoftDelete[ExternalApp, external_app.ExternalApp](i.store, FromEntity)

	i.IServiceDelete = universally.NewSoftDelete[external_app.ExternalApp](i.store)

	i.IServiceCreate = universally.NewCreator[CreateExternalApp, external_app.ExternalApp](i.store, "external_app", createEntityHandler, uniquestHandler, labelHandler)

	i.IServiceEdit = universally.NewEdit[EditExternalApp, external_app.ExternalApp](i.store, updateHandler, labelHandler)
}

func labelHandler(e *external_app.ExternalApp) []string {
	return []string{e.Name, e.UUID, e.Desc}
}
func uniquestHandler(i *CreateExternalApp) []map[string]interface{} {
	return []map[string]interface{}{{"uuid": i.Id}}
}
func createEntityHandler(i *CreateExternalApp) *external_app.ExternalApp {
	now := time.Now()
	return &external_app.ExternalApp{
		UUID:     i.Id,
		Name:     i.Name,
		Token:    i.Token,
		Desc:     i.Desc,
		Tags:     i.Tags,
		Enable:   i.Enable,
		CreateAt: now,
		UpdateAt: now,
	}
}

func updateHandler(e *external_app.ExternalApp, i *EditExternalApp) {
	if i.Name != nil {
		e.Name = *i.Name
	}
	if i.Token != nil {
		e.Token = *i.Token
	}
	if i.Desc != nil {
		e.Desc = *i.Desc
	}
	if len(i.Tags) > 0 {
		e.Tags = i.Tags
	}
	if i.Enable != nil {
		e.Enable = *i.Enable
	}
	e.UpdateAt = time.Now()
}
