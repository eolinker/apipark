package service

import (
	"github.com/gin-gonic/gin"
	"reflect"

	service_dto "github.com/eolinker/apipark/module/service/dto"

	"github.com/eolinker/go-common/autowire"
)

type IServiceController interface {
	Get(ctx *gin.Context, pid string, sid string) (*service_dto.Service, error)
	Search(ctx *gin.Context, keyword string, pid string) ([]*service_dto.ServiceItem, error)
	Create(ctx *gin.Context, pid string, input *service_dto.CreateService) (*service_dto.Service, error)
	Edit(ctx *gin.Context, pid string, sid string, input *service_dto.EditService) (*service_dto.Service, error)
	Delete(ctx *gin.Context, pid string, sid string) error
	Enable(ctx *gin.Context, pid string, sid string) error
	Disable(ctx *gin.Context, pid string, sid string) error
	ServiceDoc(ctx *gin.Context, pid string, sid string) (*service_dto.ServiceDoc, error)
	SaveServiceDoc(ctx *gin.Context, pid string, sid string, input *service_dto.SaveServiceDoc) error
	ServiceApis(ctx *gin.Context, pid string, sid string) ([]*service_dto.ServiceApi, error)
	BindServiceApi(ctx *gin.Context, pid string, sid string, apis *service_dto.BindApis) error
	UnbindServiceApi(ctx *gin.Context, pid string, sid string, api string) error
	SortApis(ctx *gin.Context, pid string, sid string, apis *service_dto.BindApis) error
	SimpleList(ctx *gin.Context, pid string) ([]*service_dto.SimpleItem, error)
}

func init() {
	autowire.Auto[IServiceController](func() reflect.Value {
		return reflect.ValueOf(new(imlServiceController))
	})
}
