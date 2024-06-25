package monitor

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/eolinker/apipark/service/service"

	"github.com/eolinker/apipark/service/subscribe"

	"github.com/eolinker/apipark/service/project"

	"github.com/eolinker/apipark/service/cluster"

	"github.com/eolinker/apipark/module/monitor/driver"

	"github.com/eolinker/go-common/auto"

	"github.com/eolinker/go-common/utils"

	"github.com/eolinker/apipark/service/api"

	"github.com/eolinker/apipark/service/partition"

	"github.com/eolinker/apipark/service/monitor"

	monitor_dto "github.com/eolinker/apipark/module/monitor/dto"
)

var (
	_ IMonitorStatisticModule = (*imlMonitorStatisticModule)(nil)
)

type imlMonitorStatisticModule struct {
	monitorStatisticCacheService monitor.IMonitorStatisticsCache   `autowired:""`
	partitionService             partition.IPartitionService       `autowired:""`
	subscribeService             subscribe.ISubscribeService       `autowired:""`
	serviceService               service.IServiceService           `autowired:""`
	serviceApiService            service.IApiService               `autowired:""`
	projectPartitionService      project.IProjectPartitionsService `autowired:""`
	clusterService               cluster.IClusterService           `autowired:""`
	monitorService               monitor.IMonitorService           `autowired:""`
	projectService               project.IProjectService           `autowired:""`
	apiService                   api.IAPIService                   `autowired:""`
}

func (i *imlMonitorStatisticModule) InvokeTrendWithSubscriberAndApi(ctx context.Context, partitionId string, apiId string, subscriberId string, input *monitor_dto.CommonInput) (*monitor_dto.MonInvokeCountTrend, string, error) {

	wheres, err := i.genCommonWheres(ctx, partitionId)
	if err != nil {
		return nil, "", err
	}
	wheres = append(wheres, monitor.MonWhereItem{
		Key:       "api",
		Operation: "=",
		Values:    []string{apiId},
	}, monitor.MonWhereItem{
		Key:       "app",
		Operation: "=",
		Values:    []string{subscriberId},
	})
	executor, err := i.getExecutor(ctx, partitionId)
	if err != nil {
		return nil, "", err
	}
	result, timeInterval, err := executor.InvokeTrend(ctx, formatTimeByMinute(input.Start), formatTimeByMinute(input.End), wheres)
	if err != nil {
		return nil, "", err
	}
	return monitor_dto.ToMonInvokeCountTrend(result), timeInterval, nil
}

func (i *imlMonitorStatisticModule) InvokeTrendWithProviderAndApi(ctx context.Context, partitionId string, providerId string, apiId string, input *monitor_dto.CommonInput) (*monitor_dto.MonInvokeCountTrend, string, error) {

	wheres, err := i.genCommonWheres(ctx, partitionId)
	if err != nil {
		return nil, "", err
	}
	wheres = append(wheres, monitor.MonWhereItem{
		Key:       "api",
		Operation: "=",
		Values:    []string{apiId},
	}, monitor.MonWhereItem{
		Key:       "provider",
		Operation: "=",
		Values:    []string{providerId},
	})
	executor, err := i.getExecutor(ctx, partitionId)
	if err != nil {
		return nil, "", err
	}
	result, timeInterval, err := executor.InvokeTrend(ctx, formatTimeByMinute(input.Start), formatTimeByMinute(input.End), wheres)
	if err != nil {
		return nil, "", err
	}
	return monitor_dto.ToMonInvokeCountTrend(result), timeInterval, nil
}

func (i *imlMonitorStatisticModule) SubscriberStatisticsOnApi(ctx context.Context, partitionId string, apiId string, input *monitor_dto.StatisticInput) ([]*monitor_dto.ProjectStatisticBasicItem, error) {
	return i.statisticOnApi(ctx, partitionId, apiId, "app", input)
}

