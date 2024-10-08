package release

import (
	"context"
	"errors"
	"fmt"

	"github.com/eolinker/apipark/service/cluster"
	"github.com/eolinker/apipark/service/service"
	"github.com/eolinker/apipark/service/service_diff"

	"github.com/eolinker/apipark/module/release/dto"
	serviceDiff "github.com/eolinker/apipark/module/service-diff"
	"github.com/eolinker/apipark/service/api"
	"github.com/eolinker/apipark/service/publish"
	"github.com/eolinker/apipark/service/release"
	"github.com/eolinker/apipark/service/universally/commit"
	"github.com/eolinker/apipark/service/upstream"
	"github.com/eolinker/go-common/auto"
	"github.com/eolinker/go-common/store"
	"github.com/eolinker/go-common/utils"
	"gorm.io/gorm"
)

var (
	_                     IReleaseModule = (*imlReleaseModule)(nil)
	projectRuleMustServer                = map[string]bool{
		"as_server": true,
	}
)

type imlReleaseModule struct {
	projectDiffModule serviceDiff.IServiceDiffModule `autowired:""`
	releaseService    release.IReleaseService        `autowired:""`
	apiService        api.IAPIService                `autowired:""`
	upstreamService   upstream.IUpstreamService      `autowired:""`
	publishService    publish.IPublishService        `autowired:""`
	transaction       store.ITransaction             `autowired:""`
	projectService    service.IServiceService        `autowired:""`
	clusterService    cluster.IClusterService        `autowired:""`
}

func (m *imlReleaseModule) Create(ctx context.Context, serviceId string, input *dto.CreateInput) (string, error) {

	proInfo, err := m.projectService.Check(ctx, serviceId, projectRuleMustServer)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("project not found")
		}
		return "", err
	}
	clusters, err := m.clusterService.List(ctx)
	if err != nil || len(clusters) == 0 {
		return "", fmt.Errorf("cluster not set:%w", err)
	}

	apis, err := m.apiService.ListForService(ctx, proInfo.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("api not found")
		}
		return "", err
	}
	if len(apis) == 0 {
		return "", errors.New("api not found")
	}
	apiUUIDS := utils.SliceToSlice(apis, func(a *api.API) string {
		return a.UUID
	})
	apiProxy, err := m.apiService.ListLatestCommitProxy(ctx, apiUUIDS...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("api  config or  document not found")
		}
		return "", err
	}
	if len(apis) != len(apiProxy) {
		return "", errors.New("api or document not found")
	}
	apiDocs, err := m.apiService.ListLatestCommitDocument(ctx, apiUUIDS...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("api  config or  document not found")
		}
		return "", err
	}
	if len(apis) != len(apiDocs) {
		return "", errors.New("api or document not found")
	}
	upstreams, err := m.upstreamService.ListLatestCommit(ctx, serviceId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("api  config or  document not found")
		}
		return "", err
	}

	apiProxyCommits := utils.SliceToMapO(apiProxy, func(c *commit.Commit[api.Proxy]) (string, string) {
		return c.Target, c.UUID
	})
	apiDocumentCommits := utils.SliceToMapO(apiDocs, func(c *commit.Commit[api.Document]) (string, string) {
		return c.Target, c.UUID
	})
	upstreamCommits := utils.SliceToMapArray(upstreams, func(c *commit.Commit[upstream.Config]) string {
		return c.Target
	})
	upstreamCommitsForUKC := utils.MapChange(upstreamCommits, func(ls []*commit.Commit[upstream.Config]) map[string]string {
		return utils.SliceToMapO(ls, func(c *commit.Commit[upstream.Config]) (string, string) {
			return c.Key, c.UUID
		})
	})
	if !m.releaseService.Completeness(utils.SliceToSlice(clusters, func(s *cluster.Cluster) string {
		return s.Uuid
	}), apiUUIDS, apiProxy, apiDocs, upstreams) {
		return "", errors.New("completeness check failed")
	}
	newRelease, err := m.releaseService.CreateRelease(ctx, serviceId, input.Version, input.Remark, apiProxyCommits, apiDocumentCommits, upstreamCommitsForUKC)
	if err != nil {
		return "", err
	}
	return newRelease.UUID, err
}

