package subscribe

import (
	"context"
	"reflect"

	"github.com/eolinker/apipark/service/universally"
	"github.com/eolinker/go-common/autowire"
)

type ISubscribeService interface {
	universally.IServiceGet[Subscribe]
	universally.IServiceDelete
	universally.IServiceCreate[CreateSubscribe]
	universally.IServiceEdit[UpdateSubscribe]
	DeleteByApplication(ctx context.Context, service string, application string) error
	ListByApplication(ctx context.Context, service string, application ...string) ([]*Subscribe, error)
	ListByServices(ctx context.Context, serviceIds ...string) ([]*Subscribe, error)

	MySubscribeServices(ctx context.Context, application string, projectIds []string, serviceIDs []string, partitionIds ...string) ([]*Subscribe, error)
	UpdateSubscribeStatus(ctx context.Context, application string, service string, status int) error
	ListBySubscribeStatus(ctx context.Context, projectId string, status int) ([]*Subscribe, error)
	SubscribersByProject(ctx context.Context, projectIds ...string) ([]*Subscribe, error)
	Subscribers(ctx context.Context, project string, status int) ([]*Subscribe, error)
	SubscriptionsByApplication(ctx context.Context, applicationIds ...string) ([]*Subscribe, error)
}

type ISubscribeApplyService interface {
	universally.IServiceGet[Apply]
	universally.IServiceDelete
	universally.IServiceCreate[CreateApply]
	universally.IServiceEdit[EditApply]
	ListByStatus(ctx context.Context, pid string, status ...int) ([]*Apply, error)
	Revoke(ctx context.Context, service string, application string) error
	RevokeById(ctx context.Context, id string) error
}

func init() {
	autowire.Auto[ISubscribeService](func() reflect.Value {
		return reflect.ValueOf(new(imlSubscribeService))
	})

	autowire.Auto[ISubscribeApplyService](func() reflect.Value {
		return reflect.ValueOf(new(imlSubscribeApplyService))
	})
}
