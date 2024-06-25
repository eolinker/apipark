package core

import (
	"net/http"

	"github.com/eolinker/go-common/pm3"
)

func (p *plugin) clusterApi() []pm3.Api {
	return []pm3.Api{
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/partition/cluster/nodes", []string{"context", "query:partition"}, []string{"nodes"}, p.partitionController.Nodes),
		pm3.CreateApiWidthDoc(http.MethodPut, "/api/v1/partition/cluster/reset", []string{"context", "query:partition", "body"}, []string{"nodes"}, p.partitionController.ResetCluster),
		pm3.CreateApiWidthDoc(http.MethodPost, "/api/v1/partition/cluster/check", []string{"context", "body"}, []string{"nodes"}, p.partitionController.Check),
	}
}
