package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/eolinker/apipark/service/partition"
	"github.com/eolinker/apipark/service/subscribe"
	"github.com/eolinker/eosc/log"

	"github.com/eolinker/apipark/gateway"

	"github.com/eolinker/apipark/service/cluster"

	"github.com/eolinker/apipark/service/api"

	"gorm.io/gorm"

	"github.com/eolinker/go-common/utils"

	"github.com/eolinker/go-common/auto"

	"github.com/eolinker/apipark/service/tag"

	"github.com/eolinker/apipark/service/project"
	"github.com/google/uuid"

	serviceDto "github.com/eolinker/apipark/module/service/dto"
	"github.com/eolinker/apipark/service/service"
	"github.com/eolinker/go-common/store"
)

var (
	_                     IServiceModule = (*imlServiceModule)(nil)
	projectRuleMustServer                = map[string]bool{
		"as_server": true,
	}
)

type imlServiceModule struct {
	projectService          project.IProjectService     `autowired:""`
	serviceService          service.IServiceService     `autowired:""`
	serviceTagService       service.ITagService         `autowired:""`
	servicePartitionService service.IPartitionsService  `autowired:""`
	serviceDocService       service.IDocService         `autowired:""`
	tagService              tag.ITagService             `autowired:""`
	serviceApiService       service.IApiService         `autowired:""`
	apiService              api.IAPIService             `autowired:""`
	subscribeService        subscribe.ISubscribeService `autowired:""`
	partitionService        partition.IPartitionService `autowired:""`
	clusterService          cluster.IClusterService     `autowired:""`

	transaction store.ITransaction `autowired:""`
}

func (i *imlServiceModule) getServiceApis(ctx context.Context, projectIds []string) ([]*gateway.ServiceRelease, error) {
	services, err := i.serviceService.ListByProject(ctx, projectIds...)
	if err != nil {
		return nil, err
	}
	serviceIds := utils.SliceToSlice(services, func(s *service.Service) string {
		return s.Id
	}, func(s *service.Service) bool {
		return s.Status == service.StatusOn
	})
	apis, err := i.serviceApiService.List(ctx, serviceIds...)
	if err != nil {
		return nil, err
	}
	serviceMap := utils.SliceToMapArrayO(apis, func(s *service.Api) (string, string) {
		return s.Sid, s.Aid
	})

	return utils.MapToSlice(serviceMap, func(k string, v []string) *gateway.ServiceRelease {
		return &gateway.ServiceRelease{
			ID:   k,
			Apis: v,
		}
	}), nil
}

func (i *imlServiceModule) initGateway(ctx context.Context, partitionId string, clientDriver gateway.IClientDriver) error {
	projects, err := i.projectService.List(ctx)
	projectIds := utils.SliceToSlice(projects, func(p *project.Project) string {
		return p.Id
	})
	releases, err := i.getServiceApis(ctx, projectIds)
	if err != nil {
		return err
	}

	return clientDriver.Service().Online(ctx, releases...)

}

func (i *imlServiceModule) SimpleList(ctx context.Context, pid string) ([]*serviceDto.SimpleItem, error) {
	_, err := i.projectService.CheckProject(ctx, pid, projectRuleMustServer)

	if err != nil {
		return nil, err
	}
	list, err := i.serviceService.Search(ctx, "", map[string]interface{}{
		"project": pid,
	})
	return utils.SliceToSlice(list, func(i *service.Service) *serviceDto.SimpleItem {
		return &serviceDto.SimpleItem{
			Id:   i.Id,
			Name: i.Name,
		}
	}), nil
}

