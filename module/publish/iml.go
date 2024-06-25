package publish

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/eolinker/apipark/service/partition"

	"github.com/eolinker/apipark/service/universally/commit"

	"github.com/eolinker/apipark/service/api"
	"github.com/eolinker/apipark/service/upstream"

	"github.com/eolinker/go-common/auto"

	"github.com/eolinker/apipark/gateway"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/apipark/service/project"

	projectDiff "github.com/eolinker/apipark/module/project_diff"
	"github.com/eolinker/apipark/module/publish/dto"
	"github.com/eolinker/apipark/service/cluster"
	"github.com/eolinker/apipark/service/publish"
	"github.com/eolinker/apipark/service/release"
	"github.com/eolinker/go-common/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	_                     IPublishModule = (*imlPublishModule)(nil)
	projectRuleMustServer                = map[string]bool{
		"as_server": true,
	}
)

type imlPublishModule struct {
	projectDiffModule       projectDiff.IProjectDiffModule    `autowired:""`
	publishService          publish.IPublishService           `autowired:""`
	apiService              api.IAPIService                   `autowired:""`
	upstreamService         upstream.IUpstreamService         `autowired:""`
	releaseService          release.IReleaseService           `autowired:""`
	projectPartitionService project.IProjectPartitionsService `autowired:""`
	clusterService          cluster.IClusterService           `autowired:""`
	partitionService        partition.IPartitionService       `autowired:""`
	projectService          project.IProjectService           `autowired:""`
}

