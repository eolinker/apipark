package release

import (
	"context"
	"github.com/eolinker/apipark/service/api"
	"github.com/eolinker/apipark/service/universally/commit"
	"github.com/eolinker/apipark/service/upstream"
	"github.com/eolinker/go-common/autowire"
	"reflect"
)

type IReleaseService interface {
	// GetRelease 获取发布信息
	GetRelease(ctx context.Context, id string) (*Release, error)
	// CreateRelease 创建发布
	CreateRelease(ctx context.Context, project string, version string, remark string, apisProxyCommits, apiDocCommits map[string]string, upstreams map[string]map[string]string) (*Release, error)
	//Diff(ctx context.Context, baseReleaseId string, targetReleaseId string) (*Diff, error)
	//DiffApis(ctx context.Context, baseAPi []*Api, targetAPiProxy []*Api) []*APiDiff
	//DiffUpstreams(ctx context.Context, baseUpstream []*UpstreamCommit, targetUpstream []*UpstreamCommit) []*UpstreamDiff
	// DeleteRelease 删除发布
	DeleteRelease(ctx context.Context, id string) error
	List(ctx context.Context, project string) ([]*Release, error)
	GetApiProxyCommit(ctx context.Context, id string, apiUUID string) (string, error)
	GetApiDocCommit(ctx context.Context, id string, apiUUID string) (string, error)
	GetReleaseInfos(ctx context.Context, id string) ([]*APIProxyCommit, []*APIDocumentCommit, []*UpstreamCommit, error)
	//GetApiProxyCommit(ctx context.Context, id string, apiUUID string) (string, error)
	//GetApiDocCommit(ctx context.Context, id string, apiUUID string) (string, error)
	GetCommits(ctx context.Context, id string) ([]*ProjectCommits, error)

	GetRunningApiDocCommit(ctx context.Context, project string, apiUUID string) (string, error)
	GetRunningApiProxyCommit(ctx context.Context, project string, apiUUID string) (string, error)
	Completeness(partitions []string, apis []string, proxyCommits []*commit.Commit[api.Proxy], documentCommits []*commit.Commit[api.Document], upstreamCommits []*commit.Commit[upstream.Config]) bool

	// GetRunning gets the running release with the given project.
	//
	// ctx: the context
	// project: the project name
	// Return type(s): *Release, error
	GetRunning(ctx context.Context, project string) (*Release, error)

	SetRunning(ctx context.Context, project string, id string) error
	CheckNewVersion(ctx context.Context, project string, version string) (bool, error)
}

func init() {
	autowire.Auto[IReleaseService](func() reflect.Value {
		return reflect.ValueOf(new(imlReleaseService))
	})
}
