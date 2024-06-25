package external_app

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/eolinker/go-common/auto"

	"github.com/eolinker/go-common/utils"

	external_app_dto "github.com/eolinker/apipark/module/external-app/dto"
	external_app "github.com/eolinker/apipark/service/external-app"
)

var _ IExternalAppModule = (*imlExternalAppModule)(nil)

type imlExternalAppModule struct {
	externalAppService external_app.IExternalAppService `autowired:""`
}

func (i *imlExternalAppModule) GetExternalApp(ctx context.Context, id string) (*external_app_dto.ExternalApp, error) {
	app, err := i.externalAppService.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &external_app_dto.ExternalApp{
		ID:   app.ID,
		Name: app.Name,
		Desc: app.Desc,
	}, nil

}

func (i *imlExternalAppModule) ListExternalApp(ctx context.Context) ([]*external_app_dto.ExternalAppItem, error) {
	list, err := i.externalAppService.Search(ctx, "", nil, "update_at desc")
	if err != nil {
		return nil, err
	}
	return utils.SliceToSlice(list, func(e *external_app.ExternalApp) *external_app_dto.ExternalAppItem {
		tags := strings.Join(e.Tags, ",")
		return &external_app_dto.ExternalAppItem{
			ID:         e.ID,
			Name:       e.Name,
			Token:      e.Token,
			Tags:       tags,
			Status:     e.Enable,
			Operator:   auto.UUID(e.Updater),
			UpdateTime: auto.TimeLabel(e.UpdateAt),
		}

	}), nil
}

func (i *imlExternalAppModule) CreateExternalApp(ctx context.Context, input *external_app_dto.CreateExternalApp) (*external_app_dto.ExternalApp, error) {
	if input.Id == "" {
		input.Id = uuid.New().String()
	}
	token := utils.Md5(uuid.New().String())
	err := i.externalAppService.Create(ctx, &external_app.CreateExternalApp{
		Id:     input.Id,
		Name:   input.Name,
		Token:  token,
		Desc:   input.Desc,
		Enable: true,
	})
	if err != nil {
		return nil, err
	}
	return i.GetExternalApp(ctx, input.Id)
}

func (i *imlExternalAppModule) EditExternalApp(ctx context.Context, id string, input *external_app_dto.EditExternalApp) (*external_app_dto.ExternalApp, error) {
	err := i.externalAppService.Save(ctx, id, &external_app.EditExternalApp{
		Name: input.Name,
		Desc: input.Desc,
	})
	if err != nil {
		return nil, err
	}
	return i.GetExternalApp(ctx, id)
}

func (i *imlExternalAppModule) DeleteExternalApp(ctx context.Context, id string) error {
	return i.externalAppService.Delete(ctx, id)
}

func (i *imlExternalAppModule) EnableExternalApp(ctx context.Context, id string) error {
	enable := true
	return i.externalAppService.Save(ctx, id, &external_app.EditExternalApp{
		Enable: &enable,
	})
}

func (i *imlExternalAppModule) DisableExternalApp(ctx context.Context, id string) error {
	enable := false
	return i.externalAppService.Save(ctx, id, &external_app.EditExternalApp{
		Enable: &enable,
	})
}

func (i *imlExternalAppModule) UpdateExternalAppToken(ctx context.Context, id string) (string, error) {
	token := utils.Md5(uuid.New().String())
	return token, i.externalAppService.Save(ctx, id, &external_app.EditExternalApp{
		Token: &token,
	})
}