func (m *imlPublishModule) initGateway(ctx context.Context, partitionId string, clientDriver gateway.IClientDriver) error {

	projectPartitions, err := m.projectPartitionService.ListByPartition(ctx, partitionId)
	if err != nil {
		return err
	}
	projectIds := utils.SliceToSlice(projectPartitions, func(p *project.Partition) string {
		return p.Project
	})
	for _, projectId := range projectIds {
		releaseInfo, err := m.getProjectRelease(ctx, projectId, partitionId)
		if err != nil {
			return err
		}
		if releaseInfo == nil {
			continue
		}

		err = clientDriver.Project().Online(ctx, releaseInfo)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *imlPublishModule) getProjectRelease(ctx context.Context, projectID string, partitionId string) (*gateway.ProjectRelease, error) {

	releaseInfo, err := m.releaseService.GetRunning(ctx, projectID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, nil
	}
	commits, err := m.releaseService.GetCommits(ctx, releaseInfo.UUID)
	if err != nil {
		return nil, err
	}
	apiIds := make([]string, 0, len(commits))
	apiProxyCommitIds := make([]string, 0, len(commits))
	upstreamCommitIds := make([]string, 0, len(commits))
	for _, c := range commits {
		switch c.Type {
		case release.CommitApiProxy:
			apiIds = append(apiIds, c.Target)
			apiProxyCommitIds = append(apiProxyCommitIds, c.Commit)
		case release.CommitUpstream:
			upstreamCommitIds = append(upstreamCommitIds, c.Commit)
		}
	}

	apiInfos, err := m.apiService.ListInfo(ctx, apiIds...)
	if err != nil {
		return nil, err
	}

	proxyCommits, err := m.apiService.ListProxyCommit(ctx, apiProxyCommitIds...)
	if err != nil {
		return nil, err
	}
	proxyCommitMap := utils.SliceToMapO(proxyCommits, func(c *commit.Commit[api.Proxy]) (string, *api.Proxy) {
		return c.Target, c.Data
	})

	upstreamCommits, err := m.upstreamService.ListCommit(ctx, upstreamCommitIds...)
	if err != nil {
		return nil, err
	}
	version := releaseInfo.UUID
	apis := make([]*gateway.ApiRelease, 0, len(apiInfos))
	for _, a := range apiInfos {
		apiInfo := &gateway.ApiRelease{
			BasicItem: &gateway.BasicItem{
				ID:          a.UUID,
				Description: a.Description,
				Version:     version,
			},
			Path:    a.Path,
			Method:  []string{a.Method},
			Service: a.Upstream,
		}
		proxy, ok := proxyCommitMap[a.UUID]
		if ok {
			apiInfo.ProxyPath = proxy.Path
			apiInfo.ProxyHeaders = utils.SliceToSlice(proxy.Headers, func(h *api.Header) *gateway.ProxyHeader {
				return &gateway.ProxyHeader{
					Key:   h.Key,
					Value: h.Value,
				}
			})
			apiInfo.Retry = proxy.Retry
			apiInfo.Timeout = proxy.Timeout
		}
		apis = append(apis, apiInfo)
	}
	var upstreamRelease *gateway.UpstreamRelease
	for _, c := range upstreamCommits {
		if c.Key != partitionId {
			continue
		}
		upstreamRelease = &gateway.UpstreamRelease{
			BasicItem: &gateway.BasicItem{
				ID:      c.Target,
				Version: version,
				MatchLabels: map[string]string{
					"project": projectID,
				},
			},
			PassHost: c.Data.PassHost,
			Scheme:   c.Data.Scheme,
			Balance:  c.Data.Balance,
			Timeout:  c.Data.Timeout,
			Nodes: utils.SliceToSlice(c.Data.Nodes, func(n *upstream.NodeConfig) string {
				return fmt.Sprintf("%s weight=%d", n.Address, n.Weight)
			}),
		}
	}

	return &gateway.ProjectRelease{
		Id:       projectID,
		Version:  version,
		Apis:     apis,
		Upstream: upstreamRelease,
	}, nil
}

func (m *imlPublishModule) getReleaseInfo(ctx context.Context, projectID, releaseId string, version string, partitionIds []string) (map[string]*gateway.ProjectRelease, error) {
	commits, err := m.releaseService.GetCommits(ctx, releaseId)
	if err != nil {
		return nil, err
	}
	apiIds := make([]string, 0, len(commits))
	apiProxyCommitIds := make([]string, 0, len(commits))
	upstreamCommitIds := make([]string, 0, len(commits))
	for _, c := range commits {
		switch c.Type {
		case release.CommitApiProxy:
			apiIds = append(apiIds, c.Target)
			apiProxyCommitIds = append(apiProxyCommitIds, c.Commit)
		case release.CommitUpstream:
			upstreamCommitIds = append(upstreamCommitIds, c.Commit)
		}
	}

	apiInfos, err := m.apiService.ListInfo(ctx, apiIds...)
	if err != nil {
		return nil, err
	}

	proxyCommits, err := m.apiService.ListProxyCommit(ctx, apiProxyCommitIds...)
	if err != nil {
		return nil, err
	}
	proxyCommitMap := utils.SliceToMapO(proxyCommits, func(c *commit.Commit[api.Proxy]) (string, *api.Proxy) {
		return c.Target, c.Data
	})

	upstreamCommits, err := m.upstreamService.ListCommit(ctx, upstreamCommitIds...)
	if err != nil {
		return nil, err
	}
	apis := make([]*gateway.ApiRelease, 0, len(apiInfos))
	for _, a := range apiInfos {
		apiInfo := &gateway.ApiRelease{
			BasicItem: &gateway.BasicItem{
				ID:          a.UUID,
				Description: a.Description,
				Version:     version,
			},
			Path:    a.Path,
			Method:  []string{a.Method},
			Service: a.Upstream,
		}
		proxy, ok := proxyCommitMap[a.UUID]
		if ok {
			apiInfo.Plugins = utils.MapChange(proxy.Plugins, func(v api.PluginSetting) *gateway.Plugin {
				return &gateway.Plugin{
					Config:  v.Config,
					Disable: v.Disable,
				}
			})
			apiInfo.Extends = proxy.Extends
			apiInfo.ProxyPath = proxy.Path
			apiInfo.ProxyHeaders = utils.SliceToSlice(proxy.Headers, func(h *api.Header) *gateway.ProxyHeader {
				return &gateway.ProxyHeader{
					Key:   h.Key,
					Value: h.Value,
				}
			})
			apiInfo.Retry = proxy.Retry
			apiInfo.Timeout = proxy.Timeout
		}
		apis = append(apis, apiInfo)
	}
	projectReleaseMap := make(map[string]*gateway.ProjectRelease)
	upstreamReleaseMap := make(map[string]*gateway.UpstreamRelease)

	for _, c := range upstreamCommits {
		for _, partitionId := range partitionIds {
			upstreamRelease := &gateway.UpstreamRelease{
				BasicItem: &gateway.BasicItem{
					ID:      c.Target,
					Version: version,
					MatchLabels: map[string]string{
						"project": projectID,
					},
				},
				PassHost: c.Data.PassHost,
				Scheme:   c.Data.Scheme,
				Balance:  c.Data.Balance,
				Timeout:  c.Data.Timeout,
				Nodes: utils.SliceToSlice(c.Data.Nodes, func(n *upstream.NodeConfig) string {
					return fmt.Sprintf("%s weight=%d", n.Address, n.Weight)
				}),
			}

			upstreamReleaseMap[partitionId] = upstreamRelease
		}
	}

	for _, clusterId := range partitionIds {
		projectReleaseMap[clusterId] = &gateway.ProjectRelease{
			Id:       projectID,
			Version:  version,
			Apis:     apis,
			Upstream: upstreamReleaseMap[clusterId],
		}
	}
	return projectReleaseMap, nil
}

func (m *imlPublishModule) PublishStatuses(ctx context.Context, project string, id string) ([]*dto.PublishStatus, error) {
	_, err := m.projectService.CheckProject(ctx, project, projectRuleMustServer)
	if err != nil {
		return nil, err
	}
	flow, err := m.publishService.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if flow.Project != project {
		return nil, errors.New("项目不一致")
	}
	list, err := m.publishService.GetPublishStatus(ctx, id)
	if err != nil {
		return nil, err
	}
	return utils.SliceToSlice(list, func(s *publish.Status) *dto.PublishStatus {
		status := s.Status
		errMsg := s.Error
		if s.Status == publish.StatusPublishing && time.Now().Sub(s.UpdateAt) > 30*time.Second {
			status = publish.StatusPublishError
			errMsg = "发布超时"
		}
		return &dto.PublishStatus{
			Partition: auto.UUID(s.Partition),
			Cluster:   auto.UUID(s.Cluster),
			Status:    status.String(),
			Error:     errMsg,
		}

	}), nil
}

// Apply applies the changes to the imlPublishModule.
//
// ctx context.Context, project string, input *dto.ApplyInput
// *dto.Publish, error
func (m *imlPublishModule) Apply(ctx context.Context, project string, input *dto.ApplyInput) (*dto.Publish, error) {
	_, err := m.projectService.CheckProject(ctx, project, projectRuleMustServer)
	if err != nil {
		return nil, err
	}
	err = m.checkPublish(ctx, project, input.Release)
	if err != nil {
		return nil, err
	}

	previous := ""
	running, err := m.releaseService.GetRunning(ctx, project)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {

		return nil, err
	}
	if running != nil {
		previous = running.UUID
	}

	releaseToPublish, err := m.releaseService.GetRelease(ctx, input.Release)
	if err != nil {
		// 目标版本不存在
		return nil, err
	}

	newPublishId := uuid.NewString()
	diff, ok, err := m.projectDiffModule.DiffForLatest(ctx, project, previous)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("latest completeness check failed")
	}
	err = m.publishService.Create(ctx, newPublishId, project, releaseToPublish.UUID, previous, releaseToPublish.Version, input.Remark, diff)
	if err != nil {
		return nil, err
	}
	np, err := m.publishService.Get(ctx, newPublishId)
	if err != nil {
		return nil, err
	}
	return dto.FromModel(np, releaseToPublish.Remark), nil
}

