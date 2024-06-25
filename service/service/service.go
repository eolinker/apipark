package service

import (
	"context"
	"reflect"

	"github.com/eolinker/apipark/service/universally"
	"github.com/eolinker/go-common/autowire"
)

type IServiceService interface {
	universally.IServiceGet[Service]
	universally.IServiceDelete
	universally.IServiceCreate[CreateService]
	universally.IServiceEdit[EditService]
	SearchServicePage(ctx context.Context, input SearchServicePage, sort ...string) ([]*Service, int64, error)
	ListByProject(ctx context.Context, pid ...string) ([]*Service, error)
	SearchByUuids(ctx context.Context, keyword string, uuids ...string) ([]*Service, error)
}

type IPartitionsService interface {
	Delete(ctx context.Context, pids []string, sids []string) error
	Create(ctx context.Context, input *CreatePartition) error
	List(ctx context.Context, sid string) ([]*Partition, error)
	PartitionsByService(ctx context.Context, sids ...string) (map[string][]string, error)
}

type ITagService interface {
	Delete(ctx context.Context, tids []string, sids []string) error
	Create(ctx context.Context, input *CreateTag) error
	List(ctx context.Context, sids []string, tids []string) ([]*Tag, error)
}

type IDocService interface {
	Get(ctx context.Context, sid string) (*Doc, error)
	Save(ctx context.Context, input *SaveDoc) error
}

type IApiService interface {
	Bind(ctx context.Context, sid string, aid string, sort int) error
	Unbind(ctx context.Context, sid string, aid string) error
	List(ctx context.Context, sids ...string) ([]*Api, error)
	Count(ctx context.Context, sid string) (int64, error)
	CountBySids(ctx context.Context, sids ...string) (map[string]int64, error)
	Clear(ctx context.Context, sid string) error
	LastSortIndex(ctx context.Context, sid string) (int, error)
}

func init() {
	autowire.Auto[IServiceService](func() reflect.Value {
		return reflect.ValueOf(new(imlServiceService))
	})

	autowire.Auto[IPartitionsService](func() reflect.Value {
		return reflect.ValueOf(new(imlPartitionsService))
	})

	autowire.Auto[ITagService](func() reflect.Value {
		return reflect.ValueOf(new(imlTagService))
	})

	autowire.Auto[IDocService](func() reflect.Value {
		return reflect.ValueOf(new(imlDocService))
	})

	autowire.Auto[IApiService](func() reflect.Value {
		return reflect.ValueOf(new(imlApiService))
	})
}
