package publish_flow

import (
	"net/http"

	"github.com/eolinker/go-common/pm3"
)

func (p *plugin) getApis() []pm3.Api {
	return []pm3.Api{
		pm3.CreateApiWidthDoc(http.MethodPost, "/api/v1/project/publish/release", []string{"context", "query:project", "body"}, []string{"publish"}, p.controller.ApplyOnRelease),
		pm3.CreateApiWidthDoc(http.MethodPost, "/api/v1/project/publish/apply", []string{"context", "query:project", "body"}, []string{"publish"}, p.controller.Apply),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/publishs", []string{"context", "query:project", "query:page", "query:page_size"}, []string{"publishs", "page", "size", "total"}, p.controller.ListPage),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/publish", []string{"context", "query:project", "query:id"}, []string{"publish"}, p.controller.Detail),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/publish/check", []string{"context", "query:project", "query:release"}, []string{"diffs"}, p.controller.CheckPublish),
		pm3.CreateApiWidthDoc(http.MethodDelete, "/api/v1/project/publish/close", []string{"context", "query:project", "query:id"}, []string{}, p.controller.Close),
		pm3.CreateApiWidthDoc(http.MethodDelete, "/api/v1/project/publish/stop", []string{"context", "query:project", "query:id"}, []string{}, p.controller.Stop),
		pm3.CreateApiWidthDoc(http.MethodPut, "/api/v1/project/publish/refuse", []string{"context", "query:project", "query:id", "body"}, []string{}, p.controller.Refuse),
		pm3.CreateApiWidthDoc(http.MethodPut, "/api/v1/project/publish/accept", []string{"context", "query:project", "query:id", "body"}, []string{}, p.controller.Accept),
		pm3.CreateApiWidthDoc(http.MethodPut, "/api/v1/project/publish/execute", []string{"context", "query:project", "query:id"}, []string{}, p.controller.Publish),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/project/publish/status", []string{"context", "query:project", "query:id"}, []string{"publish_status_list"}, p.controller.PublishStatuses),
	}
}