func (m *imlReleaseModule) Detail(ctx context.Context, project string, id string) (*dto.Detail, error) {
	r, err := m.releaseService.GetRelease(ctx, id)
	if err != nil {
		return nil, err
	}
	if r.Service != project {
		return nil, errors.New("release not found")
	}
	running, err := m.releaseService.GetRunning(ctx, project)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	runningRelease := ""
	if running != nil {
		runningRelease = running.UUID
	}
	diff, err := m.projectDiffModule.Diff(ctx, project, runningRelease, r.UUID)
	if err != nil {
		return nil, err
	}
	out, err := m.projectDiffModule.Out(ctx, diff)
	if err != nil {
		return nil, err
	}
	return &dto.Detail{
		Id:         r.UUID,
		Version:    r.Version,
		Remark:     r.Remark,
		Service:    auto.UUID(r.Service),
		CreateTime: auto.TimeLabel(r.CreateAt),
		Creator:    auto.UUID(r.Creator),
		Diffs:      out,
	}, nil
}

func (m *imlReleaseModule) List(ctx context.Context, project string) ([]*dto.Release, error) {
	_, err := m.projectService.Check(ctx, project, projectRuleMustServer)
	if err != nil {
		return nil, err
	}
	list, err := m.releaseService.List(ctx, project)
	if err != nil {
		return nil, err
	}
	running, err := m.releaseService.GetRunning(ctx, project)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	releaseIds := utils.SliceToSlice(list, func(s *release.Release) string {
		return s.UUID
	})
	flows, err := m.publishService.Latest(ctx, releaseIds...)
	if err != nil {
		return nil, err
	}
	flowMap := utils.SliceToMap(flows, func(s *publish.Publish) string {
		return s.Release
	})

	return utils.SliceToSlice(list, func(s *release.Release) *dto.Release {

		r := &dto.Release{
			Id:          s.UUID,
			Service:     auto.UUID(s.Service),
			Version:     s.Version,
			Remark:      s.Remark,
			Status:      dto.StatusNone,
			CanRollback: false,
			CanDelete:   true,
			Creator:     auto.UUID(s.Creator),
			CreateTime:  auto.TimeLabel(s.CreateAt),
		}

		if running != nil && running.UUID == s.UUID {
			r.Status = dto.StatusRunning
			r.CanRollback = true
			r.CanDelete = false
		}
		flow, has := flowMap[s.UUID]
		if has {
			r.FlowId = flow.Id

			if flow.Status == publish.StatusApply {
				r.Status = dto.StatusApply
				r.CanDelete = false
			} else if flow.Status == publish.StatusAccept {

				r.Status = dto.StatusAccept
				r.CanDelete = false
			} else if flow.Status == publish.StatusPublishError {
				r.Status = dto.StatusError
				r.CanDelete = false
			}
		}
		return r
	}), nil
}

func (m *imlReleaseModule) Delete(ctx context.Context, project string, id string) error {
	_, err := m.projectService.Check(ctx, project, projectRuleMustServer)
	if err != nil {
		return err
	}
	return m.transaction.Transaction(ctx, func(ctx context.Context) error {
		r, err := m.releaseService.GetRelease(ctx, id)
		if err != nil {
			return err
		}
		if r == nil {
			return errors.New("release not found")
		}
		if r.Service != project {
			return errors.New("project not match")
		}
		running, err := m.releaseService.GetRunning(ctx, project)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if running != nil && running.UUID == id {
			return errors.New("can not delete running release")
		}
		flow, err := m.publishService.GetLatest(ctx, id)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if flow != nil {
			if flow.Status == publish.StatusApply || flow.Status == publish.StatusAccept {
				return errors.New("can not delete  release in apply or approve flow")
			}
		}
		return m.releaseService.DeleteRelease(ctx, id)
	})

}

func (m *imlReleaseModule) Preview(ctx context.Context, project string) (*dto.Release, *service_diff.Diff, bool, error) {
	_, err := m.projectService.Check(ctx, project, projectRuleMustServer)
	if err != nil {
		return nil, nil, false, err
	}
	running, err := m.releaseService.GetRunning(ctx, project)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, false, err
	}

	if running == nil {
		running = new(release.Release)
	}

	diff, completeness, err := m.projectDiffModule.DiffForLatest(ctx, project, running.UUID)
	if err != nil {
		return nil, nil, false, err
	}
	return &dto.Release{
		Id:          running.UUID,
		Version:     running.Version,
		Service:     auto.UUID(project),
		CreateTime:  auto.TimeLabel(running.CreateAt),
		Creator:     auto.UUID(running.Creator),
		Status:      dto.StatusNone,
		Remark:      running.Remark,
		CanDelete:   false,
		CanRollback: false,
	}, diff, completeness, nil

}
