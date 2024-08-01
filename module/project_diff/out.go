package project_diff

import (
	"github.com/eolinker/apipark/service/api"
	"github.com/eolinker/apipark/service/project_diff"
	"github.com/eolinker/apipark/service/universally/commit"
	"github.com/eolinker/apipark/service/upstream"
	"github.com/eolinker/go-common/auto"
)

type DiffOut struct {
	Apis      []*ApiDiffOut      `json:"apis"`
	Upstreams []*UpstreamDiffOut `json:"upstreams"`
}

type ApiDiffOut struct {
	Api    auto.Label `json:"api,omitempty" aolabel:"api"`
	Name   string     `json:"name,omitempty"`
	Method string     `json:"method,omitempty"`
	Path   string     `json:"path,omitempty"`
	//Upstream auto.Label              `json:"upstream,omitempty" aolabel:"upstream"`
	Change project_diff.ChangeType `json:"change,omitempty"`
	Status project_diff.Status     `json:"status,omitempty"`
}
type UpstreamDiffOut struct {
	//Upstream  auto.Label              `json:"upstream,omitempty" aolabel:"upstream"`
	Partition auto.Label `json:"partition,omitempty" aolabel:"partition"`
	//Cluster   auto.Label              `json:"cluster,omitempty" aolabel:"cluster"`
	Change project_diff.ChangeType `json:"change,omitempty"`
	Status project_diff.StatusType `json:"status,omitempty"`
	Type   string                  `json:"type,omitempty"`
	Addr   []string                `json:"addr,omitempty"`
}

//
//func CreateOut(d *project_diff.Diff) *DiffOut {
//	if d == nil {
//		return nil
//	}
//	return &DiffOut{
//		Apis: utils.SliceToSlice(d.Apis, func(s *project_diff.ApiDiff) *ApiDiffOut {
//			return &ApiDiffOut{
//				Name:     s.Name,
//				Method:   s.Method,
//				Path:     s.Path,
//				Upstream: s.Upstream,
//				Change:   s.Change,
//			}
//		}),
//		Upstreams: utils.SliceToSlice(d.Upstreams, func(s *project_diff.UpstreamDiff) *UpstreamDiffOut {
//			return &UpstreamDiffOut{
//				Upstream:  s.Name,
//				Cluster: auto.UUID(s.Cluster),
//				Cluster:   auto.UUID(s.Cluster),
//				Change:    s.Change,
//				Type:      s.Type,
//				Addr:      s.Addr,
//			}
//		}),
//	}
//}

type projectInfo struct {
	id              string
	apis            []*api.Info
	apiCommits      []*commit.Commit[api.Proxy]
	apiDocs         []*commit.Commit[api.Document]
	upstreamCommits []*commit.Commit[upstream.Config]
}
