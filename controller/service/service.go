package service

import (
	"reflect"

	"github.com/gin-gonic/gin"

	service_dto "github.com/eolinker/apipark/module/service/dto"

	"github.com/eolinker/go-common/autowire"
)

type IServiceController interface {
	ServiceDoc(ctx *gin.Context, pid string) (*service_dto.ServiceDoc, error)
	SaveServiceDoc(ctx *gin.Context, pid string, input *service_dto.SaveServiceDoc) error
}

func init() {
	autowire.Auto[IServiceController](func() reflect.Value {
		return reflect.ValueOf(new(imlServiceController))
	})
}