func (m *imlPublishModule) CheckPublish(ctx context.Context, project string, releaseId string) (*dto.DiffOut, error) {
	_, err := m.projectService.CheckProject(ctx, project, projectRuleMustServer)
	if err != nil {
		return nil, err
	}
	err = m.checkPublish(ctx, project, releaseId)
	if err != nil {
		return nil, err
	}

	running, err := m.releaseService.GetRunning(ctx, project)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	runningReleaseId := ""
	if running != nil {
		runningReleaseId = running.UUID
	}
	if releaseId == "" {
		// 发布latest 版本
		diff, _, err := m.projectDiffModule.DiffForLatest(ctx, project, runningReleaseId)
		if err != nil {
			return nil, err
		}
		return m.projectDiffModule.Out(ctx, diff)
	} else {
		// 发布 releaseId 版本, 返回 与当前版本的差异
		diff, err := m.projectDiffModule.Diff(ctx, project, runningReleaseId, releaseId)
		if err != nil {
			return nil, err
		}
		return m.projectDiffModule.Out(ctx, diff)
	}

}
func (m *imlPublishModule) checkPublish(ctx context.Context, project string, releaseId string) error {
	flows, err := m.publishService.ListForStatus(ctx, project, publish.StatusApply, publish.StatusAccept)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if len(flows) > 0 {
		return errors.New("正在发布中")
	}
	running, err := m.releaseService.GetRunning(ctx, project)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if running == nil {
		return nil
	}
	if running.UUID == releaseId {
		return errors.New("不能申请发布当前版本")
	}
	return nil
}
func (m *imlPublishModule) Close(ctx context.Context, project, id string) error {
	err := m.publishService.SetStatus(ctx, project, id, publish.StatusClose)
	if err != nil {
		return err
	}

	return nil
}

func (m *imlPublishModule) Stop(ctx context.Context, project string, id string) error {
	_, err := m.projectService.CheckProject(ctx, project, projectRuleMustServer)
	if err != nil {
		return err
	}
	flow, err := m.publishService.Get(ctx, id)
	if err != nil {
		return err
	}
	if flow.Project != project {
		return errors.New("项目不一致")
	}

	if flow.Status != publish.StatusApply && flow.Status != publish.StatusAccept {
		return errors.New("只有发布中状态才能停止")
	}
	status := publish.StatusStop
	if flow.Status == publish.StatusApply {
		status = publish.StatusClose
	}
	return m.publishService.SetStatus(ctx, project, id, status)
}

