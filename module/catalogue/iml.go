package catalogue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/eolinker/apipark/service/partition"

	"github.com/eolinker/apipark/service/project"

	"github.com/eolinker/apipark/service/subscribe"

	"github.com/eolinker/go-common/store"

	"gorm.io/gorm"

	"github.com/eolinker/go-common/utils"

	"github.com/eolinker/apipark/service/release"

	"github.com/eolinker/apipark/service/api"
	"github.com/eolinker/go-common/auto"

	"github.com/eolinker/apipark/service/tag"

	"github.com/eolinker/apipark/service/service"

	"github.com/google/uuid"

	"github.com/eolinker/apipark/service/catalogue"

	catalogue_dto "github.com/eolinker/apipark/module/catalogue/dto"
)

var (
	_        ICatalogueModule = (*imlCatalogueModule)(nil)
	_sortMax                  = math.MaxInt32 / 2
)

type imlCatalogueModule struct {
	catalogueService        catalogue.ICatalogueService      `autowired:""`
	projectService          project.IProjectService          `autowired:""`
	apiService              api.IAPIService                  `autowired:""`
	serviceService          service.IServiceService          `autowired:""`
	serviceApiService       service.IApiService              `autowired:""`
	serviceTagService       service.ITagService              `autowired:""`
	servicePartitionService service.IPartitionsService       `autowired:""`
	serviceDocService       service.IDocService              `autowired:""`
	tagService              tag.ITagService                  `autowired:""`
	releaseService          release.IReleaseService          `autowired:""`
	subscribeService        subscribe.ISubscribeService      `autowired:""`
	subscribeApplyService   subscribe.ISubscribeApplyService `autowired:""`
	partitionService        partition.IPartitionService      `autowired:""`
	transaction             store.ITransaction               `autowired:""`

	root *Root
}

func (i *imlCatalogueModule) Subscribe(ctx context.Context, subscribeInfo *catalogue_dto.SubscribeService) error {
	if len(subscribeInfo.Applications) == 0 {
		return fmt.Errorf("applications is empty")
	}
	// 获取服务的基本信息
	s, err := i.serviceService.Get(ctx, subscribeInfo.Service)
	if err != nil {
		return fmt.Errorf("get service failed: %w", err)
	}

	userId := utils.UserId(ctx)
	return i.transaction.Transaction(ctx, func(ctx context.Context) error {

		projects := make([]string, 0, len(subscribeInfo.Applications))

		for _, pid := range subscribeInfo.Applications {
			if pid == s.Project {
				// 不能订阅自己
				continue
			}

			pInfo, err := i.projectService.Get(ctx, pid)
			if err != nil {
				return err
			}
			if !pInfo.AsApp {
				// 当系统不可作为订阅方时，不可订阅
				continue
			}
			applyID := uuid.New().String()
			// 创建一条审核申请
			err = i.subscribeApplyService.Create(ctx, &subscribe.CreateApply{
				Uuid:        applyID,
				Service:     subscribeInfo.Service,
				Project:     s.Project,
				Team:        s.Team,
				Application: pid,
				ApplyTeam:   pInfo.Team,
				Reason:      subscribeInfo.Reason,
				Status:      subscribe.ApplyStatusReview,
				Applier:     userId,
			})
			if err != nil {
				return err
			}
			// 修改订阅表状态
			subscriber, err := i.subscribeService.ListByApplication(ctx, subscribeInfo.Service, pid)
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}
				err = i.subscribeService.Create(ctx, &subscribe.CreateSubscribe{
					Uuid:        uuid.New().String(),
					Service:     subscribeInfo.Service,
					Project:     s.Project,
					Application: pid,
					ApplyStatus: subscribe.ApplyStatusReview,
					From:        subscribe.FromSubscribe,
				})
				if err != nil {
					return err
				}

			} else {
				subscriberMap := utils.SliceToMap(subscriber, func(t *subscribe.Subscribe) string {
					return t.Application
				})
				v, has := subscriberMap[pid]
				if !has {
					err = i.subscribeService.Create(ctx, &subscribe.CreateSubscribe{
						Uuid:        uuid.New().String(),
						Service:     subscribeInfo.Service,
						Project:     s.Project,
						Application: pid,
						ApplyStatus: subscribe.ApplyStatusReview,
						From:        subscribe.FromSubscribe,
					})
					if err != nil {
						return err
					}
				} else if v.ApplyStatus != subscribe.ApplyStatusSubscribe {
					status := subscribe.ApplyStatusReview
					err = i.subscribeService.Save(ctx, v.Id, &subscribe.UpdateSubscribe{
						ApplyStatus: &status,
					})
				}

			}

			projects = append(projects, pid)
		}
		if len(projects) == 0 {
			return fmt.Errorf("no available projects")
		}
		return nil
	})

}

