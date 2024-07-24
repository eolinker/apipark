package api

import (
	"context"
	"reflect"

	"github.com/eolinker/go-common/autowire"

	api_dto "github.com/eolinker/apipark/module/api/dto"
)

type IApiModule interface {
	// Detail 获取API详情
	Detail(ctx context.Context, pid string, aid string) (*api_dto.ApiDetail, error)
	// SimpleDetail 获取API简要详情
	SimpleDetail(ctx context.Context, pid string, aid string) (*api_dto.ApiSimpleDetail, error)
	// Search 获取API列表
	Search(ctx context.Context, keyword string, pid string) ([]*api_dto.ApiItem, error)
	// SimpleSearch 获取API简要列表
	SimpleSearch(ctx context.Context, keyword string, pid string) ([]*api_dto.ApiSimpleItem, error)
	SimpleList(ctx context.Context, input *api_dto.ListInput) ([]*api_dto.ApiSimpleItem, error)
	// Create 创建API
	Create(ctx context.Context, pid string, dto *api_dto.CreateApi) (*api_dto.ApiSimpleDetail, error)
	// Edit 编辑API
	Edit(ctx context.Context, pid string, aid string, dto *api_dto.EditApi) (*api_dto.ApiSimpleDetail, error)
	// Delete 删除API
	Delete(ctx context.Context, pid string, aid string) error
	// Copy 复制API
	Copy(ctx context.Context, pid string, aid string, dto *api_dto.CreateApi) (*api_dto.ApiSimpleDetail, error)
	// ApiDocDetail 获取API文档详情
	ApiDocDetail(ctx context.Context, pid string, aid string) (*api_dto.ApiDocDetail, error)
	// ApiProxyDetail 获取API代理详情
	ApiProxyDetail(ctx context.Context, pid string, aid string) (*api_dto.ApiProxyDetail, error)
	// Prefix 获取API前缀
	Prefix(ctx context.Context, pid string) (string, error)
}

func init() {
	autowire.Auto[IApiModule](func() reflect.Value {
		return reflect.ValueOf(new(imlApiModule))
	})
}
