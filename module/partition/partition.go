package partition

import (
	"context"
	"reflect"
	
	parition_dto "github.com/eolinker/apipark/module/partition/dto"
	"github.com/eolinker/go-common/autowire"
)

type IPartitionModule interface {
	CreatePartition(ctx context.Context, partition *parition_dto.Create) (*parition_dto.Detail, error)
	Search(ctx context.Context, keyword string) ([]*parition_dto.Item, error)
	Get(ctx context.Context, id string) (*parition_dto.Detail, error)
	Update(ctx context.Context, id string, edit *parition_dto.Edit) (*parition_dto.Detail, error)
	Delete(ctx context.Context, id string) error
	Simple(ctx context.Context) ([]*parition_dto.Simple, error)
	SimpleByIds(ctx context.Context, ids []string) ([]*parition_dto.Simple, error)
	SimpleWithCluster(ctx context.Context) ([]*parition_dto.SimpleWithCluster, error)
	CheckCluster(ctx context.Context, address ...string) ([]*parition_dto.Node, error)
	ResetCluster(ctx context.Context, partitionId string, address string) ([]*parition_dto.Node, error)
	ClusterNodes(ctx context.Context, partitionId string) ([]*parition_dto.Node, error)
}

func init() {
	autowire.Auto[IPartitionModule](func() reflect.Value {
		return reflect.ValueOf(new(imlPartition))
	})
}
