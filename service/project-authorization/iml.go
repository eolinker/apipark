package project_authorization

import (
	"context"
	"time"

	"github.com/eolinker/go-common/utils"

	"github.com/eolinker/go-common/auto"

	"github.com/eolinker/apipark/service/universally"
	"github.com/eolinker/apipark/stores/project"
)

var (
	_ IProjectAuthorizationService = (*imlProjectAuthorizationService)(nil)
)

type imlProjectAuthorizationService struct {
	store project.IAuthorizationStore `autowired:""`
	universally.IServiceGet[Authorization]
	universally.IServiceDelete
	universally.IServiceCreate[CreateAuthorization]
	universally.IServiceEdit[EditAuthorization]
}

func (i *imlProjectAuthorizationService) ListByProject(ctx context.Context, pid ...string) ([]*Authorization, error) {
	w := map[string]interface{}{}
	if len(pid) > 0 {
		w["project"] = pid
	}
	list, err := i.store.List(ctx, w, "update_at desc")
	if err != nil {
		return nil, err
	}
	return utils.SliceToSlice(list, FromEntity), nil
}

func (i *imlProjectAuthorizationService) GetLabels(ctx context.Context, ids ...string) map[string]string {
	if len(ids) == 0 {
		return nil
	}
	list, err := i.store.ListQuery(ctx, "`uuid` in (?)", []interface{}{ids}, "id")
	if err != nil {
		return nil
	}
	return utils.SliceToMapO(list, func(i *project.Authorization) (string, string) {
		return i.UUID, i.Name
	})
}

func (i *imlProjectAuthorizationService) OnComplete() {
	i.IServiceGet = universally.NewGet[Authorization, project.Authorization](i.store, FromEntity)

	i.IServiceDelete = universally.NewDelete[project.Authorization](i.store)

	i.IServiceCreate = universally.NewCreator[CreateAuthorization, project.Authorization](i.store, "project_authorization", createEntityHandler, uniquestHandler, labelHandler)

	i.IServiceEdit = universally.NewEdit[EditAuthorization, project.Authorization](i.store, updateHandler, labelHandler)
	auto.RegisterService("project_authorization", i)
}

func labelHandler(e *project.Authorization) []string {
	return []string{e.Name, e.UUID}
}
func uniquestHandler(i *CreateAuthorization) []map[string]interface{} {
	return []map[string]interface{}{{"uuid": i.UUID}}
}
func createEntityHandler(i *CreateAuthorization) *project.Authorization {
	now := time.Now()
	return &project.Authorization{
		UUID:           i.UUID,
		Name:           i.Name,
		Project:        i.Project,
		Type:           i.Type,
		Position:       i.Position,
		TokenName:      i.TokenName,
		Config:         i.Config,
		ExpireTime:     i.ExpireTime,
		CreateAt:       now,
		UpdateAt:       now,
		HideCredential: i.HideCredential,
	}
}

func updateHandler(e *project.Authorization, i *EditAuthorization) {
	if i.Name != nil {
		e.Name = *i.Name
	}
	if i.Position != nil {
		e.Position = *i.Position
	}
	if i.TokenName != nil {
		e.TokenName = *i.TokenName
	}
	if i.Config != nil {
		e.Config = *i.Config
	}
	if i.ExpireTime != nil {
		e.ExpireTime = *i.ExpireTime
	}
	if i.HideCredential != nil {
		e.HideCredential = *i.HideCredential
	}
	e.UpdateAt = time.Now()
}
