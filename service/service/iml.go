package service

import (
	"context"
	"fmt"
	"time"

	"github.com/eolinker/go-common/utils"

	"github.com/eolinker/go-common/auto"

	"github.com/eolinker/apipark/service/universally"
	"github.com/eolinker/apipark/stores/service"
)

var (
	_ IServiceService = (*imlServiceService)(nil)
)

type imlServiceService struct {
	store service.IServiceStore `autowired:""`
	universally.IServiceGet[Service]
	universally.IServiceDelete
	universally.IServiceCreate[Create]
	universally.IServiceEdit[Edit]
}

func (i *imlServiceService) ServiceCountByTeam(ctx context.Context, teamId ...string) (map[string]int64, error) {
	w := map[string]interface{}{
		"as_server": true,
	}
	if len(teamId) > 0 {
		w["team"] = teamId
	}
	return i.store.CountByGroup(ctx, "", w, "team")
}

func (i *imlServiceService) AppCountByTeam(ctx context.Context, teamId ...string) (map[string]int64, error) {
	w := map[string]interface{}{
		"as_app": true,
	}
	if len(teamId) > 0 {
		w["team"] = teamId
	}
	return i.store.CountByGroup(ctx, "", w, "team")
}

func (i *imlServiceService) AppList(ctx context.Context, appIds ...string) ([]*Service, error) {
	w := make(map[string]interface{})
	if len(appIds) > 0 {
		w["uuid"] = appIds
	}
	w["as_app"] = true
	list, err := i.store.List(ctx, w)
	if err != nil {
		return nil, err
	}
	return utils.SliceToSlice(list, FromEntity), nil
}

func (i *imlServiceService) Check(ctx context.Context, id string, rule map[string]bool) (*Service, error) {
	pro, err := i.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if rule == nil || len(rule) == 0 {
		return pro, nil
	}
	if rule["as_server"] && !pro.AsServer {
		return nil, fmt.Errorf("project %s is not as server", id)
	}
	if rule["as_app"] && !pro.AsApp {
		return nil, fmt.Errorf("project %s is not as app", id)
	}
	return pro, nil
}

func (i *imlServiceService) GetLabels(ctx context.Context, ids ...string) map[string]string {
	if len(ids) == 0 {
		return nil
	}
	list, err := i.store.ListQuery(ctx, "`uuid` in (?)", []interface{}{ids}, "id")
	if err != nil {
		return nil
	}
	return utils.SliceToMapO(list, func(i *service.Service) (string, string) {
		return i.UUID, i.Name
	})
}

func (i *imlServiceService) OnComplete() {
	i.IServiceGet = universally.NewGetSoftDelete[Service, service.Service](i.store, FromEntity)

	i.IServiceDelete = universally.NewSoftDelete[service.Service](i.store)

	i.IServiceCreate = universally.NewCreatorSoftDelete[Create, service.Service](i.store, "service", createEntityHandler, uniquestHandler, labelHandler)

	i.IServiceEdit = universally.NewEdit[Edit, service.Service](i.store, updateHandler, labelHandler)
	auto.RegisterService("service", i)
}

func labelHandler(e *service.Service) []string {
	return []string{e.Name, e.UUID, e.Description}
}
func uniquestHandler(i *Create) []map[string]interface{} {
	return []map[string]interface{}{{"uuid": i.Id}}
}
func createEntityHandler(i *Create) *service.Service {
	now := time.Now()
	return &service.Service{
		Id:          0,
		UUID:        i.Id,
		Name:        i.Name,
		CreateAt:    now,
		UpdateAt:    now,
		Description: i.Description,
		Prefix:      i.Prefix,
		Team:        i.Team,
		ServiceType: i.ServiceType.Int(),
		AsServer:    i.AsServer,
		AsApp:       i.AsApp,
	}
}
func updateHandler(e *service.Service, i *Edit) {
	if i.Name != nil {
		e.Name = *i.Name
	}
	if i.Description != nil {
		e.Description = *i.Description
	}
	if i.ServiceType != nil {
		e.ServiceType = (*i.ServiceType).Int()
	}
	if i.Logo != nil {
		e.Logo = *i.Logo
	}

}
