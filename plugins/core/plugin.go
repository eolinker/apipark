package core

import (
	"github.com/eolinker/go-common/pm3"
	"net/http"
)

func (p *plugin) PartitionPluginApi() []pm3.Api {
	return []pm3.Api{
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/apinto/plugin/:plugin", []string{"context", "query:partition", "rest:plugin"}, []string{"plugin", "render"}, p.pluginPartitionController.Get),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/apinto/plugins", []string{"context", "query:partition"}, []string{"plugins"}, p.pluginPartitionController.List),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/simple/apinto/plugins/project", []string{"context", "query:project"}, []string{"plugins"}, p.pluginPartitionController.Option),
		pm3.CreateApiWidthDoc(http.MethodGet, "/api/v1/simple/apinto/plugin/:name", []string{"context", "rest:name"}, []string{"plugins"}, p.pluginPartitionController.Info),
		pm3.CreateApiWidthDoc(http.MethodPut, "/api/v1/apinto/plugin/:plugin", []string{"context", "query:partition", "rest:plugin", "body"}, []string{}, p.pluginPartitionController.Set),
	}
}