func (i *imlMonitorStatisticModule) ProviderStatisticsOnApi(ctx context.Context, partitionId string, apiId string, input *monitor_dto.StatisticInput) ([]*monitor_dto.ProjectStatisticBasicItem, error) {
	return i.statisticOnApi(ctx, partitionId, apiId, "provider", input)
}

func (i *imlMonitorStatisticModule) statisticOnApi(ctx context.Context, partitionId string, apiId string, groupBy string, input *monitor_dto.StatisticInput) ([]*monitor_dto.ProjectStatisticBasicItem, error) {
	_, err := i.partitionService.Get(ctx, partitionId)
	if err != nil {
		return nil, err
	}
	var projects []*project.Project
	switch groupBy {
	case "app":
		projects, err = i.projectService.AppList(ctx)

	case "provider":
		pp, err := i.projectPartitionService.ListByPartition(ctx, partitionId)
		if err != nil {
			return nil, err
		}

		projects, err = i.projectService.List(ctx, utils.SliceToSlice(pp, func(t *project.Partition) string { return t.Project })...)
	default:
		return nil, errors.New("invalid group by")
	}
	if err != nil {
		return nil, err
	}

	wheres, err := i.genCommonWheres(ctx, partitionId)
	if err != nil {
		return nil, err
	}
	wheres = append(wheres, monitor.MonWhereItem{
		Key:       "api",
		Operation: "=",
		Values:    []string{apiId},
	})

	statisticMap, err := i.statistics(ctx, partitionId, groupBy, formatTimeByMinute(input.Start), formatTimeByMinute(input.End), wheres, 0)
	if err != nil {
		return nil, err
	}

	result := make([]*monitor_dto.ProjectStatisticBasicItem, 0, len(statisticMap))
	for _, item := range projects {

		statisticItem := &monitor_dto.ProjectStatisticBasicItem{
			Id:            item.Id,
			Name:          item.Name,
			MonCommonData: new(monitor_dto.MonCommonData),
		}
		if val, ok := statisticMap[item.Id]; ok {
			statisticItem.MonCommonData = monitor_dto.ToMonCommonData(val)
			delete(statisticMap, item.Id)
		}
		result = append(result, statisticItem)
	}
	for key, item := range statisticMap {
		statisticItem := &monitor_dto.ProjectStatisticBasicItem{
			Id:            key,
			Name:          "未知系统-" + key,
			MonCommonData: monitor_dto.ToMonCommonData(item),
		}

		if key == "-" {
			statisticItem.Name = "无系统"
		}
		result = append(result, statisticItem)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].RequestTotal > result[j].RequestTotal
	})
	return result, nil
}

func (i *imlMonitorStatisticModule) ApiStatisticsOnSubscriber(ctx context.Context, partitionId string, subscriberId string, input *monitor_dto.StatisticInput) ([]*monitor_dto.ApiStatisticBasicItem, error) {
	_, err := i.partitionService.Get(ctx, partitionId)
	if err != nil {
		return nil, err
	}
	// 根据订阅ID查询订阅的服务列表
	subscriptions, err := i.subscribeService.MySubscribeServices(ctx, subscriberId, nil, nil)
	if err != nil {
		return nil, err
	}
	serviceIds := utils.SliceToSlice(subscriptions, func(t *subscribe.Subscribe) string {
		return t.Service
	})
	serviceApis, err := i.serviceApiService.List(ctx, serviceIds...)
	if err != nil {
		return nil, err
	}
	serviceApiMap := make(map[string]struct{})
	apiIds := utils.SliceToSlice(serviceApis, func(t *service.Api) string {
		return t.Aid
	}, func(a *service.Api) bool {
		_, ok := serviceApiMap[a.Aid]
		if ok {
			return false
		}
		serviceApiMap[a.Aid] = struct{}{}
		return true
	})

	wheres, err := i.genCommonWheres(ctx, partitionId)
	if err != nil {
		return nil, err
	}
	wheres = append(wheres, monitor.MonWhereItem{
		Key:       "app",
		Operation: "=",
		Values:    []string{subscriberId},
	})

	apiInfos, err := i.apiService.ListInfo(ctx, apiIds...)
	if err != nil {
		return nil, err
	}
	return i.apiStatistics(ctx, partitionId, apiInfos, formatTimeByMinute(input.Start), formatTimeByMinute(input.End), wheres, 0)
}

