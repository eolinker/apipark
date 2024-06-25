package core

import (
	"net/http"

	"github.com/eolinker/go-common/pm3"
)

func (p *plugin) TopologyApis() []pm3.Api {
	return []pm3.Api{
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/topology", []string{"context"}, []string{"projects", "services"}, p.topologyController.SystemTopology),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/topology", []string{"context", "query:project"}, []string{"services", "subscribers", "invoke"}, p.topologyController.ProjectTopology),
	}
}
