package partition

import (
	"errors"
	
	"github.com/eolinker/apipark/module/partition"
	parition_dto "github.com/eolinker/apipark/module/partition/dto"
	"github.com/eolinker/go-common/auto"
	"github.com/gin-gonic/gin"
)

var (
	_ IPartitionController = (*imlPartition)(nil)
)

type imlPartition struct {
	partitionModule partition.IPartitionModule `autowired:""`
}

func (p *imlPartition) Nodes(ctx *gin.Context, partitionId string) ([]*parition_dto.Node, error) {
	return p.partitionModule.ClusterNodes(ctx, partitionId)
}

func (p *imlPartition) ResetCluster(ctx *gin.Context, partitionId string, input *parition_dto.ResetCluster) ([]*parition_dto.Node, error) {
	return p.partitionModule.ResetCluster(ctx, partitionId, input.ManagerAddress)
}

func (p *imlPartition) Check(ctx *gin.Context, input *parition_dto.CheckCluster) ([]*parition_dto.Node, error) {
	return p.partitionModule.CheckCluster(ctx, input.Address)
}

func (p *imlPartition) SimpleWithCluster(ctx *gin.Context) ([]*parition_dto.SimpleWithCluster, error) {
	return p.partitionModule.SimpleWithCluster(ctx)
}

func (p *imlPartition) Delete(ctx *gin.Context, id string) (string, error) {
	err := p.partitionModule.Delete(ctx, id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (p *imlPartition) Search(ctx *gin.Context, keyword string) ([]*parition_dto.Item, error) {
	return p.partitionModule.Search(ctx, keyword)
}

func (p *imlPartition) Simple(ctx *gin.Context) ([]*parition_dto.Simple, error) {
	return p.partitionModule.Simple(ctx)
}

func (p *imlPartition) Info(ctx *gin.Context, id string) (*parition_dto.Detail, error) {
	if id == "" {
		return nil, errors.New("id is empty")
	}
	return p.partitionModule.Get(ctx, id)
}

func (p *imlPartition) Update(ctx *gin.Context, id string, input *parition_dto.Edit) (*parition_dto.Detail, error) {
	return p.partitionModule.Update(ctx, id, input)
}

func (p *imlPartition) Create(ctx *gin.Context, input *parition_dto.Create) (*parition_dto.Detail, string, auto.TimeLabel, error) {
	detail, err := p.partitionModule.CreatePartition(ctx, input)
	if err != nil {
		return nil, "", auto.TimeLabel{}, err
	}
	return detail, detail.Id, detail.UpdateTime, nil
}
