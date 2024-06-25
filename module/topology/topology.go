package topology

import (
	"context"
	"reflect"

	"github.com/eolinker/go-common/autowire"

	topology_dto "github.com/eolinker/apipark/module/topology/dto"
)

type ITopologyModule interface {
	SystemTopology(ctx context.Context) ([]*topology_dto.ProjectItem, []*topology_dto.ServiceItem, error)
	ProjectTopology(ctx context.Context, projectID string) ([]*topology_dto.ServiceItem, []*topology_dto.TopologyItem, []*topology_dto.TopologyItem, error)
}

func init() {
	autowire.Auto[ITopologyModule](func() reflect.Value {
		return reflect.ValueOf(new(imlTopologyModule))
	})
}