func (i *imlMonitorStatisticModule) ApiStatisticsOnProvider(ctx context.Context, partitionId string, providerId string, input *monitor_dto.StatisticInput) ([]*monitor_dto.ApiStatisticBasicItem, error) {
	_, err := i.partitionService.Get(ctx, partitionId)
	if err != nil {
		return nil, err
	}

	apiInfos, err := i.apiService.ListInfoForProject(ctx, providerId)
	if err != nil {
		return nil, err
	}
	wheres, err := i.genCommonWheres(ctx, partitionId)
	if err != nil {
		return nil, err
	}
	wheres = append(wheres, monitor.MonWhereItem{
		Key:       "provider",
		Operation: "=",
		Values:    []string{providerId},
	})

	return i.apiStatistics(ctx, partitionId, apiInfos, formatTimeByMinute(input.Start), formatTimeByMinute(input.End), wheres, 0)
}

func (i *imlMonitorStatisticModule) apiStatistics(ctx context.Context, partitionId string, apiInfos []*api.APIInfo, start time.Time, end time.Time, wheres []monitor.MonWhereItem, limit int) ([]*monitor_dto.ApiStatisticBasicItem, error) {
	statisticMap, err := i.statistics(ctx, partitionId, "api", start, end, wheres, limit)
	if err != nil {
		return nil, err
	}

	result := make([]*monitor_dto.ApiStatisticBasicItem, 0, len(statisticMap))
	for _, item := range apiInfos {

		statisticItem := &monitor_dto.ApiStatisticBasicItem{
			Id:            item.UUID,
			Name:          item.Name,
			Path:          item.Path,
			Project:       auto.UUID(item.Project),
			MonCommonData: new(monitor_dto.MonCommonData),
		}
		if val, ok := statisticMap[item.UUID]; ok {
			statisticItem.MonCommonData = monitor_dto.ToMonCommonData(val)
			delete(statisticMap, item.UUID)
		}
		result = append(result, statisticItem)
	}
	for key, item := range statisticMap {
		statisticItem := &monitor_dto.ApiStatisticBasicItem{
			Id:            key,
			Name:          "未知API-" + key,
			MonCommonData: monitor_dto.ToMonCommonData(item),
		}

		if key == "-" {
			statisticItem.Name = "无API"
		}
		result = append(result, statisticItem)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].RequestTotal > result[j].RequestTotal
	})
	return result, nil
}

func (i *imlMonitorStatisticModule) APITrend(ctx context.Context, partitionId string, apiId string, input *monitor_dto.CommonInput) (*monitor_dto.MonInvokeCountTrend, string, error) {

	wheres, err := i.genCommonWheres(ctx, partitionId)
	if err != nil {
		return nil, "", err
	}
	wheres = append(wheres, monitor.MonWhereItem{
		Key:       "api",
		Operation: "=",
		Values:    []string{apiId},
	})
	executor, err := i.getExecutor(ctx, partitionId)
	if err != nil {
		return nil, "", err
	}
	result, timeInterval, err := executor.InvokeTrend(ctx, formatTimeByMinute(input.Start), formatTimeByMinute(input.End), wheres)
	if err != nil {
		return nil, "", err
	}
	return monitor_dto.ToMonInvokeCountTrend(result), timeInterval, nil
}

func (i *imlMonitorStatisticModule) ProviderTrend(ctx context.Context, partitionId string, providerId string, input *monitor_dto.CommonInput) (*monitor_dto.MonInvokeCountTrend, string, error) {

	wheres, err := i.genCommonWheres(ctx, partitionId)
	if err != nil {
		return nil, "", err
	}
	wheres = append(wheres, monitor.MonWhereItem{
		Key:       "provider",
		Operation: "=",
		Values:    []string{providerId},
	})
	executor, err := i.getExecutor(ctx, partitionId)
	if err != nil {
		return nil, "", err
	}
	result, timeInterval, err := executor.InvokeTrend(ctx, formatTimeByMinute(input.Start), formatTimeByMinute(input.End), wheres)
	if err != nil {
		return nil, "", err
	}
	return monitor_dto.ToMonInvokeCountTrend(result), timeInterval, nil
}

