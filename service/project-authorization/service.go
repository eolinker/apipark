package project_authorization

import (
	"context"
	"reflect"

	"github.com/eolinker/apipark/service/universally"
	"github.com/eolinker/go-common/autowire"
)

type IProjectAuthorizationService interface {
	universally.IServiceGet[Authorization]
	universally.IServiceDelete
	universally.IServiceCreate[CreateAuthorization]
	universally.IServiceEdit[EditAuthorization]
	ListByProject(ctx context.Context, pid ...string) ([]*Authorization, error)
}

func init() {
	autowire.Auto[IProjectAuthorizationService](func() reflect.Value {
		return reflect.ValueOf(new(imlProjectAuthorizationService))
	})
}
