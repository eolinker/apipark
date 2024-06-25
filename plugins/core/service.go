package core

import (
	"net/http"

	"github.com/eolinker/go-common/pm3"
)

func (p *plugin) serviceApi() []pm3.Api {
	return []pm3.Api{
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/services", []string{"context", "query:keyword", "query:project"}, []string{"services"}, p.serviceController.Search),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/service/info", []string{"context", "query:project", "query:service"}, []string{"service"}, p.serviceController.Get),
		pm3.CreateApiWidthDoc(http.MethodPost, "/api/v1/project/service", []string{"context", "query:project", "body"}, []string{"service"}, p.serviceController.Create),
		pm3.CreateApiWidthDoc(http.MethodPut, "/api/v1/project/service/info", []string{"context", "query:project", "query:service", "body"}, []string{"service"}, p.serviceController.Edit),
		pm3.CreateApiWidthDoc(http.MethodDelete, "/api/v1/project/service", []string{"context", "query:project", "query:service"}, nil, p.serviceController.Delete),
		pm3.CreateApiWidthDoc(http.MethodPut, "/api/v1/project/service/enable", []string{"context", "query:project", "query:service"}, nil, p.serviceController.Enable),
		pm3.CreateApiWidthDoc(http.MethodPut, "/api/v1/project/service/disable", []string{"context", "query:project", "query:service"}, nil, p.serviceController.Disable),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/service/doc", []string{"context", "query:project", "query:service"}, []string{"doc"}, p.serviceController.ServiceDoc),
		pm3.CreateApiWidthDoc(http.MethodPut, "/api/v1/project/service/doc", []string{"context", "query:project", "query:service", "body"}, nil, p.serviceController.SaveServiceDoc),

		// Service API
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/service/apis", []string{"context", "query:project", "query:service"}, []string{"apis"}, p.serviceController.ServiceApis),
		pm3.CreateApiWidthDoc(http.MethodPost, "/api/v1/project/service/bind", []string{"context", "query:project", "query:service", "body"}, nil, p.serviceController.BindServiceApi),
		pm3.CreateApiWidthDoc(http.MethodDelete, "/api/v1/project/service/unbind", []string{"context", "query:project", "query:service", "query:apis"}, nil, p.serviceController.UnbindServiceApi),
		pm3.CreateApiWidthDoc(http.MethodPut, "/api/v1/project/service/api/sort", []string{"context", "query:project", "query:service", "body"}, nil, p.serviceController.SortApis),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/simple/project/services", []string{"context", "query:project"}, []string{"services"}, p.serviceController.SimpleList),
	}
}