func (i *imlMonitorStatisticModule) SubscriberTrend(ctx context.Context, partitionId string, subscriberId string, input *monitor_dto.CommonInput) (*monitor_dto.MonInvokeCountTrend, string, error) {

	wheres, err := i.genCommonWheres(ctx, partitionId)
	if err != nil {
		return nil, "", err
	}
	wheres = append(wheres, monitor.MonWhereItem{
		Key:       "app",
		Operation: "=",
		Values:    []string{subscriberId},
	})
	executor, err := i.getExecutor(ctx, partitionId)
	if err != nil {
		return nil, "", err
	}
	result, timeInterval, err := executor.InvokeTrend(ctx, formatTimeByMinute(input.Start), formatTimeByMinute(input.End), wheres)
	if err != nil {
		return nil, "", err
	}
	return monitor_dto.ToMonInvokeCountTrend(result), timeInterval, nil
}

func (i *imlMonitorStatisticModule) projectStatistics(ctx context.Context, partitionId string, projects []*project.Project, groupBy string, start time.Time, end time.Time, wheres []monitor.MonWhereItem, limit int) ([]*monitor_dto.ProjectStatisticBasicItem, error) {
	statisticMap, err := i.statistics(ctx, partitionId, groupBy, start, end, wheres, limit)
	if err != nil {
		return nil, err
	}

	result := make([]*monitor_dto.ProjectStatisticBasicItem, 0, len(statisticMap))
	for _, item := range projects {
		statisticItem := &monitor_dto.ProjectStatisticBasicItem{
			Id:            item.Id,
			Name:          item.Name,
			MonCommonData: new(monitor_dto.MonCommonData),
		}
		if val, ok := statisticMap[item.Id]; ok {
			statisticItem.MonCommonData = monitor_dto.ToMonCommonData(val)
			delete(statisticMap, item.Id)
		}
		result = append(result, statisticItem)
	}
	for key, item := range statisticMap {
		statisticItem := &monitor_dto.ProjectStatisticBasicItem{
			Id:            key,
			Name:          "未知系统-" + key,
			MonCommonData: monitor_dto.ToMonCommonData(item),
		}

		if key == "-" {
			statisticItem.Name = "无系统"
		}
		result = append(result, statisticItem)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].RequestTotal > result[j].RequestTotal
	})
	return result, nil
}

func (i *imlMonitorStatisticModule) SubscriberStatistics(ctx context.Context, partitionId string, input *monitor_dto.StatisticInput) ([]*monitor_dto.ProjectStatisticBasicItem, error) {
	_, err := i.partitionService.Get(ctx, partitionId)
	if err != nil {
		return nil, err
	}

	projects, err := i.projectService.AppList(ctx, input.Projects...)
	if err != nil {
		return nil, err
	}
	projectIds := utils.SliceToSlice(projects, func(p *project.Project) string {
		return p.Id
	})

	wheres, err := i.genCommonWheres(ctx, partitionId)
	if err != nil {
		return nil, err
	}

	if len(projectIds) > 0 {
		wheres = append(wheres, monitor.MonWhereItem{
			Key:       "app",
			Operation: "in",
			Values:    projectIds,
		})
	}

	return i.projectStatistics(ctx, partitionId, projects, "app", formatTimeByMinute(input.Start), formatTimeByMinute(input.End), wheres, 0)
}

