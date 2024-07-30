package project_diff

import (
	"context"
	"errors"
	"fmt"

	"github.com/eolinker/apipark/service/api"
	"github.com/eolinker/apipark/service/cluster"
	"github.com/eolinker/apipark/service/project_diff"
	"github.com/eolinker/apipark/service/release"
	"github.com/eolinker/apipark/service/universally/commit"
	"github.com/eolinker/apipark/service/upstream"
	"github.com/eolinker/go-common/auto"
	"github.com/eolinker/go-common/utils"
)

type imlProjectDiff struct {
	apiService      api.IAPIService           `autowired:""`
	upstreamService upstream.IUpstreamService `autowired:""`
	releaseService  release.IReleaseService   `autowired:""`
	clusterService  cluster.IClusterService   `autowired:""`
}

func (m *imlProjectDiff) Diff(ctx context.Context, projectId string, baseRelease, targetRelease string) (*project_diff.Diff, error) {
	if targetRelease == "" {
		return nil, fmt.Errorf("target release is required")
	}

	var target *projectInfo

	targetReleaseValue, err := m.releaseService.GetRelease(ctx, targetRelease)
	if err != nil {
		return nil, fmt.Errorf("get target release  failed:%w", err)
	}
	if targetReleaseValue.Project != projectId {
		return nil, errors.New("project not match")
	}

	target, err = m.getReleaseInfo(ctx, targetRelease)
	if err != nil {
		return nil, err
	}
	base, err := m.getBaseInfo(ctx, projectId, baseRelease)
	if err != nil {
		return nil, err
	}
	target.id = projectId
	clusters, err := m.clusterService.List(ctx)
	if err != nil {
		return nil, err
	}
	clusterIds := utils.SliceToSlice(clusters, func(i *cluster.Cluster) string {
		return i.Uuid
	})
	diff := m.diff(clusterIds, base, target)
	return diff, nil

}
func (m *imlProjectDiff) getBaseInfo(ctx context.Context, projectId, baseRelease string) (*projectInfo, error) {
	if baseRelease == "" {
		return &projectInfo{}, nil
	}
	baseReleaseValue, err := m.releaseService.GetRelease(ctx, baseRelease)
	if err != nil {
		return nil, fmt.Errorf("get base release  failed:%w", err)
	}
	if baseReleaseValue.Project != projectId {
		return nil, errors.New("project not match")
	}
	base, err := m.getReleaseInfo(ctx, baseRelease)
	if err != nil {
		return nil, fmt.Errorf("get base release info failed:%w", err)
	}

	return base, nil
}
func (m *imlProjectDiff) DiffForLatest(ctx context.Context, projectId string, baseRelease string) (*project_diff.Diff, bool, error) {

	apis, err := m.apiService.ListForProject(ctx, projectId)
	if err != nil {
		return nil, false, err
	}

	apiIds := utils.SliceToSlice(apis, func(i *api.API) string {
		return i.UUID
	})
	apiInfos, err := m.apiService.ListInfo(ctx, apiIds...)
	if err != nil {
		return nil, false, err
	}
	proxy, err := m.apiService.ListLatestCommitProxy(ctx, apiIds...)
	if err != nil {
		return nil, false, fmt.Errorf("diff for api commit %v", err)
	}
	documents, err := m.apiService.ListLatestCommitDocument(ctx, apiIds...)
	if err != nil {
		return nil, false, err
	}

	upstreamCommits, err := m.upstreamService.ListLatestCommit(ctx, projectId)
	if err != nil {
		return nil, false, err
	}

	base, err := m.getBaseInfo(ctx, projectId, baseRelease)
	if err != nil {
		return nil, false, err
	}
	target := &projectInfo{
		id:              projectId,
		apis:            apiInfos,
		apiCommits:      proxy,
		apiDocs:         documents,
		upstreamCommits: upstreamCommits,
	}
	clusters, err := m.clusterService.List(ctx)
	if err != nil {
		return nil, false, err
	}
	clusterIds := utils.SliceToSlice(clusters, func(i *cluster.Cluster) string {
		return i.Uuid
	})
	return m.diff(clusterIds, base, target), true, nil
}
func (m *imlProjectDiff) getReleaseInfo(ctx context.Context, releaseId string) (*projectInfo, error) {
	commits, err := m.releaseService.GetCommits(ctx, releaseId)
	if err != nil {
		return nil, err
	}

	apiIds := utils.SliceToSlice(commits, func(i *release.ProjectCommits) string {
		return i.Target
	}, func(c *release.ProjectCommits) bool {
		return c.Type == release.CommitApiProxy || c.Type == release.CommitApiDocument
	})
	apiInfos, err := m.apiService.ListInfo(ctx, apiIds...)
	if err != nil {
		return nil, err
	}
	apiProxyCommitIds := utils.SliceToSlice(commits, func(i *release.ProjectCommits) string {
		return i.Commit
	}, func(c *release.ProjectCommits) bool {
		return c.Type == release.CommitApiProxy
	})
	apiDocumentCommitIds := utils.SliceToSlice(commits, func(i *release.ProjectCommits) string {
		return i.Commit
	}, func(c *release.ProjectCommits) bool {
		return c.Type == release.CommitApiDocument
	})
	upstreamCommitIds := utils.SliceToSlice(commits, func(i *release.ProjectCommits) string {
		return i.Commit
	}, func(c *release.ProjectCommits) bool {
		return c.Type == release.CommitUpstream
	})
	proxyCommits, err := m.apiService.ListProxyCommit(ctx, apiProxyCommitIds...)
	if err != nil {
		return nil, err
	}
	documentCommits, err := m.apiService.ListDocumentCommit(ctx, apiDocumentCommitIds...)
	if err != nil {
		return nil, err
	}
	upstreamCommits, err := m.upstreamService.ListCommit(ctx, upstreamCommitIds...)
	if err != nil {
		return nil, err
	}
	return &projectInfo{
		apis:            apiInfos,
		apiCommits:      proxyCommits,
		apiDocs:         documentCommits,
		upstreamCommits: upstreamCommits,
	}, nil
}
func (m *imlProjectDiff) diff(partitions []string, base, target *projectInfo) *project_diff.Diff {
	out := &project_diff.Diff{
		Apis:      nil,
		Upstreams: nil,
		//Clusters: partitions,
	}
	baseApis := utils.NewSet(utils.SliceToSlice(base.apis, func(i *api.APIInfo) string {
		return i.UUID
	})...)
	baseApiProxy := utils.SliceToMap(base.apiCommits, func(i *commit.Commit[api.Proxy]) string {
		return i.Target
	})
	baseAPIDoc := utils.SliceToMap(base.apiDocs, func(i *commit.Commit[api.Document]) string {
		return i.Target
	})

	targetApiProxy := utils.SliceToMap(target.apiCommits, func(i *commit.Commit[api.Proxy]) string {
		return i.Target
	})
	targetAPIDoc := utils.SliceToMap(target.apiDocs, func(i *commit.Commit[api.Document]) string {
		return i.Target
	})

	for _, apiInfo := range target.apis {
		apiId := apiInfo.UUID
		a := &project_diff.ApiDiff{
			APi:    apiInfo.UUID,
			Name:   apiInfo.Name,
			Method: apiInfo.Method,
			Path:   apiInfo.Path,
			Status: project_diff.Status{},
		}

		pc, hasPc := targetApiProxy[apiId]
		dc, hasDC := targetAPIDoc[apiId]
		if !hasPc {
			// 未设置proxy信息
			a.Status.Proxy = project_diff.StatusUnset
		}
		if !hasDC {
			// 未设置文档
			a.Status.Doc = project_diff.StatusUnset
		}

		if !baseApis.Has(apiId) {
			a.Change = project_diff.ChangeTypeNew
		} else {
			a.Change = project_diff.ChangeTypeNone

			baseProxy, hasBaseProxy := baseApiProxy[apiId]
			baseDoc, hasBaseDoc := baseAPIDoc[apiId]
			if hasBaseDoc != hasDC || hasBaseProxy != hasPc {
				// 文档或者proxy变更
				a.Change = project_diff.ChangeTypeUpdate
			} else if (hasPc && pc.UUID != baseProxy.UUID) || (hasDC && dc.UUID != baseDoc.UUID) {
				// 文档 或者 proxy 变更
				a.Change = project_diff.ChangeTypeUpdate
			}
		}
		out.Apis = append(out.Apis, a)

	}
	baseApis.Remove(utils.SliceToSlice(out.Apis, func(i *project_diff.ApiDiff) string {
		return i.APi
	})...)
	for _, apiInfo := range base.apis {
		if baseApis.Has(apiInfo.UUID) {
			out.Apis = append(out.Apis, &project_diff.ApiDiff{
				APi:    apiInfo.UUID,
				Name:   apiInfo.Name,
				Method: apiInfo.Method,
				Path:   apiInfo.Path,
				Status: project_diff.Status{},
				Change: project_diff.ChangeTypeDelete,
			})
		}

	}
	// upstream diff
	targetUpstreamMap := utils.SliceToMap(target.upstreamCommits, func(i *commit.Commit[upstream.Config]) string {
		return fmt.Sprintf("%s-%s", i.Target, i.Key)
	})
	baseUpstreamMap := utils.SliceToMap(base.upstreamCommits, func(i *commit.Commit[upstream.Config]) string {
		return fmt.Sprintf("%s-%s", i.Target, i.Key)
	})

	for _, partitionId := range partitions {
		key := fmt.Sprintf("%s-%s", target.id, partitionId)
		o := &project_diff.UpstreamDiff{
			Upstream:  target.id,
			Partition: partitionId,
			Data:      nil,
			Change:    project_diff.ChangeTypeNone,
			Status:    0,
		}
		out.Upstreams = append(out.Upstreams, o)
		bu, hasBu := baseUpstreamMap[key]
		tu, hasTu := targetUpstreamMap[key]
		if hasTu {
			o.Data = tu.Data
			if !hasBu {
				o.Change = project_diff.ChangeTypeNew
			} else if tu.UUID != bu.UUID {
				o.Change = project_diff.ChangeTypeUpdate
			}

		} else {
			o.Status = project_diff.StatusLoss
			if hasBu {
				o.Change = project_diff.ChangeTypeDelete
			}
		}
	}

	return out
}

