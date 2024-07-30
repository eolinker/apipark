package core

import (
	"net/http"

	plugin_cluster "github.com/eolinker/apipark/controller/plugin-cluster"

	"github.com/eolinker/apipark/controller/cluster"

	"github.com/eolinker/ap-account/controller/role"

	"github.com/eolinker/apipark/controller/common"

	"github.com/eolinker/apipark/controller/topology"

	dynamic_module "github.com/eolinker/apipark/controller/dynamic-module"

	"github.com/eolinker/apipark/controller/release"

	project_authorization "github.com/eolinker/apipark/controller/project-authorization"

	"github.com/eolinker/apipark/controller/subscribe"

	"github.com/eolinker/apipark/controller/api"

	"github.com/eolinker/apipark/controller/upstream"

	"github.com/eolinker/apipark/controller/service"

	"github.com/eolinker/apipark/controller/catalogue"

	"github.com/eolinker/apipark/controller/project"

	"github.com/eolinker/apipark/controller/my_team"

	"github.com/eolinker/apipark/controller/certificate"
	"github.com/eolinker/apipark/controller/team_manager"
	"github.com/eolinker/go-common/autowire"
	"github.com/eolinker/go-common/pm3"
)

func init() {
	pm3.Register("core", new(Driver))
}

type Driver struct {
}

func (d *Driver) Access() map[string][]string {
	return map[string][]string{}
}

func (d *Driver) Create() (pm3.IPlugin, error) {

	p := new(plugin)
	autowire.Autowired(p)
	return p, nil
}

type plugin struct {
	clusterController              cluster.IClusterController                            `autowired:""`
	certificateController          certificate.ICertificateController                    `autowired:""`
	teamManagerController          team_manager.ITeamManagerController                   `autowired:""`
	myTeamController               my_team.ITeamController                               `autowired:""`
	appController                  project.IAppController                                `autowired:""`
	projectController              project.IProjectController                            `autowired:""`
	serviceController              service.IServiceController                            `autowired:""`
	catalogueController            catalogue.ICatalogueController                        `autowired:""`
	upstreamController             upstream.IUpstreamController                          `autowired:""`
	apiController                  api.IAPIController                                    `autowired:""`
	subscribeController            subscribe.ISubscribeController                        `autowired:""`
	projectAuthorizationController project_authorization.IProjectAuthorizationController `autowired:""`
	releaseController              release.IReleaseController                            `autowired:""`
	roleController                 role.IRoleController                                  `autowired:""`
	subscribeApprovalController    subscribe.ISubscribeApprovalController                `autowired:""`
	dynamicModuleController        dynamic_module.IDynamicModuleController               `autowired:""`
	topologyController             topology.ITopologyController                          `autowired:""`
	pluginClusterController        plugin_cluster.IPluginClusterController               `autowired:""`
	commonController               common.ICommonController                              `autowired:""`
	apis                           []pm3.Api
}

func (p *plugin) OnComplete() {
	p.apis = append(p.apis, p.partitionApi()...)
	p.apis = append(p.apis, p.certificateApi()...)
	p.apis = append(p.apis, p.clusterApi()...)
	p.apis = append(p.apis, p.TeamManagerApi()...)
	p.apis = append(p.apis, p.MyTeamApi()...)
	p.apis = append(p.apis, p.ProjectApi()...)
	p.apis = append(p.apis, p.catalogueApi()...)
	p.apis = append(p.apis, p.serviceApi()...)
	p.apis = append(p.apis, p.upstreamApis()...)
	p.apis = append(p.apis, p.apiApis()...)
	p.apis = append(p.apis, p.subscribeApis()...)
	p.apis = append(p.apis, p.projectAuthorizationApis()...)
	p.apis = append(p.apis, p.releaseApis()...)
	p.apis = append(p.apis, p.DynamicModuleApis()...)
	p.apis = append(p.apis, p.TopologyApis()...)
	p.apis = append(p.apis, p.PartitionPluginApi()...)
	p.apis = append(p.apis, p.commonApis()...)
}

func (p *plugin) Name() string {
	return "core"
}

func (p *plugin) APis() []pm3.Api {
	return p.apis
}

func (p *plugin) Assets() map[string]http.FileSystem {
	return nil
}