func (i *imlMonitorStatisticModule) projectsByPartition(ctx context.Context, partitionId string, inputProjects []string) ([]*project.Project, error) {
	projectIds := make([]string, 0)
	// 获取当前分区的系统列表
	pp, err := i.projectPartitionService.ListByPartition(ctx, partitionId)
	if err != nil {
		return nil, err
	}
	projectPartitionMap := utils.SliceToMap(pp, func(t *project.Partition) string {
		projectIds = append(projectIds, t.Project)
		return t.Project
	})
	if len(inputProjects) > 0 {
		projectIds = utils.SliceToSlice(inputProjects, func(s string) string {
			return s
		}, func(s string) bool {
			_, ok := projectPartitionMap[s]
			return ok
		})
		if len(projectIds) == 0 {
			return nil, nil
		}
	}

	return i.projectService.List(ctx, projectIds...)
}

func (i *imlMonitorStatisticModule) ProviderStatistics(ctx context.Context, partitionId string, input *monitor_dto.StatisticInput) ([]*monitor_dto.ProjectStatisticBasicItem, error) {
	_, err := i.partitionService.Get(ctx, partitionId)
	if err != nil {
		return nil, err
	}

	projects, err := i.projectsByPartition(ctx, partitionId, input.Projects)
	if err != nil {
		return nil, err
	}

	wheres, err := i.genCommonWheres(ctx, partitionId)
	if err != nil {
		return nil, err
	}

	if len(input.Projects) > 0 {
		wheres = append(wheres, monitor.MonWhereItem{
			Key:       "provider",
			Operation: "in",
			Values:    input.Projects,
		})
	}

	return i.projectStatistics(ctx, partitionId, projects, "provider", formatTimeByMinute(input.Start), formatTimeByMinute(input.End), wheres, 0)
}

func (i *imlMonitorStatisticModule) ApiStatistics(ctx context.Context, partitionId string, input *monitor_dto.StatisticInput) ([]*monitor_dto.ApiStatisticBasicItem, error) {
	_, err := i.partitionService.Get(ctx, partitionId)
	if err != nil {
		return nil, err
	}
	wheres, err := i.genCommonWheres(ctx, partitionId)
	if err != nil {
		return nil, err
	}

	wm := make(map[string]interface{})
	if len(input.Apis) > 0 {
		wm["uuid"] = input.Apis
		wheres = append(wheres, monitor.MonWhereItem{
			Key:       "api",
			Operation: "in",
			Values:    input.Apis,
		})
	}
	if len(input.Projects) > 0 {
		wm["project"] = input.Projects
		wheres = append(wheres, monitor.MonWhereItem{
			Key:       "project",
			Operation: "in",
			Values:    input.Projects,
		})
	}
	// 查询符合条件的API
	apis, err := i.apiService.Search(ctx, input.Path, wm)
	if err != nil {
		return nil, err
	}
	if len(apis) < 1 {
		// 没有符合条件的API
		return make([]*monitor_dto.ApiStatisticBasicItem, 0), nil
	}
	apiIds := utils.SliceToSlice(apis, func(t *api.API) string {
		return t.UUID
	})

	apiInfos, err := i.apiService.ListInfo(ctx, apiIds...)
	if err != nil {
		return nil, err
	}
	return i.apiStatistics(ctx, partitionId, apiInfos, formatTimeByMinute(input.Start), formatTimeByMinute(input.End), wheres, 0)
}

func (i *imlMonitorStatisticModule) MessageTrend(ctx context.Context, partitionId string, input *monitor_dto.CommonInput) (*monitor_dto.MonMessageTrend, string, error) {

	wheres, err := i.genCommonWheres(ctx, partitionId)
	if err != nil {
		return nil, "", err
	}
	executor, err := i.getExecutor(ctx, partitionId)
	if err != nil {
		return nil, "", err
	}
	result, timeInterval, err := executor.MessageTrend(ctx, formatTimeByMinute(input.Start), formatTimeByMinute(input.End), wheres)
	if err != nil {
		return nil, "", err
	}
	return monitor_dto.ToMonMessageTrend(result), timeInterval, nil
}

