package service

import (
	"github.com/gin-gonic/gin"

	"github.com/eolinker/apipark/module/service"
	service_dto "github.com/eolinker/apipark/module/service/dto"
)

var (
	_ IServiceController = (*imlServiceController)(nil)
)

type imlServiceController struct {
	module service.IServiceModule `autowired:""`
}

func (i *imlServiceController) ServiceDoc(ctx *gin.Context, pid string) (*service_dto.ServiceDoc, error) {
	return i.module.ServiceDoc(ctx, pid)
}

func (i *imlServiceController) SaveServiceDoc(ctx *gin.Context, pid string, input *service_dto.SaveServiceDoc) error {
	return i.module.SaveServiceDoc(ctx, pid, input)
}
