package project_authorization

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/eolinker/apipark/service/partition"

	"github.com/eolinker/eosc/log"

	authDriver "github.com/eolinker/apipark/module/project-authorization/auth-driver"

	"github.com/eolinker/go-common/utils"

	"github.com/eolinker/apipark/gateway"

	"github.com/eolinker/apipark/service/cluster"

	"github.com/eolinker/go-common/auto"

	"github.com/google/uuid"

	"github.com/eolinker/go-common/store"

	"github.com/eolinker/apipark/service/project"
	projectAuthorization "github.com/eolinker/apipark/service/project-authorization"

	projectAuthorizationDto "github.com/eolinker/apipark/module/project-authorization/dto"
)

var _ IProjectAuthorizationModule = (*imlProjectAuthorizationModule)(nil)

type imlProjectAuthorizationModule struct {
	projectService              project.IProjectService                           `autowired:""`
	projectAuthorizationService projectAuthorization.IProjectAuthorizationService `autowired:""`
	partitionService            partition.IPartitionService                       `autowired:""`
	clusterService              cluster.IClusterService                           `autowired:""`
	transaction                 store.ITransaction                                `autowired:""`
}

func (i *imlProjectAuthorizationModule) getApplications(ctx context.Context, projectIds []string, projectMap map[string]*project.Project) ([]*gateway.ApplicationRelease, error) {
	authorizations, err := i.projectAuthorizationService.ListByProject(ctx, projectIds...)
	if err != nil {
		return nil, err
	}
	authMap := utils.SliceToMapArray(authorizations, func(a *projectAuthorization.Authorization) string {
		return a.Project
	})
	return utils.SliceToSlice(projectIds, func(projectId string) *gateway.ApplicationRelease {
		auths := authMap[projectId]
		description := ""
		projectInfo, ok := projectMap[projectId]
		if ok {
			description = projectInfo.Description
		}
		return &gateway.ApplicationRelease{
			BasicItem: &gateway.BasicItem{
				ID:          projectId,
				Description: description,
				Version:     time.Now().Format("20060102150405"),
				MatchLabels: map[string]string{
					"project": projectId,
				},
			},

			Authorizations: utils.SliceToSlice(auths, func(a *projectAuthorization.Authorization) *gateway.Authorization {
				authCfg := make(map[string]interface{})
				_ = json.Unmarshal([]byte(a.Config), &authCfg)
				return &gateway.Authorization{
					Type:           a.Type,
					Position:       a.Position,
					TokenName:      a.TokenName,
					Expire:         a.ExpireTime,
					Config:         authCfg,
					HideCredential: a.HideCredential,
				}
			}),
		}
	}), nil
}

func (i *imlProjectAuthorizationModule) initGateway(ctx context.Context, partitionId string, clientDriver gateway.IClientDriver) error {
	projects, err := i.projectService.List(ctx)
	if err != nil {
		return err
	}
	projectIds := make([]string, 0, len(projects))
	projectMap := make(map[string]*project.Project)
	for _, p := range projects {
		projectIds = append(projectIds, p.Id)
		projectMap[p.Id] = p
	}

	applications, err := i.getApplications(ctx, projectIds, projectMap)
	if err != nil {
		return err
	}
	return clientDriver.Application().Online(ctx, applications...)
}