func (i *imlCatalogueModule) ServiceDetail(ctx context.Context, sid string) (*catalogue_dto.ServiceDetail, error) {
	// 获取服务的基本信息
	s, err := i.serviceService.Get(ctx, sid)
	if err != nil {
		return nil, fmt.Errorf("get service failed: %w", err)
	}
	docStr := ""
	doc, err := i.serviceDocService.Get(ctx, sid)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get service doc failed: %w", err)
		}
	} else {
		docStr = doc.Doc
	}

	globalPartitions, err := i.partitionService.List(ctx)
	if err != nil {
		return nil, err
	}

	partitions := utils.SliceToSlice(globalPartitions, func(t *partition.Partition) *catalogue_dto.Partition {
		return &catalogue_dto.Partition{
			Id:     t.UUID,
			Name:   t.Name,
			Prefix: t.Prefix,
		}
	})
	r, err := i.releaseService.GetRunning(ctx, s.Project)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &catalogue_dto.ServiceDetail{
				Name:        s.Name,
				Description: s.Description,
				Document:    docStr,
				Basic: &catalogue_dto.ServiceBasic{
					Project: auto.UUID(s.Project),
					Team:    auto.UUID(s.Team),
					ApiNum:  0,
				},
				Partition: partitions,
			}, nil
		}

		return nil, fmt.Errorf("get running release failed: %w", err)
	}
	_, docCommits, _, err := i.releaseService.GetReleaseInfos(ctx, r.UUID)
	if err != nil {

		return nil, fmt.Errorf("get release apis failed: %w", err)
	}
	apiMap := utils.SliceToMap(docCommits, func(t *release.APIDocumentCommit) string {
		return t.API
	})

	apiList, err := i.serviceApiService.List(ctx, sid)
	if err != nil {
		return nil, err
	}
	apis := make([]*catalogue_dto.ServiceApi, 0, len(apiList))
	disableApis := make([]*catalogue_dto.ServiceApiBasic, 0, len(apiList))
	for _, a := range apiList {

		apiInfo, err := i.apiService.GetInfo(ctx, a.Aid)
		if err != nil {
			return nil, err
		}
		basicApi := &catalogue_dto.ServiceApiBasic{
			Id:          apiInfo.UUID,
			Name:        apiInfo.Name,
			Description: apiInfo.Description,
			Method:      apiInfo.Method,
			Path:        apiInfo.Path,
			Creator:     auto.UUID(apiInfo.Creator),
			Updater:     auto.UUID(apiInfo.Updater),
			CreateTime:  auto.TimeLabel(apiInfo.CreateAt),
			UpdateTime:  auto.TimeLabel(apiInfo.UpdateAt),
		}
		v, ok := apiMap[a.Aid]
		if !ok {
			disableApis = append(disableApis, basicApi)
			continue
		}
		commit, err := i.apiService.GetDocumentCommit(ctx, v.Commit)
		if err != nil {
			return nil, err
		}
		tmp := make(map[string]interface{})
		if commit.Data != nil {
			err = json.Unmarshal([]byte(commit.Data.Content), &tmp)
			if err != nil {
				return nil, err
			}
		}

		apis = append(apis, &catalogue_dto.ServiceApi{
			ServiceApiBasic: basicApi,
			Doc:             tmp,
		})
	}
	subscribers, err := i.subscribeService.ListByServices(ctx, sid)
	if err != nil {
		return nil, err
	}
	subscribeCount := map[string]int{}
	tmp := map[string]struct{}{}
	for _, s := range subscribers {
		key := fmt.Sprintf("%s-%s", s.Service, s.Application)
		if _, ok := tmp[key]; !ok {
			tmp[key] = struct{}{}
			subscribeCount[s.Service]++
		}
	}
	return &catalogue_dto.ServiceDetail{
		Name:        s.Name,
		Description: s.Description,
		Document:    docStr,
		Basic: &catalogue_dto.ServiceBasic{
			Project:       auto.UUID(s.Project),
			Team:          auto.UUID(s.Team),
			ApiNum:        len(apis),
			SubscriberNum: subscribeCount[s.Id],
		},
		Apis:      apis,
		Partition: partitions,
	}, nil
}