func (m *imlPublishModule) Refuse(ctx context.Context, project string, id string, commits string) error {
	_, err := m.projectService.CheckProject(ctx, project, projectRuleMustServer)
	if err != nil {
		return err
	}
	return m.publishService.Refuse(ctx, project, id, commits)
}

func (m *imlPublishModule) Accept(ctx context.Context, project string, id string, commits string) error {
	_, err := m.projectService.CheckProject(ctx, project, projectRuleMustServer)
	if err != nil {
		return err
	}
	return m.publishService.Accept(ctx, project, id, commits)
}

func (m *imlPublishModule) publish(ctx context.Context, id string, partitionId string, projectRelease *gateway.ProjectRelease) error {
	pInfo, err := m.partitionService.Get(ctx, partitionId)
	if err != nil {
		return err
	}
	publishStatus := &publish.Status{
		Publish:   id,
		Partition: partitionId,
		Status:    publish.StatusPublishing,
		UpdateAt:  time.Now(),
	}
	err = m.publishService.SetPublishStatus(ctx, publishStatus)
	if err != nil {
		return fmt.Errorf("set publishing publishStatus error: %v", err)
	}
	defer func() {
		err := m.publishService.SetPublishStatus(ctx, publishStatus)
		if err != nil {
			log.Errorf("set publishing publishStatus error: %v", err)
		}
	}()

	client, err := m.clusterService.GatewayClient(ctx, pInfo.Cluster)
	if err != nil {
		publishStatus.Status = publish.StatusPublishError
		publishStatus.Error = err.Error()
		publishStatus.UpdateAt = time.Now()
		return fmt.Errorf("get gateway client error: %v", err)
	}
	defer func() {
		err := client.Close(ctx)
		if err != nil {
			log.Warn("close apinto client:", err)
		}
	}()
	err = client.Project().Online(ctx, projectRelease)
	if err != nil {
		publishStatus.Status = publish.StatusPublishError
		publishStatus.Error = err.Error()
		publishStatus.UpdateAt = time.Now()
		return fmt.Errorf("online error: %v", err)
	}
	publishStatus.Status = publish.StatusDone
	publishStatus.UpdateAt = time.Now()
	return nil
}

func (m *imlPublishModule) Publish(ctx context.Context, project string, id string) error {
	_, err := m.projectService.CheckProject(ctx, project, projectRuleMustServer)
	if err != nil {
		return err
	}
	flow, err := m.publishService.Get(ctx, id)
	if err != nil {
		return err
	}
	if flow.Project != project {
		return errors.New("项目不一致")
	}
	if flow.Status != publish.StatusAccept {
		return errors.New("只有通过状态才能发布")
	}
	partitionIds, err := m.projectPartitionService.GetByProject(ctx, project)
	if err != nil {
		return err
	}
	projectReleaseMap, err := m.getReleaseInfo(ctx, project, flow.Release, flow.Release, partitionIds)
	if err != nil {
		return err
	}
	hasError := false
	for _, partitionId := range partitionIds {
		err = m.publish(ctx, flow.Id, partitionId, projectReleaseMap[partitionId])
		if err != nil {
			hasError = true
			log.Error(err)
			continue
		}
	}
	err = m.releaseService.SetRunning(ctx, project, flow.Release)
	if err != nil {
		return err
	}
	status := publish.StatusDone
	if hasError {
		status = publish.StatusPublishError
	}
	return m.publishService.SetStatus(ctx, project, id, status)
}

func (m *imlPublishModule) List(ctx context.Context, project string, page, pageSize int) ([]*dto.Publish, int64, error) {
	_, err := m.projectService.CheckProject(ctx, project, projectRuleMustServer)
	if err != nil {
		return nil, 0, err
	}
	list, total, err := m.publishService.ListProjectPage(ctx, project, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return utils.SliceToSlice(list, func(s *publish.Publish) *dto.Publish {
		return dto.FromModel(s, "")
	}), total, nil
}

func (m *imlPublishModule) Detail(ctx context.Context, project string, id string) (*dto.PublishDetail, error) {
	_, err := m.projectService.CheckProject(ctx, project, projectRuleMustServer)
	if err != nil {
		return nil, err
	}
	flow, err := m.publishService.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if flow.Project != project {
		return nil, errors.New("项目不一致")
	}
	diff, err := m.publishService.GetDiff(ctx, id)
	if err != nil {
		return nil, err
	}
	out, err := m.projectDiffModule.Out(ctx, diff)
	if err != nil {
		return nil, err
	}
	publishStatuses, err := m.PublishStatuses(ctx, project, id)
	if err != nil {
		return nil, err
	}
	releaseInfo, err := m.releaseService.GetRelease(ctx, flow.Release)
	if err != nil {
		return nil, err
	}
	return &dto.PublishDetail{
		Publish:         dto.FromModel(flow, releaseInfo.Remark),
		Diffs:           out,
		PublishStatuses: publishStatuses,
	}, nil

}
