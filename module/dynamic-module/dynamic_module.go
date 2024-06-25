package dynamic_module

import (
	"context"
	"reflect"

	"github.com/eolinker/apipark/gateway"

	"github.com/eolinker/go-common/autowire"

	dynamic_module_dto "github.com/eolinker/apipark/module/dynamic-module/dto"
)

type IDynamicModuleModule interface {
	Create(ctx context.Context, module string, input *dynamic_module_dto.CreateDynamicModule) (*dynamic_module_dto.DynamicModule, error)
	Edit(ctx context.Context, module string, id string, input *dynamic_module_dto.EditDynamicModule) (*dynamic_module_dto.DynamicModule, error)
	Delete(ctx context.Context, module string, ids []string) error
	Get(ctx context.Context, module string, id string) (*dynamic_module_dto.DynamicModule, error)
	List(ctx context.Context, module string, keyword string, page int, pageSize int) ([]map[string]interface{}, int64, error)
	PluginInfo(ctx context.Context, module string, partitionIds ...string) (*dynamic_module_dto.PluginInfo, error)
	Render(ctx context.Context, module string) (map[string]interface{}, error)
	ModuleDrivers(ctx context.Context, group string) ([]*dynamic_module_dto.ModuleDriver, error)

	Online(ctx context.Context, module string, id string, partitionInput *dynamic_module_dto.PartitionInput) error
	Offline(ctx context.Context, module string, id string, partitionInput *dynamic_module_dto.PartitionInput) error
	PartitionStatuses(ctx context.Context, module string, keyword string, page int, pageSize int) (map[string]map[string]string, error)
	PartitionStatus(ctx context.Context, module string, id string) (*dynamic_module_dto.OnlineInfo, error)
}

func init() {
	autowire.Auto[IDynamicModuleModule](func() reflect.Value {
		m := new(imlDynamicModule)
		gateway.RegisterInitHandleFunc(m.initGateway)
		return reflect.ValueOf(m)
	})
}
