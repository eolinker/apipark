package external_app

import (
	"reflect"

	"github.com/eolinker/apipark/service/universally"
	"github.com/eolinker/go-common/autowire"
)

type IExternalAppService interface {
	universally.IServiceGet[ExternalApp]
	universally.IServiceDelete
	universally.IServiceCreate[CreateExternalApp]
	universally.IServiceEdit[EditExternalApp]
}

func init() {
	autowire.Auto[IExternalAppService](func() reflect.Value {
		return reflect.ValueOf(new(imlExternalAppService))
	})
}