func (m *imlProjectDiff) Out(ctx context.Context, diff *project_diff.Diff) (*DiffOut, error) {

	clusters, err := m.clusterService.List(ctx, diff.Clusters...)
	if err != nil {
		return nil, err
	}
	if len(clusters) == 0 {
		return nil, fmt.Errorf("unset gateway for clusters %v", diff.Clusters)
	}
	//// 检查分区是否配置集群，若没有配置，则报错
	//requirePartition := make([]*partition.Partition, 0, len(partitions))
	//for _, p := range partitions {
	//	if p.Cluster == "" {
	//		requirePartition = append(requirePartition, p)
	//		continue
	//	}
	//	_, err = m.clusterService.Get(ctx, p.Cluster)
	//	if err != nil {
	//		if !errors.Is(err, gorm.ErrRecordNotFound) {
	//			requirePartition = append(requirePartition, p)
	//			continue
	//		}
	//		return nil, err
	//	}
	//
	//}
	//
	//if len(requirePartition) > 0 {
	//	return nil, fmt.Errorf("unset gateway for partitions %v", requirePartition)
	//}

	out := &DiffOut{}
	out.Apis = utils.SliceToSlice(diff.Apis, func(i *project_diff.ApiDiff) *ApiDiffOut {
		return &ApiDiffOut{
			Api:    auto.UUID(i.APi),
			Name:   i.Name,
			Method: i.Method,
			Path:   i.Path,
			Change: i.Change,
			Status: i.Status,
		}
	})

	for _, u := range diff.Upstreams {
		typeValue := u.Data.Type

		if typeValue == "" {
			typeValue = "static"
		}
		out.Upstreams = append(out.Upstreams, &UpstreamDiffOut{
			Partition: auto.UUID(u.Partition),
			Change:    u.Change,
			Type:      typeValue,
			Status:    u.Status,
			Addr: utils.SliceToSlice(u.Data.Nodes, func(i *upstream.NodeConfig) string {
				return i.Address
			}),
		})
	}
	return out, nil
}
