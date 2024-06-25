package external_app

import (
	"reflect"

	"github.com/eolinker/go-common/autowire"
	"github.com/eolinker/go-common/store"
)

type IExternalAppStore interface {
	store.ISearchStore[ExternalApp]
}
type storeExternalApp struct {
	store.SearchStore[ExternalApp] // 用struct方式继承,会自动填充并初始化表
}

func init() {
	autowire.Auto[IExternalAppStore](func() reflect.Value {
		return reflect.ValueOf(new(storeExternalApp))
	})
}
