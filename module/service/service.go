package service

import (
	"context"
	"reflect"

	service_dto "github.com/eolinker/apipark/module/service/dto"

	"github.com/eolinker/go-common/autowire"
)

type IServiceModule interface {
	// ServiceDoc 服务文档
	ServiceDoc(ctx context.Context, pid string) (*service_dto.ServiceDoc, error)
	// SaveServiceDoc 保存服务文档
	SaveServiceDoc(ctx context.Context, pid string, input *service_dto.SaveServiceDoc) error
}

func init() {
	autowire.Auto[IServiceModule](func() reflect.Value {
		return reflect.ValueOf(new(imlServiceModule))
	})

}
