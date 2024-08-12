package application_authorization

import (
	"context"
	"reflect"

	"github.com/eolinker/apipark/service/universally"
	"github.com/eolinker/go-common/autowire"
)

type IAuthorizationService interface {
	universally.IServiceGet[Authorization]
	universally.IServiceDelete
	universally.IServiceCreate[Create]
	universally.IServiceEdit[Edit]
	ListByApp(ctx context.Context, appId ...string) ([]*Authorization, error)
}

func init() {
	autowire.Auto[IAuthorizationService](func() reflect.Value {
		return reflect.ValueOf(new(imlAuthorizationService))
	})
}
