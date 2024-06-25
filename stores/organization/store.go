package organization

import (
	"github.com/eolinker/go-common/autowire"
	"github.com/eolinker/go-common/store"
	"reflect"
)

type IOrganizationStore interface {
	store.ISearchStore[Organization]
}
type IOrganizationPartitionStore interface {
	store.IBaseStore[Partition]
}

type imlOrganizationStore struct {
	store.SearchStore[Organization]
}
type imlOrganizationPartitionStore struct {
	store.Store[Partition]
}

func init() {
	autowire.Auto[IOrganizationStore](func() reflect.Value {
		return reflect.ValueOf(new(imlOrganizationStore))
	})
	autowire.Auto[IOrganizationPartitionStore](func() reflect.Value {
		return reflect.ValueOf(new(imlOrganizationPartitionStore))
	})
}