func (i *imlCatalogueModule) Services(ctx context.Context, keyword string) ([]*catalogue_dto.ServiceItem, error) {

	serviceTags, err := i.serviceTagService.List(ctx, nil, nil)
	if err != nil {
		return nil, err
	}
	serviceTagMap := utils.SliceToMapArrayO(serviceTags, func(t *service.Tag) (string, string) {
		return t.Sid, t.Tid
	})

	items, err := i.serviceService.Search(ctx, keyword, nil)
	if err != nil {
		return nil, err
	}
	serviceApiCountMap, err := i.serviceApiService.CountBySids(ctx)
	if err != nil {
		return nil, err
	}
	subscribers, err := i.subscribeService.ListByServices(ctx)
	if err != nil {
		return nil, err
	}
	subscribeCount := map[string]int64{}
	tmp := map[string]struct{}{}
	for _, s := range subscribers {
		if s.ApplyStatus != subscribe.ApplyStatusSubscribe {
			continue
		}
		key := fmt.Sprintf("%s-%s", s.Service, s.Application)
		if _, ok := tmp[key]; !ok {
			tmp[key] = struct{}{}
			subscribeCount[s.Service]++
		}
	}
	result := make([]*catalogue_dto.ServiceItem, 0, len(items))
	for _, v := range items {
		apiNum, ok := serviceApiCountMap[v.Id]
		if !ok || apiNum < 1 {
			continue
		}
		//ps := utils.Intersection(servicePartitionMap[v.Id], projectPartitionMap[v.Project])
		//if len(ps) < 1 {
		//	continue
		//}
		result = append(result, &catalogue_dto.ServiceItem{
			Id:        v.Id,
			Name:      v.Name,
			Tags:      auto.List(serviceTagMap[v.Id]),
			Catalogue: auto.UUID(v.Catalogue),
			//Partition:     auto.List(ps),
			ApiNum:        apiNum,
			SubscriberNum: subscribeCount[v.Id],
			Description:   v.Description,
			Logo:          v.Logo,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].SubscriberNum != result[j].SubscriberNum {
			return result[i].SubscriberNum > result[j].SubscriberNum
		}
		if result[i].ApiNum != result[j].ApiNum {
			return result[i].ApiNum > result[j].ApiNum
		}
		return result[i].Name < result[j].Name
	})
	return result, nil
}

func (i *imlCatalogueModule) recurseUpdateSort(ctx context.Context, parent string, sorts []*catalogue_dto.SortItem) error {
	for index, item := range sorts {
		err := i.catalogueService.Save(ctx, item.Id, &catalogue.EditCatalogue{
			Parent: &parent,
			Sort:   &index,
		})
		if err != nil {
			return err
		}
		if len(item.Children) < 1 {
			continue
		}
		return i.recurseUpdateSort(ctx, item.Id, item.Children)
	}
	return nil
}

