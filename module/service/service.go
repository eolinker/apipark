package service

import (
	"context"
	"reflect"

	service_dto "github.com/eolinker/apipark/module/service/dto"

	"github.com/eolinker/go-common/autowire"
)

type IServiceModule interface {
	// Search 关键字获取服务列表
	Search(ctx context.Context, keyword string, pid string) ([]*service_dto.ServiceItem, error)
	// Create 创建服务
	Create(ctx context.Context, pid string, input *service_dto.CreateService) (*service_dto.Service, error)
	// Edit 编辑服务
	Edit(ctx context.Context, pid string, sid string, input *service_dto.EditService) (*service_dto.Service, error)
	// Delete 删除服务
	Delete(ctx context.Context, pid string, sid string) error
	// Get 获取服务详情
	Get(ctx context.Context, pid string, sid string) (*service_dto.Service, error)
	// Enable 启用服务
	Enable(ctx context.Context, pid string, sid string) error
	// Disable 禁用服务
	Disable(ctx context.Context, pid string, sid string) error
	// ServiceDoc 服务文档
	ServiceDoc(ctx context.Context, pid string, sid string) (*service_dto.ServiceDoc, error)
	// SaveServiceDoc 保存服务文档
	SaveServiceDoc(ctx context.Context, pid string, sid string, input *service_dto.SaveServiceDoc) error
	// ServiceApis 服务API
	ServiceApis(ctx context.Context, pid string, sid string) ([]*service_dto.ServiceApi, error)
	// BindServiceApi 绑定服务API
	BindServiceApi(ctx context.Context, pid string, sid string, apis *service_dto.BindApis) error
	// UnbindServiceApi 解绑服务API
	UnbindServiceApi(ctx context.Context, pid string, sid string, apis []string) error
	// SortApis 排序API
	SortApis(ctx context.Context, pid string, sid string, apis *service_dto.BindApis) error
	// SimpleList 简单列表
	SimpleList(ctx context.Context, pid string) ([]*service_dto.SimpleItem, error)
}

func init() {
	autowire.Auto[IServiceModule](func() reflect.Value {
		return reflect.ValueOf(new(imlServiceModule))
	})

}