func (i *imlMonitorStatisticModule) InvokeTrend(ctx context.Context, partitionId string, input *monitor_dto.CommonInput) (*monitor_dto.MonInvokeCountTrend, string, error) {
	wheres, err := i.genCommonWheres(ctx, partitionId)
	if err != nil {
		return nil, "", err
	}
	executor, err := i.getExecutor(ctx, partitionId)
	if err != nil {
		return nil, "", err
	}
	result, timeInterval, err := executor.InvokeTrend(ctx, formatTimeByMinute(input.Start), formatTimeByMinute(input.End), wheres)
	if err != nil {
		return nil, "", err
	}
	return monitor_dto.ToMonInvokeCountTrend(result), timeInterval, nil
}

func (i *imlMonitorStatisticModule) genCommonWheres(ctx context.Context, partitionId string) ([]monitor.MonWhereItem, error) {

	info, err := i.partitionService.Get(ctx, partitionId)
	if err != nil {
		return nil, err
	}

	wheres := make([]monitor.MonWhereItem, 0, 1)

	wheres = append(wheres, monitor.MonWhereItem{
		Key:       "cluster",
		Operation: "=",
		Values:    []string{info.Cluster},
	})

	return wheres, nil
}

func (i *imlMonitorStatisticModule) statistics(ctx context.Context, partitionId string, groupBy string, start, end time.Time, wheres []monitor.MonWhereItem, limit int) (map[string]monitor.MonCommonData, error) {
	statisticMap, _ := i.monitorStatisticCacheService.GetStatisticsCache(ctx, partitionId, start, end, groupBy, wheres, limit)
	if len(statisticMap) > 0 {
		return statisticMap, nil
	}

	executor, err := i.getExecutor(ctx, partitionId)
	if err != nil {
		return nil, err
	}

	result, err := executor.CommonStatistics(ctx, start, end, groupBy, limit, wheres)
	if err != nil {
		return nil, err
	}
	i.monitorStatisticCacheService.SetStatisticsCache(ctx, partitionId, start, end, groupBy, wheres, limit, result)
	return result, nil
}

func (i *imlMonitorStatisticModule) TopAPIStatistics(ctx context.Context, partitionId string, limit int, input *monitor_dto.CommonInput) ([]*monitor_dto.ApiStatisticItem, error) {

	wheres, err := i.genCommonWheres(ctx, partitionId)
	if err != nil {
		return nil, err
	}

	statisticMap, err := i.statistics(ctx, partitionId, "api", formatTimeByMinute(input.Start), formatTimeByMinute(input.End), wheres, limit)
	if err != nil {
		return nil, err
	}

	uuids := utils.MapToSlice(statisticMap, func(key string, value monitor.MonCommonData) string {
		return value.ID
	})
	apis, err := i.apiService.ListInfo(ctx, uuids...)
	if err != nil {
		return nil, err
	}
	apiMap := utils.SliceToMap(apis, func(t *api.APIInfo) string {
		return t.UUID
	})
	result := make([]*monitor_dto.ApiStatisticItem, 0, len(statisticMap))
	for key, item := range statisticMap {
		statisticItem := &monitor_dto.ApiStatisticItem{
			ApiStatisticBasicItem: &monitor_dto.ApiStatisticBasicItem{
				Id:            key,
				MonCommonData: monitor_dto.ToMonCommonData(item),
			},
		}
		if a, ok := apiMap[item.ID]; ok {
			statisticItem.Name = a.Name
			statisticItem.Path = a.Path
			statisticItem.Project = auto.UUID(a.Project)
		} else {
			statisticItem.IsRed = true
			if key == "-" {
				statisticItem.Name = "无API"
			} else {
				statisticItem.Name = fmt.Sprintf("未知API-%s", key)
			}
		}
		result = append(result, statisticItem)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].RequestTotal > result[j].RequestTotal
	})
	return result, nil
}

func (i *imlMonitorStatisticModule) TopSubscriberStatistics(ctx context.Context, partitionId string, limit int, input *monitor_dto.CommonInput) ([]*monitor_dto.ProjectStatisticItem, error) {
	_, err := i.partitionService.Get(ctx, partitionId)
	if err != nil {
		return nil, err
	}
	return i.topProjectStatistics(ctx, partitionId, "app", input, limit)
}

