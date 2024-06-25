package external_app

import (
	external_app "github.com/eolinker/apipark/module/external-app"
	external_app_dto "github.com/eolinker/apipark/module/external-app/dto"
	"github.com/gin-gonic/gin"
)

var _ IExternalAppController = (*imlExternalAppController)(nil)

type imlExternalAppController struct {
	module external_app.IExternalAppModule `autowired:""`
}

func (i *imlExternalAppController) GetExternalApp(ctx *gin.Context, id string) (*external_app_dto.ExternalApp, error) {
	return i.module.GetExternalApp(ctx, id)
}

func (i *imlExternalAppController) ListExternalApp(ctx *gin.Context) ([]*external_app_dto.ExternalAppItem, error) {
	return i.module.ListExternalApp(ctx)
}

func (i *imlExternalAppController) CreateExternalApp(ctx *gin.Context, input *external_app_dto.CreateExternalApp) (*external_app_dto.ExternalApp, error) {
	return i.module.CreateExternalApp(ctx, input)
}

func (i *imlExternalAppController) EditExternalApp(ctx *gin.Context, id string, input *external_app_dto.EditExternalApp) (*external_app_dto.ExternalApp, error) {
	return i.module.EditExternalApp(ctx, id, input)
}

func (i *imlExternalAppController) DeleteExternalApp(ctx *gin.Context, id string) error {
	return i.module.DeleteExternalApp(ctx, id)
}

func (i *imlExternalAppController) EnableExternalApp(ctx *gin.Context, id string) error {
	return i.module.EnableExternalApp(ctx, id)
}

func (i *imlExternalAppController) DisableExternalApp(ctx *gin.Context, id string) error {
	return i.module.DisableExternalApp(ctx, id)
}

func (i *imlExternalAppController) UpdateExternalAppToken(ctx *gin.Context, id string) (string, error) {
	return i.module.UpdateExternalAppToken(ctx, id)
}
