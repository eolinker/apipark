package topology

import (
	"github.com/gin-gonic/gin"

	"github.com/eolinker/apipark/module/topology"
	topology_dto "github.com/eolinker/apipark/module/topology/dto"
)

var _ ITopologyController = (*imlTopologyController)(nil)

type imlTopologyController struct {
	module topology.ITopologyModule `autowired:""`
}

func (i *imlTopologyController) SystemTopology(ctx *gin.Context) ([]*topology_dto.ProjectItem, []*topology_dto.ServiceItem, error) {
	return i.module.SystemTopology(ctx)
}

func (i *imlTopologyController) ProjectTopology(ctx *gin.Context, projectID string) ([]*topology_dto.ServiceItem, []*topology_dto.TopologyItem, []*topology_dto.TopologyItem, error) {
	return i.module.ProjectTopology(ctx, projectID)
}
