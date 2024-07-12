package partition

import (
	"reflect"
	
	parition_dto "github.com/eolinker/apipark/module/partition/dto"
	
	"github.com/eolinker/go-common/auto"
	"github.com/eolinker/go-common/autowire"
	"github.com/gin-gonic/gin"
)

type IPartitionController interface {
	Create(ctx *gin.Context, input *parition_dto.Create) (*parition_dto.Detail, string, auto.TimeLabel, error)
	Search(ctx *gin.Context, keyword string) ([]*parition_dto.Item, error)
	Simple(ctx *gin.Context) ([]*parition_dto.Simple, error)
	Info(ctx *gin.Context, id string) (*parition_dto.Detail, error)
	Update(ctx *gin.Context, id string, input *parition_dto.Edit) (*parition_dto.Detail, error)
	Delete(ctx *gin.Context, id string) (string, error)
	SimpleWithCluster(ctx *gin.Context) ([]*parition_dto.SimpleWithCluster, error)
	
	Nodes(ctx *gin.Context, partitionId string) ([]*parition_dto.Node, error)
	ResetCluster(ctx *gin.Context, partitionId string, input *parition_dto.ResetCluster) ([]*parition_dto.Node, error)
	Check(ctx *gin.Context, input *parition_dto.CheckCluster) ([]*parition_dto.Node, error)
}

func init() {
	autowire.Auto[IPartitionController](func() reflect.Value {
		return reflect.ValueOf(new(imlPartition))
	})
}
