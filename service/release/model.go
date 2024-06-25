package release

import (
	"github.com/eolinker/apipark/stores/release"
	"time"
)

type Release struct {
	UUID     string
	Project  string
	Version  string
	Remark   string
	Creator  string
	CreateAt time.Time
}

func FromEntity(e *release.Release) *Release {
	return &Release{
		UUID:     e.UUID,
		Project:  e.Project,
		Version:  e.Name,
		Remark:   e.Remark,
		Creator:  e.Creator,
		CreateAt: e.CreateAt,
	}
}

type APIProxyCommit struct {
	Release string
	API     string
	Commit  string
}
type APIDocumentCommit struct {
	Release string
	API     string
	Commit  string
}
type UpstreamCommit struct {
	Release   string
	Upstream  string
	Partition string
	Commit    string
}

type ProjectCommits struct {
	Release string
	Type    string
	Target  string
	Key     string
	Commit  string
}

//type Diff struct {
//	Apis      []*APiDiff      `json:"apis"`
//	Upstreams []*UpstreamDiff `json:"upstream"`
//}

//type APiDiff struct {
//	Api string `json:"api,omitempty"`
//
//	Change project_diff.ChangeType `json:"change,omitempty"`
//}
//
//type UpstreamDiff struct {
//	UpstreamCommit  string                  `json:"upstream,omitempty"`
//	Partition string                  `json:"partition,omitempty"`
//	Commit    string                  `json:"commit,omitempty"`
//	Change    project_diff.ChangeType `json:"change,omitempty"`
//}
