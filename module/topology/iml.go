package topology

import (
	"context"

	"github.com/eolinker/go-common/auto"

	"github.com/eolinker/apipark/service/subscribe"

	"github.com/eolinker/go-common/utils"

	"github.com/eolinker/apipark/service/service"

	topology_dto "github.com/eolinker/apipark/module/topology/dto"

	"github.com/eolinker/apipark/service/project"
)

var _ ITopologyModule = (*imlTopologyModule)(nil)

type imlTopologyModule struct {
	projectService   project.IProjectService     `autowired:""`
	subscribeService subscribe.ISubscribeService `autowired:""`
	serviceService   service.IServiceService     `autowired:""`
}

func (i *imlTopologyModule) SystemTopology(ctx context.Context) ([]*topology_dto.ProjectItem, []*topology_dto.ServiceItem, error) {
	projects, err := i.projectService.List(ctx)
	if err != nil {
		return nil, nil, err
	}
	subscriptions, err := i.subscribeService.ListBySubscribeStatus(ctx, "", subscribe.ApplyStatusSubscribe)
	if err != nil {
		return nil, nil, err
	}

	services, err := i.serviceService.List(ctx)
	if err != nil {
		return nil, nil, err
	}
	subscriptionMap := utils.SliceToMapArrayO(subscriptions, func(s *subscribe.Subscribe) (string, string) {
		return s.Application, s.Service
	})
	ps := utils.SliceToSlice(projects, func(p *project.Project) *topology_dto.ProjectItem {
		return &topology_dto.ProjectItem{
			ID:             p.Id,
			Name:           p.Name,
			InvokeServices: subscriptionMap[p.Id],
			IsApp:          p.AsApp,
			IsServer:       p.AsServer,
		}
	})
	ss := utils.SliceToSlice(services, func(s *service.Service) *topology_dto.ServiceItem {
		return &topology_dto.ServiceItem{
			ID:      s.Id,
			Name:    s.Name,
			Project: s.Project,
		}
	})
	return ps, ss, nil
}

func (i *imlTopologyModule) ProjectTopology(ctx context.Context, projectID string) ([]*topology_dto.ServiceItem, []*topology_dto.TopologyItem, []*topology_dto.TopologyItem, error) {
	// 获取系统中所有服务
	services, err := i.serviceService.ListByProject(ctx, projectID)
	if err != nil {
		return nil, nil, nil, err
	}
	// 获取当前项目的订阅关系
	mySubscribeServices, err := i.subscribeService.MySubscribeServices(ctx, projectID, nil, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	subscribers, err := i.subscribeService.Subscribers(ctx, projectID, subscribe.ApplyStatusSubscribe)
	if err != nil {
		return nil, nil, nil, err
	}
	projectMap := utils.SliceToMapArrayO(mySubscribeServices, func(s *subscribe.Subscribe) (string, auto.Label) {
		return s.Project, auto.UUID(s.Service)
	})
	subscriberMap := utils.SliceToMapArrayO(subscribers, func(s *subscribe.Subscribe) (string, auto.Label) {
		return s.Application, auto.UUID(s.Service)
	})

	return utils.SliceToSlice(services, func(s *service.Service) *topology_dto.ServiceItem {
			return &topology_dto.ServiceItem{
				ID:      s.Id,
				Name:    s.Name,
				Project: s.Project,
			}
		}), utils.MapToSlice(subscriberMap, func(k string, ss []auto.Label) *topology_dto.TopologyItem {
			return &topology_dto.TopologyItem{
				Project:  auto.UUID(k),
				Services: ss,
			}
		}), utils.MapToSlice(projectMap, func(k string, ss []auto.Label) *topology_dto.TopologyItem {
			return &topology_dto.TopologyItem{
				Project:  auto.UUID(k),
				Services: ss,
			}
		}), nil

}