func (i *imlCatalogueModule) Sort(ctx context.Context, sorts []*catalogue_dto.SortItem) error {
	return i.transaction.Transaction(ctx, func(ctx context.Context) error {
		err := i.recurseUpdateSort(ctx, "", sorts)
		if err != nil {
			return err
		}
		all, err := i.catalogueService.List(ctx)
		if err != nil {
			return err
		}
		i.root = NewRoot(all)
		return nil
	})

}

func (i *imlCatalogueModule) Search(ctx context.Context, keyword string) ([]*catalogue_dto.Item, error) {
	all, err := i.catalogueService.List(ctx)
	if err != nil {
		return nil, err
	}
	if keyword == "" {
		parentMap := make(map[string][]*catalogue.Catalogue)
		nodeMap := make(map[string]*catalogue.Catalogue)
		for _, v := range all {
			if _, ok := parentMap[v.Parent]; !ok {
				parentMap[v.Parent] = make([]*catalogue.Catalogue, 0)
			}
			parentMap[v.Parent] = append(parentMap[v.Parent], v)
			nodeMap[v.Id] = v
		}
		return treeItems("", parentMap), nil
	}

	catalogues, err := i.catalogueService.Search(ctx, keyword, nil)
	if err != nil {
		return nil, err
	}
	if i.root == nil {
		// 初始化
		i.root = NewRoot(all)
	}
	items := make([]*catalogue_dto.Item, 0, len(catalogues))

	return items, nil
}

func (i *imlCatalogueModule) Create(ctx context.Context, input *catalogue_dto.CreateCatalogue) error {
	parent := ""
	if input.Parent != nil {
		parent = *input.Parent
	}
	if input.Id == "" {
		input.Id = uuid.New().String()
	}
	err := i.catalogueService.Create(ctx, &catalogue.CreateCatalogue{
		Id:     input.Id,
		Name:   input.Name,
		Parent: parent,
		Sort:   _sortMax,
	})
	if err != nil {
		return err
	}
	// 重新初始化
	catalogues, err := i.catalogueService.List(ctx)
	if err != nil {
		return err
	}
	i.root = NewRoot(catalogues)
	return nil
}

func (i *imlCatalogueModule) Edit(ctx context.Context, id string, input *catalogue_dto.EditCatalogue) error {
	err := i.catalogueService.Save(ctx, id, &catalogue.EditCatalogue{
		Name:   input.Name,
		Parent: input.Parent,
	})
	if err != nil {
		return err
	}
	// 重新初始化
	catalogues, err := i.catalogueService.List(ctx)
	if err != nil {
		return err
	}
	i.root = NewRoot(catalogues)
	return nil
}

func (i *imlCatalogueModule) Delete(ctx context.Context, id string) error {
	if id == "" {
		return nil
	}
	list, err := i.catalogueService.Search(ctx, "", map[string]interface{}{
		"parent": id,
	})
	if err != nil {
		return err
	}
	if len(list) > 0 {
		return fmt.Errorf("该目录下存在子目录")
	}
	err = i.catalogueService.Delete(ctx, id)
	if err != nil {
		return err
	}
	// 重新初始化
	catalogues, err := i.catalogueService.List(ctx)
	if err != nil {
		return err
	}
	i.root = NewRoot(catalogues)
	return nil
}

// treeItems 获取子树
func treeItems(parentId string, parentMap map[string][]*catalogue.Catalogue) []*catalogue_dto.Item {
	items := make([]*catalogue_dto.Item, 0)
	if _, ok := parentMap[parentId]; ok {
		for _, v := range parentMap[parentId] {
			childItems := treeItems(v.Id, parentMap)
			items = append(items, &catalogue_dto.Item{
				Id:       v.Id,
				Name:     v.Name,
				Children: childItems,
			})
		}
	}
	return items
}