func (i *imlMonitorStatisticModule) TopProviderStatistics(ctx context.Context, partitionId string, limit int, input *monitor_dto.CommonInput) ([]*monitor_dto.ProjectStatisticItem, error) {
	_, err := i.partitionService.Get(ctx, partitionId)
	if err != nil {
		return nil, err
	}
	return i.topProjectStatistics(ctx, partitionId, "provider", input, limit)
}
func (i *imlMonitorStatisticModule) topProjectStatistics(ctx context.Context, partitionId string, groupBy string, input *monitor_dto.CommonInput, limit int) ([]*monitor_dto.ProjectStatisticItem, error) {
	wheres, err := i.genCommonWheres(ctx, partitionId)
	if err != nil {
		return nil, err
	}
	statisticMap, err := i.statistics(ctx, partitionId, groupBy, formatTimeByMinute(input.Start), formatTimeByMinute(input.End), wheres, limit)
	if err != nil {
		return nil, err
	}
	var projects []*project.Project
	switch groupBy {
	case "app":
		projects, err = i.projectService.AppList(ctx)
	case "provider":
		projects, err = i.projectsByPartition(ctx, partitionId, nil)
	default:
		return nil, errors.New("invalid group by")
	}
	if err != nil {
		return nil, err
	}
	projectMap := utils.SliceToMap(projects, func(t *project.Project) string {
		return t.Id
	})
	result := make([]*monitor_dto.ProjectStatisticItem, 0, len(statisticMap))
	for key, item := range statisticMap {
		statisticItem := &monitor_dto.ProjectStatisticItem{
			ProjectStatisticBasicItem: &monitor_dto.ProjectStatisticBasicItem{
				Id:            key,
				MonCommonData: monitor_dto.ToMonCommonData(item),
			},
		}
		if a, ok := projectMap[item.ID]; ok {
			statisticItem.Name = a.Name
		} else {
			statisticItem.IsRed = true
			if key == "-" {
				statisticItem.Name = "无系统"
			} else {
				statisticItem.Name = fmt.Sprintf("未知系统-%s", key)
			}
		}
		result = append(result, statisticItem)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].RequestTotal > result[j].RequestTotal
	})
	return result, nil
}

func (i *imlMonitorStatisticModule) getExecutor(ctx context.Context, partitionId string) (driver.IExecutor, error) {
	info, err := i.monitorService.GetByPartition(ctx, partitionId)
	if err != nil {
		return nil, err
	}
	return driver.CreateExecutor(info.Driver, info.Config)
}

func (i *imlMonitorStatisticModule) RequestSummary(ctx context.Context, partitionId string, input *monitor_dto.CommonInput) (*monitor_dto.MonSummaryOutput, error) {
	wheres, err := i.genCommonWheres(ctx, partitionId)
	if err != nil {
		return nil, err
	}
	executor, err := i.getExecutor(ctx, partitionId)
	if err != nil {
		return nil, err

	}
	summary, err := executor.RequestSummary(ctx, formatTimeByMinute(input.Start), formatTimeByMinute(input.End), wheres)
	if err != nil {
		return nil, err
	}

	return monitor_dto.ToMonSummaryOutput(summary), nil
}

func (i *imlMonitorStatisticModule) ProxySummary(ctx context.Context, partitionId string, input *monitor_dto.CommonInput) (*monitor_dto.MonSummaryOutput, error) {
	wheres, err := i.genCommonWheres(ctx, partitionId)
	if err != nil {
		return nil, err
	}
	executor, err := i.getExecutor(ctx, partitionId)
	if err != nil {
		return nil, err

	}
	summary, err := executor.ProxySummary(ctx, formatTimeByMinute(input.Start), formatTimeByMinute(input.End), wheres)
	if err != nil {
		return nil, err
	}

	return monitor_dto.ToMonSummaryOutput(summary), nil
}
