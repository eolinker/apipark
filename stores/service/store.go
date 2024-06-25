package service

import (
	"reflect"

	"github.com/eolinker/go-common/autowire"
	"github.com/eolinker/go-common/store"
)

type IServiceStore interface {
	store.ISearchStore[Service]
}

type imlServiceStore struct {
	store.SearchStoreSoftDelete[Service]
}

type IServicePartitionStore interface {
	store.IBaseStore[Partition]
}

type imlServicePartitionStore struct {
	store.Store[Partition]
}

type IServiceTagStore interface {
	store.IBaseStore[Tag]
}

type imlServiceTagStore struct {
	store.Store[Tag]
}

type IServiceDocStore interface {
	store.ISearchStore[Doc]
}

type imlServiceDocStore struct {
	store.SearchStore[Doc]
}

type IServiceApiStore interface {
	store.IBaseStore[Api]
}

type imlServiceApiStore struct {
	store.Store[Api]
}

func init() {
	autowire.Auto[IServiceStore](func() reflect.Value {
		return reflect.ValueOf(new(imlServiceStore))
	})

	autowire.Auto[IServicePartitionStore](func() reflect.Value {
		return reflect.ValueOf(new(imlServicePartitionStore))
	})

	autowire.Auto[IServiceTagStore](func() reflect.Value {
		return reflect.ValueOf(new(imlServiceTagStore))
	})

	autowire.Auto[IServiceDocStore](func() reflect.Value {
		return reflect.ValueOf(new(imlServiceDocStore))
	})

	autowire.Auto[IServiceApiStore](func() reflect.Value {
		return reflect.ValueOf(new(imlServiceApiStore))
	})
}
