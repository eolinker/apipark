package external_app

import (
	"github.com/gin-gonic/gin"
	"reflect"

	"github.com/eolinker/go-common/autowire"

	external_app_dto "github.com/eolinker/apipark/module/external-app/dto"
)

type IExternalAppController interface {
	// GetExternalApp get external app by id
	GetExternalApp(ctx *gin.Context, id string) (*external_app_dto.ExternalApp, error)
	// ListExternalApp list external app
	ListExternalApp(ctx *gin.Context) ([]*external_app_dto.ExternalAppItem, error)
	// CreateExternalApp create external app
	CreateExternalApp(ctx *gin.Context, input *external_app_dto.CreateExternalApp) (*external_app_dto.ExternalApp, error)
	// EditExternalApp edit external app
	EditExternalApp(ctx *gin.Context, id string, input *external_app_dto.EditExternalApp) (*external_app_dto.ExternalApp, error)
	// DeleteExternalApp delete external app
	DeleteExternalApp(ctx *gin.Context, id string) error
	// EnableExternalApp enable external app
	EnableExternalApp(ctx *gin.Context, id string) error
	// DisableExternalApp disable external app
	DisableExternalApp(ctx *gin.Context, id string) error
	// UpdateExternalAppToken update external app token
	UpdateExternalAppToken(ctx *gin.Context, id string) (string, error)
}

func init() {
	autowire.Auto[IExternalAppController](func() reflect.Value {
		return reflect.ValueOf(new(imlExternalAppController))
	})
}
