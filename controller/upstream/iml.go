package upstream

import (
	"github.com/eolinker/apipark/module/cluster"
	"github.com/eolinker/apipark/module/project"
	"github.com/eolinker/apipark/module/upstream"
	upstream_dto "github.com/eolinker/apipark/module/upstream/dto"
	"github.com/gin-gonic/gin"
)

var (
	_ IUpstreamController = (*imlUpstreamController)(nil)
)

type imlUpstreamController struct {
	upstreamModule  upstream.IUpstreamModule `autowired:""`
	projectModule   project.IProjectModule   `autowired:""`
	partitionModule cluster.IClusterModule   `autowired:""`
}

func (i *imlUpstreamController) Get(ctx *gin.Context, pid string) (upstream_dto.UpstreamConfig, error) {
	return i.upstreamModule.Get(ctx, pid)
}

func (i *imlUpstreamController) Save(ctx *gin.Context, pid string, upstream *upstream_dto.UpstreamConfig) (upstream_dto.UpstreamConfig, error) {
	return i.upstreamModule.Save(ctx, pid, *upstream)
}
