package dto

import (
	"github.com/eolinker/apipark/module/publish/dto"
	"github.com/eolinker/go-common/auto"
)

type Release struct {
	Id          string         `json:"id,omitempty"`
	Version     string         `json:"version,omitempty"`
	Project     auto.Label     `json:"project,omitempty" aolabel:"project"`
	CreateTime  auto.TimeLabel `json:"create_time"`
	Creator     auto.Label     `json:"creator" aolabel:"user"`
	Status      Status         `json:"status,omitempty"`
	FlowId      string         `json:"flowId,omitempty"`
	Remark      string         `json:"remark,omitempty"`
	CanDelete   bool           `json:"can_delete,omitempty"`
	CanRollback bool           `json:"can_rollback,omitempty"`
}

type Detail struct {
	Id         string         `json:"id,omitempty"`
	Version    string         `json:"version,omitempty"`
	Remark     string         `json:"remark,omitempty"`
	Project    auto.Label     `json:"project,omitempty" aolabel:"project"`
	CreateTime auto.TimeLabel `json:"createTime"`
	Creator    auto.Label     `json:"creator" aolabel:"user"`
	//Apis       []*project_diff.ApiDiff      `json:"apis,omitempty"`
	//Upstreams  []*project_diff.UpstreamDiff `json:"upstreams,omitempty"`
	Diffs *dto.DiffOut `json:"diffs,omitempty"`
}
