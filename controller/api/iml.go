package api

import (
	"github.com/eolinker/apipark/module/api"
	api_dto "github.com/eolinker/apipark/module/api/dto"
	"github.com/gin-gonic/gin"
)

var _ IAPIController = (*imlAPIController)(nil)

type imlAPIController struct {
	module api.IApiModule `autowired:""`
}

func (i *imlAPIController) SimpleList(ctx *gin.Context, partitionId string, input *api_dto.ListInput) ([]*api_dto.ApiSimpleItem, error) {
	return i.module.SimpleList(ctx, partitionId, input)
}

func (i *imlAPIController) Detail(ctx *gin.Context, pid string, aid string) (*api_dto.ApiDetail, error) {
	return i.module.Detail(ctx, pid, aid)
}

func (i *imlAPIController) SimpleDetail(ctx *gin.Context, pid string, aid string) (*api_dto.ApiSimpleDetail, error) {
	return i.module.SimpleDetail(ctx, pid, aid)
}

func (i *imlAPIController) Search(ctx *gin.Context, keyword string, pid string) ([]*api_dto.ApiItem, error) {
	return i.module.Search(ctx, keyword, pid)
}

func (i *imlAPIController) SimpleSearch(ctx *gin.Context, keyword string, pid string) ([]*api_dto.ApiSimpleItem, error) {
	return i.module.SimpleSearch(ctx, keyword, pid)
}

func (i *imlAPIController) Create(ctx *gin.Context, pid string, dto *api_dto.CreateApi) (*api_dto.ApiSimpleDetail, error) {
	return i.module.Create(ctx, pid, dto)
}

func (i *imlAPIController) Edit(ctx *gin.Context, pid string, aid string, dto *api_dto.EditApi) (*api_dto.ApiSimpleDetail, error) {
	return i.module.Edit(ctx, pid, aid, dto)
}

func (i *imlAPIController) Delete(ctx *gin.Context, pid string, aid string) error {
	return i.module.Delete(ctx, pid, aid)
}

func (i *imlAPIController) Copy(ctx *gin.Context, pid string, aid string, dto *api_dto.CreateApi) (*api_dto.ApiSimpleDetail, error) {
	return i.module.Copy(ctx, pid, aid, dto)
}

func (i *imlAPIController) ApiDocDetail(ctx *gin.Context, pid string, aid string) (*api_dto.ApiDocDetail, error) {
	return i.module.ApiDocDetail(ctx, pid, aid)
}

func (i *imlAPIController) ApiProxyDetail(ctx *gin.Context, pid string, aid string) (*api_dto.ApiProxyDetail, error) {
	return i.module.ApiProxyDetail(ctx, pid, aid)
}

func (i *imlAPIController) Prefix(ctx *gin.Context, pid string) (string, bool, error) {
	prefix, err := i.module.Prefix(ctx, pid)
	if err != nil {
		return "", false, err
	}
	return prefix, true, nil
}
