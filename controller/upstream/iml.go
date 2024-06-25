package upstream

import (
	"github.com/eolinker/apipark/module/partition"
	partition_dto "github.com/eolinker/apipark/module/partition/dto"
	"github.com/eolinker/apipark/module/project"
	"github.com/eolinker/apipark/module/upstream"
	upstream_dto "github.com/eolinker/apipark/module/upstream/dto"
	"github.com/gin-gonic/gin"
)

var (
	_ IUpstreamController = (*imlUpstreamController)(nil)
)

type imlUpstreamController struct {
	upstreamModule  upstream.IUpstreamModule   `autowired:""`
	projectModule   project.IProjectModule     `autowired:""`
	partitionModule partition.IPartitionModule `autowired:""`
}

func (i *imlUpstreamController) Get(ctx *gin.Context, pid string) (upstream_dto.UpstreamConfig, []*partition_dto.Simple, error) {
	info, err := i.upstreamModule.Get(ctx, pid)
	if err != nil {
		return nil, nil, err
	}
	projectInfo, err := i.projectModule.GetProject(ctx, pid)
	if err != nil {
		return nil, nil, err
	}
	ids := make([]string, 0)
	for _, p := range projectInfo.Partition {
		ids = append(ids, p.Id)
	}
	items, err := i.partitionModule.SimpleByIds(ctx, ids)
	if err != nil {
		return nil, nil, err
	}
	return info, items, nil
}

func (i *imlUpstreamController) Save(ctx *gin.Context, pid string, upstream *upstream_dto.UpstreamConfig) (upstream_dto.UpstreamConfig, error) {
	return i.upstreamModule.Save(ctx, pid, *upstream)
}
