package upstream

import (
	"github.com/eolinker/apipark/module/cluster"
	"github.com/eolinker/apipark/module/service"
	"github.com/eolinker/apipark/module/upstream"
	upstream_dto "github.com/eolinker/apipark/module/upstream/dto"
	"github.com/gin-gonic/gin"
)

var (
	_ IUpstreamController = (*imlUpstreamController)(nil)
)

type imlUpstreamController struct {
	upstreamModule  upstream.IUpstreamModule `autowired:""`
	projectModule   service.IServiceModule   `autowired:""`
	partitionModule cluster.IClusterModule   `autowired:""`
}

func (i *imlUpstreamController) Get(ctx *gin.Context, serviceId string) (upstream_dto.UpstreamConfig, error) {
	return i.upstreamModule.Get(ctx, serviceId)
}

func (i *imlUpstreamController) Save(ctx *gin.Context, serviceId string, upstream *upstream_dto.UpstreamConfig) (upstream_dto.UpstreamConfig, error) {
	return i.upstreamModule.Save(ctx, serviceId, *upstream)
}