func (i *imlProjectAuthorizationModule) online(ctx context.Context, projectInfo *project.Project) error {

	partitions, err := i.partitionService.List(ctx)
	if err != nil {
		return err
	}
	partitionIds := utils.SliceToSlice(partitions, func(p *partition.Partition) string {
		return p.UUID
	})
	clusters, err := i.clusterService.List(ctx, partitionIds...)
	if err != nil {
		return err
	}
	authorizations, err := i.projectAuthorizationService.ListByProject(ctx, projectInfo.Id)
	if err != nil {
		return err
	}
	app := &gateway.ApplicationRelease{
		BasicItem: &gateway.BasicItem{
			ID:          projectInfo.Id,
			Description: projectInfo.Description,
			Version:     time.Now().Format("20060102150405"),
			MatchLabels: map[string]string{
				"project": projectInfo.Id,
			},
		},
		Authorizations: utils.SliceToSlice(authorizations, func(a *projectAuthorization.Authorization) *gateway.Authorization {
			authCfg := make(map[string]interface{})
			_ = json.Unmarshal([]byte(a.Config), &authCfg)
			return &gateway.Authorization{
				Type:           a.Type,
				Position:       a.Position,
				TokenName:      a.TokenName,
				Expire:         a.ExpireTime,
				Config:         authCfg,
				HideCredential: a.HideCredential,
			}
		}),
	}

	for _, c := range clusters {
		err := i.doOnline(ctx, c.Uuid, app)
		if err != nil {
			log.Warnf("project authorization online for partition[%s] %v", c.Partition, err)
		}
	}
	return nil
}
func (i *imlProjectAuthorizationModule) doOnline(ctx context.Context, clusterId string, app *gateway.ApplicationRelease) error {
	client, err := i.clusterService.GatewayClient(ctx, clusterId)
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Close(ctx)
	}()
	return client.Application().Online(ctx, app)

}
func (i *imlProjectAuthorizationModule) AddAuthorization(ctx context.Context, pid string, info *projectAuthorizationDto.CreateAuthorization) (*projectAuthorizationDto.Authorization, error) {
	authFactory, has := authDriver.GetAuthFactory(info.Driver)
	if !has {
		return nil, fmt.Errorf("unknown driver %s", info.Driver)
	}
	auth, err := authFactory.Create(info.Config)
	if err != nil {
		return nil, err
	}
	cfg, err := auth.AuthConfig().Valid()
	if err != nil {
		return nil, err
	}

	projectInfo, err := i.projectService.Get(ctx, pid)
	if err != nil {
		return nil, err
	}

	if info.UUID == "" {
		info.UUID = uuid.New().String()
	}

	// 缺少配置查重操作
	//cfg, _ := json.Marshal(info.Config)
	err = i.transaction.Transaction(ctx, func(ctx context.Context) error {
		err = i.projectAuthorizationService.Create(ctx, &projectAuthorization.CreateAuthorization{
			UUID:           info.UUID,
			Project:        pid,
			Name:           info.Name,
			Type:           info.Driver,
			Position:       info.Position,
			TokenName:      info.TokenName,
			Config:         string(cfg),
			ExpireTime:     info.ExpireTime,
			HideCredential: info.HideCredential,
			AuthID:         auth.GenerateID(info.Position, info.TokenName),
		})
		if err != nil {
			return err
		}

		return i.online(ctx, projectInfo)
	})
	if err != nil {
		return nil, err
	}

	return i.Info(ctx, pid, info.UUID)
}

func (i *imlProjectAuthorizationModule) EditAuthorization(ctx context.Context, pid string, aid string, info *projectAuthorizationDto.EditAuthorization) (*projectAuthorizationDto.Authorization, error) {
	authInfo, err := i.projectAuthorizationService.Get(ctx, aid)
	if err != nil {
		return nil, err
	}
	authFactory, has := authDriver.GetAuthFactory(authInfo.Type)
	if !has {
		return nil, fmt.Errorf("unknown driver %s", authInfo.Type)
	}
	auth, err := authFactory.Create(info.Config)
	if err != nil {
		return nil, err
	}
	cfg, err := auth.AuthConfig().Valid()
	if err != nil {
		return nil, err
	}

	projectInfo, err := i.projectService.Get(ctx, pid)
	if err != nil {
		return nil, err
	}

	err = i.transaction.Transaction(ctx, func(ctx context.Context) error {
		authId := auth.GenerateID(authInfo.Position, authInfo.TokenName)
		cfgStr := string(cfg)
		err = i.projectAuthorizationService.Save(ctx, aid, &projectAuthorization.EditAuthorization{
			Name:           info.Name,
			Position:       info.Position,
			TokenName:      info.TokenName,
			ExpireTime:     info.ExpireTime,
			HideCredential: info.HideCredential,
			AuthID:         &authId,
			Config:         &cfgStr,
		})
		if err != nil {
			return err
		}
		return i.online(ctx, projectInfo)
	})

	if err != nil {
		return nil, err
	}
	return i.Info(ctx, pid, aid)
}

