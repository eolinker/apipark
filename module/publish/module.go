package publish

import (
	"context"
	"reflect"

	"github.com/eolinker/apipark/gateway"

	"github.com/eolinker/apipark/module/publish/dto"
	"github.com/eolinker/go-common/autowire"
)

type IPublishModule interface {
	CheckPublish(ctx context.Context, serviceId string, releaseId string) (*dto.DiffOut, error)
	//ReleaseDo(ctx context.Context, serviceId string, input *dto.ApplyOnReleaseInput) error

	Apply(ctx context.Context, serviceId string, input *dto.ApplyInput) (*dto.Publish, error)
	Stop(ctx context.Context, serviceId string, id string) error
	Refuse(ctx context.Context, serviceId string, id string, commits string) error
	Accept(ctx context.Context, serviceId string, id string, commits string) error
	Publish(ctx context.Context, serviceId string, id string) error
	List(ctx context.Context, serviceId string, page, pageSize int) ([]*dto.Publish, int64, error)
	Detail(ctx context.Context, serviceId string, id string) (*dto.PublishDetail, error)
	PublishStatuses(ctx context.Context, serviceId string, id string) ([]*dto.PublishStatus, error)
}

func init() {
	autowire.Auto[IPublishModule](func() reflect.Value {
		m := new(imlPublishModule)
		gateway.RegisterInitHandleFunc(m.initGateway)
		return reflect.ValueOf(m)
	})
}
