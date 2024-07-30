package core

import (
	"net/http"

	"github.com/eolinker/go-common/pm3"
)

func (p *plugin) serviceApi() []pm3.Api {
	return []pm3.Api{

		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/service/doc", []string{"context", "query:project"}, []string{"doc"}, p.serviceController.ServiceDoc),
		pm3.CreateApiWidthDoc(http.MethodPut, "/api/v1/project/service/doc", []string{"context", "query:project", "body"}, nil, p.serviceController.SaveServiceDoc),
	}
}
