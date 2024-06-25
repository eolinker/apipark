package service

import (
	"encoding/json"
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

func (i *imlServiceController) SimpleList(ctx *gin.Context, pid string) ([]*service_dto.SimpleItem, error) {
	return i.module.SimpleList(ctx, pid)
}

func (i *imlServiceController) ServiceApis(ctx *gin.Context, pid string, sid string) ([]*service_dto.ServiceApi, error) {
	return i.module.ServiceApis(ctx, pid, sid)
}

func (i *imlServiceController) BindServiceApi(ctx *gin.Context, pid string, sid string, apis *service_dto.BindApis) error {
	return i.module.BindServiceApi(ctx, pid, sid, apis)
}

func (i *imlServiceController) UnbindServiceApi(ctx *gin.Context, pid string, sid string, api string) error {
	apis := make([]string, 0)
	err := json.Unmarshal([]byte(api), &apis)
	if err != nil {
		return err
	}
	return i.module.UnbindServiceApi(ctx, pid, sid, apis)
}

func (i *imlServiceController) SortApis(ctx *gin.Context, pid string, sid string, apis *service_dto.BindApis) error {
	return i.module.SortApis(ctx, pid, sid, apis)
}

func (i *imlServiceController) ServiceDoc(ctx *gin.Context, pid string, sid string) (*service_dto.ServiceDoc, error) {
	return i.module.ServiceDoc(ctx, pid, sid)
}

func (i *imlServiceController) SaveServiceDoc(ctx *gin.Context, pid string, sid string, input *service_dto.SaveServiceDoc) error {
	return i.module.SaveServiceDoc(ctx, pid, sid, input)
}

func (i *imlServiceController) Get(ctx *gin.Context, pid string, sid string) (*service_dto.Service, error) {
	return i.module.Get(ctx, pid, sid)
}

func (i *imlServiceController) Search(ctx *gin.Context, keyword string, pid string) ([]*service_dto.ServiceItem, error) {
	return i.module.Search(ctx, keyword, pid)
}

func (i *imlServiceController) Create(ctx *gin.Context, pid string, input *service_dto.CreateService) (*service_dto.Service, error) {
	return i.module.Create(ctx, pid, input)
}

func (i *imlServiceController) Edit(ctx *gin.Context, pid string, sid string, input *service_dto.EditService) (*service_dto.Service, error) {
	return i.module.Edit(ctx, pid, sid, input)
}

func (i *imlServiceController) Delete(ctx *gin.Context, pid string, sid string) error {
	return i.module.Delete(ctx, pid, sid)
}

func (i *imlServiceController) Enable(ctx *gin.Context, pid string, sid string) error {
	return i.module.Enable(ctx, pid, sid)
}

func (i *imlServiceController) Disable(ctx *gin.Context, pid string, sid string) error {
	return i.module.Disable(ctx, pid, sid)
}