func (i *imlServiceModule) syncServiceApi(ctx context.Context, sid string, apiIds []string, partitionIds []string, isOnline bool) error {
	partitions, err := i.partitionService.List(ctx, partitionIds...)
	if err != nil {
		return err
	}

	for _, p := range partitions {
		err := i.syncServiceApiForCluster(ctx, p.Cluster, sid, apiIds, isOnline)
		if err != nil {
			log.Warnf("sync service api for partition[%s] %s", p.UUID, err.Error())
		}
	}
	return nil
}
func (i *imlServiceModule) syncServiceApiForCluster(ctx context.Context, clusterId string, sid string, apiIds []string, isOnline bool) error {
	client, err := i.clusterService.GatewayClient(ctx, clusterId)
	if err != nil {
		return err
	}
	defer func() {
		err := client.Close(ctx)
		if err != nil {
			log.Warn("close apinto client:", err)
		}
	}()
	if !isOnline {
		err = client.Service().Offline(ctx, &gateway.ServiceRelease{
			ID:   sid,
			Apis: apiIds,
		})
		if err != nil {
			return err
		}
	} else {
		err = client.Service().Online(ctx, &gateway.ServiceRelease{
			ID:   sid,
			Apis: apiIds,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
func (i *imlServiceModule) BindServiceApi(ctx context.Context, pid string, sid string, apis *serviceDto.BindApis) error {
	if len(apis.Apis) == 0 {
		return nil
	}
	_, err := i.projectService.CheckProject(ctx, pid, projectRuleMustServer)
	if err != nil {
		return fmt.Errorf("get project(%s) error: %v", pid, err)
	}
	sInfo, err := i.serviceService.Get(ctx, sid)
	if err != nil {
		return fmt.Errorf("get service(%s) error: %v", sid, err)
	}

	partitions, err := i.servicePartitionService.List(ctx, sid)
	if err != nil {
		return err
	}
	partitionIds := utils.SliceToSlice(partitions, func(p *service.Partition) string {
		return p.Pid
	})
	return i.transaction.Transaction(ctx, func(ctx context.Context) error {
		lastSort, err := i.serviceApiService.LastSortIndex(ctx, sid)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}
		for _, aid := range apis.Apis {
			lastSort++
			err = i.serviceApiService.Bind(ctx, sid, aid, lastSort)
			if err != nil {
				return err
			}
		}
		if sInfo.Status == service.StatusOff {
			// 未开启，不上线服务
			return nil
		}
		return i.syncServiceApi(ctx, sid, apis.Apis, partitionIds, true)
	})
}

func (i *imlServiceModule) UnbindServiceApi(ctx context.Context, pid string, sid string, apiIds []string) error {
	if len(apiIds) == 0 {
		return nil
	}
	_, err := i.projectService.CheckProject(ctx, pid, projectRuleMustServer)

	if err != nil {
		return err
	}
	sInfo, err := i.serviceService.Get(ctx, sid)
	if err != nil {
		return err
	}

	partitions, err := i.servicePartitionService.List(ctx, sid)
	if err != nil {
		return err
	}
	partitionIds := utils.SliceToSlice(partitions, func(p *service.Partition) string {
		return p.Pid
	})
	return i.transaction.Transaction(ctx, func(ctx context.Context) error {
		for _, aid := range apiIds {
			err = i.serviceApiService.Unbind(ctx, sid, aid)
			if err != nil {
				return err
			}
		}

		if sInfo.Status == service.StatusOff {
			// 未开启，不下线api
			return nil
		}
		return i.syncServiceApi(ctx, sid, apiIds, partitionIds, false)
	})
}

func (i *imlServiceModule) SortApis(ctx context.Context, pid string, sid string, apis *serviceDto.BindApis) error {
	_, err := i.projectService.CheckProject(ctx, pid, projectRuleMustServer)

	if err != nil {
		return err
	}
	if len(apis.Apis) == 0 {
		return nil
	}
	return i.transaction.Transaction(ctx, func(ctx context.Context) error {
		err = i.serviceApiService.Clear(ctx, sid)
		if err != nil {
			return err
		}
		for idx, aid := range apis.Apis {
			err = i.serviceApiService.Bind(ctx, sid, aid, idx+1)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (i *imlServiceModule) ServiceApis(ctx context.Context, pid string, sid string) ([]*serviceDto.ServiceApi, error) {
	_, err := i.projectService.CheckProject(ctx, pid, projectRuleMustServer)

	if err != nil {
		return nil, err
	}
	apis, err := i.serviceApiService.List(ctx, sid)
	if err != nil {
		return nil, err
	}
	sorts := make([]*serviceDto.ServiceApi, 0, len(apis))
	for _, a := range apis {
		apiInfo, err := i.apiService.GetInfo(ctx, a.Aid)
		if err != nil {
			return nil, err
		}
		sorts = append(sorts, &serviceDto.ServiceApi{
			Id:          a.Aid,
			Name:        apiInfo.Name,
			Method:      apiInfo.Method,
			Path:        apiInfo.Path,
			Description: apiInfo.Description,
		})
	}
	return sorts, nil
}

func (i *imlServiceModule) ServiceDoc(ctx context.Context, pid string, sid string) (*serviceDto.ServiceDoc, error) {
	_, err := i.projectService.CheckProject(ctx, pid, projectRuleMustServer)

	if err != nil {
		return nil, err
	}
	info, err := i.serviceService.Get(ctx, sid)
	if err != nil {
		return nil, err
	}
	doc, err := i.serviceDocService.Get(ctx, sid)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return &serviceDto.ServiceDoc{
			Id:   sid,
			Name: info.Name,
			Doc:  "",
		}, nil
	}
	return &serviceDto.ServiceDoc{
		Id:         sid,
		Name:       info.Name,
		Doc:        doc.Doc,
		Creator:    auto.UUID(doc.Creator),
		CreateTime: auto.TimeLabel(doc.CreateTime),
		Updater:    auto.UUID(doc.Updater),
		UpdateTime: auto.TimeLabel(doc.UpdateTime),
	}, nil
}

func (i *imlServiceModule) SaveServiceDoc(ctx context.Context, pid string, sid string, input *serviceDto.SaveServiceDoc) error {
	_, err := i.projectService.CheckProject(ctx, pid, projectRuleMustServer)

	if err != nil {
		return err
	}
	return i.serviceDocService.Save(ctx, &service.SaveDoc{
		Sid: sid,
		Doc: input.Doc,
	})
}

func (i *imlServiceModule) Search(ctx context.Context, keyword string, pid string) ([]*serviceDto.ServiceItem, error) {
	_, err := i.projectService.CheckProject(ctx, pid, projectRuleMustServer)

	if err != nil {
		return nil, err
	}

	services, err := i.serviceService.Search(ctx, keyword, map[string]interface{}{
		"project": pid,
	}, "update_at desc", "create_at desc")
	if err != nil {
		return nil, err
	}
	serviceIds := utils.SliceToSlice(services, func(s *service.Service) string {
		return s.Id
	})
	partitionMap, err := i.servicePartitionService.PartitionsByService(ctx, serviceIds...)
	if err != nil {
		return nil, err
	}
	subscribers, err := i.subscribeService.ListByServices(ctx, serviceIds...)
	if err != nil {
		return nil, err
	}
	subscribeMap := make(map[string]struct{})
	for _, s := range subscribers {
		if s.ApplyStatus != subscribe.ApplyStatusReview && s.ApplyStatus != subscribe.ApplyStatusSubscribe {
			continue
		}
		if _, ok := subscribeMap[s.Service]; ok {
			continue
		}
		subscribeMap[s.Service] = struct{}{}
	}

	items := make([]*serviceDto.ServiceItem, 0, len(serviceIds))
	for _, s := range services {

		apiCount, err := i.serviceApiService.Count(ctx, s.Id)
		if err != nil {
			return nil, err
		}
		_, notDelete := subscribeMap[s.Id]
		items = append(items, &serviceDto.ServiceItem{
			Id:          s.Id,
			Name:        s.Name,
			Partition:   auto.List(partitionMap[s.Id]),
			ServiceType: s.ServiceType,
			ApiNum:      apiCount,
			Status:      s.Status,
			CreateTime:  auto.TimeLabel(s.CreateTime),
			UpdateTime:  auto.TimeLabel(s.UpdateTime),
			CanDelete:   !notDelete,
		})

	}

	return items, nil
}

func (i *imlServiceModule) Create(ctx context.Context, pid string, input *serviceDto.CreateService) (*serviceDto.Service, error) {
	pInfo, err := i.projectService.CheckProject(ctx, pid, projectRuleMustServer)

	if err != nil {
		return nil, err
	}
	if input.ID == "" {
		input.ID = uuid.New().String()
	}
	if input.ServiceType != "inner" && input.Catalogue == nil {
		return nil, fmt.Errorf("group is required")
	}
	catalogue := ""
	if input.Catalogue != nil {
		catalogue = *input.Catalogue
	}

	err = i.transaction.Transaction(ctx, func(ctx context.Context) error {
		t := strings.Join(input.Tags, ",")
		err = i.serviceService.Create(ctx, &service.CreateService{
			Uuid:        input.ID,
			Name:        input.Name,
			Description: input.Description,
			Logo:        input.Logo,
			ServiceType: input.ServiceType,
			Project:     pid,
			Team:        pInfo.Team,
			Catalogue:   catalogue,
			Status:      service.StatusOff,
			Tag:         t,
		})
		if err != nil {
			return err
		}
		err = i.servicePartitionService.Delete(ctx, nil, []string{input.ID})
		if err != nil {
			return err
		}
		for _, p := range input.Partition {
			err = i.servicePartitionService.Create(ctx, &service.CreatePartition{
				Sid: input.ID,
				Pid: p,
			})
			if err != nil {
				return err
			}
		}
		tags, err := i.getTagUuids(ctx, input.Tags)
		if err != nil {
			return err
		}
		err = i.serviceTagService.Delete(ctx, nil, []string{input.ID})
		if err != nil {
			return err
		}
		for _, t := range tags {
			err = i.serviceTagService.Create(ctx, &service.CreateTag{
				Tid: t,
				Sid: input.ID,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return i.Get(ctx, pid, input.ID)
}

func (i *imlServiceModule) getTagUuids(ctx context.Context, tags []string) ([]string, error) {
	list, err := i.tagService.Search(ctx, "", map[string]interface{}{"name": tags})
	if err != nil {
		return nil, err
	}
	tagMap := make(map[string]string)
	for _, t := range list {
		tagMap[t.Name] = t.Id
	}
	tagList := make([]string, 0, len(tags))
	repeatTag := make(map[string]struct{})
	for _, t := range tags {
		if _, ok := repeatTag[t]; ok {
			continue
		}
		repeatTag[t] = struct{}{}
		v := &tag.CreateTag{
			Name: t,
		}
		id, ok := tagMap[t]
		if !ok {
			v.Id = uuid.New().String()
			err = i.tagService.Create(ctx, v)
			if err != nil {
				return nil, err
			}
			tagMap[t] = v.Id
		} else {
			v.Id = id
		}
		tagList = append(tagList, v.Id)
	}
	return tagList, nil
}

func (i *imlServiceModule) Edit(ctx context.Context, pid string, sid string, input *serviceDto.EditService) (*serviceDto.Service, error) {
	_, err := i.projectService.CheckProject(ctx, pid, projectRuleMustServer)

	if err != nil {
		return nil, err
	}

	err = i.transaction.Transaction(ctx, func(ctx context.Context) error {
		t := strings.Join(input.Tags, ",")
		err = i.serviceService.Save(ctx, sid, &service.EditService{
			Name:        input.Name,
			Description: input.Description,
			Logo:        input.Logo,
			ServiceType: input.ServiceType,
			Catalogue:   input.Catalogue,
			Tag:         &t,
		})
		if err != nil {
			return err
		}
		tags, err := i.getTagUuids(ctx, input.Tags)
		if err != nil {
			return err
		}
		err = i.serviceTagService.Delete(ctx, nil, []string{sid})
		if err != nil {
			return err
		}
		for _, t := range tags {
			err = i.serviceTagService.Create(ctx, &service.CreateTag{
				Tid: t,
				Sid: sid,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return i.Get(ctx, pid, sid)
}

func (i *imlServiceModule) Delete(ctx context.Context, pid string, sid string) error {
	_, err := i.projectService.CheckProject(ctx, pid, projectRuleMustServer)

	if err != nil {
		return err
	}
	list, err := i.subscribeService.ListByServices(ctx, sid)
	if err != nil {
		return err
	}
	for _, l := range list {
		if l.ApplyStatus == subscribe.ApplyStatusSubscribe || l.ApplyStatus == subscribe.ApplyStatusReview {
			return fmt.Errorf("service %s is used", sid)
		}
	}
	partitions, err := i.servicePartitionService.List(ctx, sid)
	if err != nil {
		return err
	}
	partitionIds := make([]string, 0, len(partitions))
	for _, p := range partitions {
		partitionIds = append(partitionIds, p.Pid)
	}
	return i.transaction.Transaction(ctx, func(ctx context.Context) error {
		err := i.serviceService.Delete(ctx, sid)
		if err != nil {
			return err
		}
		err = i.servicePartitionService.Delete(ctx, nil, []string{sid})
		if err != nil {
			return err
		}
		err = i.serviceTagService.Delete(ctx, nil, []string{sid})
		if err != nil {
			return err
		}
		apis, err := i.serviceApiService.List(ctx, sid)
		if err != nil {
			return err
		}
		err = i.serviceApiService.Clear(ctx, sid)
		if err != nil {
			return err
		}
		apiIds := utils.SliceToSlice(apis, func(s *service.Api) string {
			return s.Aid
		})

		return i.syncServiceApi(ctx, sid, apiIds, partitionIds, false)
	})
}

func (i *imlServiceModule) Get(ctx context.Context, pid string, sid string) (*serviceDto.Service, error) {
	info, err := i.projectService.CheckProject(ctx, pid, projectRuleMustServer)

	if err != nil {
		return nil, err
	}
	s, err := i.serviceService.Get(ctx, sid)
	if err != nil {
		return nil, err
	}
	partitions, err := i.servicePartitionService.List(ctx, sid)
	if err != nil {
		return nil, err
	}
	partitionIds := make([]string, 0, len(partitions))
	for _, p := range partitions {
		partitionIds = append(partitionIds, p.Pid)
	}
	tags, err := i.serviceTagService.List(ctx, []string{sid}, nil)
	if err != nil {
		return nil, err
	}
	tagIds := make([]string, 0, len(tags))
	for _, t := range tags {
		tagIds = append(tagIds, t.Tid)
	}

	out := &serviceDto.Service{
		Id:          s.Id,
		Name:        s.Name,
		Description: s.Description,
		Logo:        s.Logo,
		ServiceType: s.ServiceType,
		Team:        auto.UUID(info.Team),
		Project:     auto.UUID(s.Project),
		Catalogue:   auto.UUID(s.Catalogue),
		Partition:   auto.List(partitionIds),
		Tags:        auto.List(tagIds),
		Status:      s.Status,
	}
	return out, nil
}

func (i *imlServiceModule) Enable(ctx context.Context, pid string, sid string) error {
	_, err := i.projectService.CheckProject(ctx, pid, projectRuleMustServer)

	if err != nil {
		return err
	}
	partitions, err := i.servicePartitionService.List(ctx, sid)
	if err != nil {
		return err
	}
	partitionIds := utils.SliceToSlice(partitions, func(p *service.Partition) string {
		return p.Pid
	})
	status := service.StatusOn
	return i.transaction.Transaction(ctx, func(ctx context.Context) error {
		err = i.serviceService.Save(ctx, sid, &service.EditService{
			Status: &status,
		})
		if err != nil {
			return err
		}
		apis, err := i.serviceApiService.List(ctx, sid)
		if err != nil {
			return err
		}
		if len(apis) == 0 {
			return nil
		}
		apiIds := utils.SliceToSlice(apis, func(i *service.Api) string {
			return i.Aid
		})
		return i.syncServiceApi(ctx, sid, apiIds, partitionIds, true)
	})

}

func (i *imlServiceModule) Disable(ctx context.Context, pid string, sid string) error {
	_, err := i.projectService.CheckProject(ctx, pid, projectRuleMustServer)

	if err != nil {
		return err
	}
	partitions, err := i.servicePartitionService.List(ctx, sid)
	if err != nil {
		return err
	}
	partitionIds := utils.SliceToSlice(partitions, func(p *service.Partition) string {
		return p.Pid
	})
	status := service.StatusOff
	return i.transaction.Transaction(ctx, func(ctx context.Context) error {
		err = i.serviceService.Save(ctx, sid, &service.EditService{
			Status: &status,
		})
		if err != nil {
			return err
		}
		apis, err := i.serviceApiService.List(ctx, sid)
		if err != nil {
			return err
		}
		if len(apis) == 0 {
			return nil
		}
		apiIds := utils.SliceToSlice(apis, func(i *service.Api) string {
			return i.Aid
		})
		return i.syncServiceApi(ctx, sid, apiIds, partitionIds, false)
	})
}