func (i *imlProjectAuthorizationModule) DeleteAuthorization(ctx context.Context, pid string, aid string) error {
	_, err := i.projectService.Get(ctx, pid)
	if err != nil {
		return err
	}
	partitions, err := i.partitionService.List(ctx)
	if err != nil {
		return err
	}
	partitionIds := utils.SliceToSlice(partitions, func(p *partition.Partition) string {
		return p.UUID
	})
	return i.transaction.Transaction(ctx, func(ctx context.Context) error {
		err = i.projectAuthorizationService.Delete(ctx, aid)
		if err != nil {
			return err
		}
		clusters, err := i.clusterService.List(ctx, partitionIds...)
		if err != nil {
			return err
		}
		app := &gateway.ApplicationRelease{
			BasicItem: &gateway.BasicItem{
				ID: pid,
			},
		}
		for _, c := range clusters {
			err := i.doOffline(ctx, c.Uuid, app)
			if err != nil {
				log.Warnf("project authorization offline for partition[%s] %v", c.Partition, err)
			}
		}
		return nil
	})
}
func (i *imlProjectAuthorizationModule) doOffline(ctx context.Context, clusterId string, app *gateway.ApplicationRelease) error {
	client, err := i.clusterService.GatewayClient(ctx, clusterId)
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Close(ctx)
	}()
	return client.Application().Offline(ctx, app)

}
func (i *imlProjectAuthorizationModule) Authorizations(ctx context.Context, pid string) ([]*projectAuthorizationDto.AuthorizationItem, error) {
	_, err := i.projectService.Get(ctx, pid)
	if err != nil {
		return nil, err
	}
	authorizations, err := i.projectAuthorizationService.ListByProject(ctx, pid)
	if err != nil {
		return nil, err
	}
	result := make([]*projectAuthorizationDto.AuthorizationItem, 0, len(authorizations))
	for _, a := range authorizations {
		result = append(result, &projectAuthorizationDto.AuthorizationItem{
			Id:             a.UUID,
			Name:           a.Name,
			Driver:         a.Type,
			ExpireTime:     a.ExpireTime,
			Position:       a.Position,
			TokenName:      a.TokenName,
			Creator:        auto.UUID(a.Creator),
			Updater:        auto.UUID(a.Updater),
			CreateTime:     auto.TimeLabel(a.CreateTime),
			UpdateTime:     auto.TimeLabel(a.UpdateTime),
			HideCredential: a.HideCredential,
		})
	}
	return result, nil
}

func (i *imlProjectAuthorizationModule) Detail(ctx context.Context, pid string, aid string) ([]projectAuthorizationDto.DetailItem, error) {
	_, err := i.projectService.Get(ctx, pid)
	if err != nil {
		return nil, err
	}
	authInfo, err := i.projectAuthorizationService.Get(ctx, aid)
	if err != nil {
		return nil, err
	}
	authFactory, has := authDriver.GetAuthFactory(authInfo.Type)
	if !has {
		return nil, errors.New("unknown driver")
	}
	auth, err := authFactory.Create(authInfo.Config)
	if err != nil {
		return nil, err
	}
	cfgItems := auth.AuthConfig().Detail()
	details := make([]projectAuthorizationDto.DetailItem, 0, 6+len(cfgItems))
	details = append(details, projectAuthorizationDto.DetailItem{Key: "名称", Value: authInfo.Name})
	details = append(details, projectAuthorizationDto.DetailItem{Key: "鉴权类型", Value: authInfo.Type})
	details = append(details, projectAuthorizationDto.DetailItem{Key: "参数位置", Value: authInfo.Position})
	details = append(details, projectAuthorizationDto.DetailItem{Key: "参数名", Value: authInfo.TokenName})
	details = append(details, cfgItems...)
	dateStr := "永久"
	if authInfo.ExpireTime != 0 {
		dateStr = time.Unix(authInfo.ExpireTime, 0).Format("2006-01-02")
	}
	details = append(details, projectAuthorizationDto.DetailItem{Key: "过期日期", Value: dateStr})
	hideAuthStr := "是"
	if !authInfo.HideCredential {
		hideAuthStr = "否"
	}
	details = append(details, projectAuthorizationDto.DetailItem{Key: "隐藏鉴权信息", Value: hideAuthStr})

	return details, nil
}

func (i *imlProjectAuthorizationModule) Info(ctx context.Context, pid string, aid string) (*projectAuthorizationDto.Authorization, error) {
	_, err := i.projectService.Get(ctx, pid)
	if err != nil {
		return nil, err
	}
	auth, err := i.projectAuthorizationService.Get(ctx, aid)
	if err != nil {
		return nil, err
	}
	var cfg map[string]interface{}
	if auth.Config != "" {
		_ = json.Unmarshal([]byte(auth.Config), &cfg)
	}

	return &projectAuthorizationDto.Authorization{
		UUID:           auth.UUID,
		Name:           auth.Name,
		Driver:         auth.Type,
		Position:       auth.Position,
		TokenName:      auth.TokenName,
		ExpireTime:     auth.ExpireTime,
		HideCredential: auth.HideCredential,
		Config:         cfg,
	}, nil
}
