package publish

import (
	"context"
	"reflect"

	"github.com/eolinker/apipark/gateway"

	"github.com/eolinker/apipark/module/publish/dto"
	"github.com/eolinker/go-common/autowire"
)

type IPublishModule interface {
	CheckPublish(ctx context.Context, project string, releaseId string) (*dto.DiffOut, error)

	Apply(ctx context.Context, project string, input *dto.ApplyInput) (*dto.Publish, error)
	Stop(ctx context.Context, project string, id string) error
	Refuse(ctx context.Context, project string, id string, commits string) error
	Accept(ctx context.Context, project string, id string, commits string) error
	Publish(ctx context.Context, project string, id string) error
	List(ctx context.Context, project string, page, pageSize int) ([]*dto.Publish, int64, error)
	Detail(ctx context.Context, project string, id string) (*dto.PublishDetail, error)
	PublishStatuses(ctx context.Context, project string, id string) ([]*dto.PublishStatus, error)
}

func init() {
	autowire.Auto[IPublishModule](func() reflect.Value {
		m := new(imlPublishModule)
		gateway.RegisterInitHandleFunc(m.initGateway)
		return reflect.ValueOf(m)
	})
}
