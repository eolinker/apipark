package external_app

import (
	"context"
	"reflect"

	"github.com/eolinker/go-common/autowire"

	external_app_dto "github.com/eolinker/apipark/module/external-app/dto"
)

type IExternalAppModule interface {
	// GetExternalApp get external app by id
	GetExternalApp(ctx context.Context, id string) (*external_app_dto.ExternalApp, error)
	// ListExternalApp list external app
	ListExternalApp(ctx context.Context) ([]*external_app_dto.ExternalAppItem, error)
	// CreateExternalApp create external app
	CreateExternalApp(ctx context.Context, input *external_app_dto.CreateExternalApp) (*external_app_dto.ExternalApp, error)
	// EditExternalApp edit external app
	EditExternalApp(ctx context.Context, id string, input *external_app_dto.EditExternalApp) (*external_app_dto.ExternalApp, error)
	// DeleteExternalApp delete external app
	DeleteExternalApp(ctx context.Context, id string) error
	// EnableExternalApp enable external app
	EnableExternalApp(ctx context.Context, id string) error
	// DisableExternalApp disable external app
	DisableExternalApp(ctx context.Context, id string) error
	// UpdateExternalAppToken update external app token
	UpdateExternalAppToken(ctx context.Context, id string) (string, error)
}

func init() {
	autowire.Auto[IExternalAppModule](func() reflect.Value {
		return reflect.ValueOf(new(imlExternalAppModule))
	})
}
