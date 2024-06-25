package topology

import (
	"github.com/gin-gonic/gin"
	"reflect"

	"github.com/eolinker/go-common/autowire"

	topology_dto "github.com/eolinker/apipark/module/topology/dto"
)

type ITopologyController interface {
	SystemTopology(ctx *gin.Context) ([]*topology_dto.ProjectItem, []*topology_dto.ServiceItem, error)
	ProjectTopology(ctx *gin.Context, projectID string) ([]*topology_dto.ServiceItem, []*topology_dto.TopologyItem, []*topology_dto.TopologyItem, error)
}

func init() {
	autowire.Auto[ITopologyController](func() reflect.Value {
		return reflect.ValueOf(new(imlTopologyController))
	})
}
