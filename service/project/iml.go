package project

import (
	"context"
	"fmt"
	"time"

	"github.com/eolinker/go-common/utils"

	"github.com/eolinker/go-common/auto"

	"github.com/eolinker/apipark/service/universally"
	"github.com/eolinker/apipark/stores/project"
)

var (
	_ IProjectService = (*imlProjectService)(nil)
)

type imlProjectService struct {
	projectStore project.IProjectStore `autowired:""`
	universally.IServiceGet[Project]
	universally.IServiceDelete
	universally.IServiceCreate[CreateProject]
	universally.IServiceEdit[EditProject]
}

func (i *imlProjectService) AppList(ctx context.Context, appIds ...string) ([]*Project, error) {
	w := make(map[string]interface{})
	if len(appIds) > 0 {
		w["uuid"] = appIds
	}
	w["as_app"] = true
	list, err := i.projectStore.List(ctx, w)
	if err != nil {
		return nil, err
	}
	return utils.SliceToSlice(list, FromEntity), nil
}

func (i *imlProjectService) CheckProject(ctx context.Context, pid string, rule map[string]bool) (*Project, error) {
	pro, err := i.Get(ctx, pid)
	if err != nil {
		return nil, err
	}
	if rule == nil || len(rule) == 0 {
		return pro, nil
	}
	if rule["as_server"] && !pro.AsServer {
		return nil, fmt.Errorf("project %s is not as server", pid)
	}
	if rule["as_app"] && !pro.AsApp {
		return nil, fmt.Errorf("project %s is not as app", pid)
	}
	return pro, nil
}

func (i *imlProjectService) CountTeam(ctx context.Context, teamID string, keyword string) (int64, error) {
	counts, err := i.projectStore.CountByGroup(ctx, keyword, map[string]interface{}{
		"team":      teamID,
		"as_server": true,
	}, "team")
	if err != nil {
		return 0, err
	}
	return counts[teamID], nil
}

func (i *imlProjectService) CountByTeam(ctx context.Context, keyword string) (map[string]int64, error) {
	return i.projectStore.CountByGroup(ctx, keyword, map[string]interface{}{"as_server": true}, "team")
}

//func (i *imlProjectService) ListByProject(ctx context.Context, projectId string) ([]*Project, error) {
//	list, err := i.projectStore.List(ctx, map[string]interface{}{
//		"project": projectId,
//	}, "update_at desc")
//	if err != nil {
//		return nil, err
//	}
//	return utils.SliceToSlice(list, FromEntity), nil
//}

func (i *imlProjectService) GetLabels(ctx context.Context, ids ...string) map[string]string {
	if len(ids) == 0 {
		return nil
	}
	list, err := i.projectStore.ListQuery(ctx, "`uuid` in (?)", []interface{}{ids}, "id")
	if err != nil {
		return nil
	}
	return utils.SliceToMapO(list, func(i *project.Project) (string, string) {
		return i.UUID, i.Name
	})
}

func (i *imlProjectService) OnComplete() {
	i.IServiceGet = universally.NewGetSoftDelete[Project, project.Project](i.projectStore, FromEntity)

	i.IServiceDelete = universally.NewSoftDelete[project.Project](i.projectStore)

	i.IServiceCreate = universally.NewCreatorSoftDelete[CreateProject, project.Project](i.projectStore, "project", createEntityHandler, uniquestHandler, labelHandler)

	i.IServiceEdit = universally.NewEdit[EditProject, project.Project](i.projectStore, updateHandler, labelHandler)
	auto.RegisterService("project", i)
}

func labelHandler(e *project.Project) []string {
	return []string{e.Name, e.UUID, e.Description}
}
func uniquestHandler(i *CreateProject) []map[string]interface{} {
	return []map[string]interface{}{{"uuid": i.Id}}
}
func createEntityHandler(i *CreateProject) *project.Project {
	now := time.Now()
	return &project.Project{
		Id:           0,
		UUID:         i.Id,
		Name:         i.Name,
		CreateAt:     now,
		UpdateAt:     now,
		Description:  i.Description,
		Prefix:       i.Prefix,
		Team:         i.Team,
		Organization: i.Organization,
		Master:       i.Master,
		AsServer:     i.AsServer,
		AsApp:        i.AsApp,
	}
}
func updateHandler(e *project.Project, i *EditProject) {
	if i.Name != nil {
		e.Name = *i.Name
	}
	if i.Description != nil {
		e.Description = *i.Description
	}
	if i.Master != nil {
		e.Master = *i.Master
	}
}

var _ IProjectPartitionsService = (*imlProjectPartitionService)(nil)

type imlProjectPartitionService struct {
	store project.IProjectPartitionStore `autowired:""`
}

func (i *imlProjectPartitionService) ListByPartition(ctx context.Context, partitionIds ...string) ([]*Partition, error) {
	w := make(map[string]interface{})
	if len(partitionIds) > 0 {
		w["partition"] = partitionIds
	}
	list, err := i.store.List(ctx, w)
	if err != nil {
		return nil, err
	}
	return utils.SliceToSlice(list, func(t *project.Partition) *Partition {
		return &Partition{
			Project:   t.Project,
			Partition: t.Partition,
		}
	}), nil
}

func (i *imlProjectPartitionService) Delete(ctx context.Context, projectID string) error {
	_, err := i.store.DeleteWhere(ctx, map[string]interface{}{
		"project": projectID,
	})
	return err
}

func (i *imlProjectPartitionService) ListByProject(ctx context.Context, projectIDs ...string) ([]*Partition, error) {
	w := make(map[string]interface{})
	if len(projectIDs) > 0 {
		w["project"] = projectIDs
	}
	list, err := i.store.List(ctx, w)
	if err != nil {
		return nil, err
	}
	return utils.SliceToSlice(list, func(t *project.Partition) *Partition {
		return &Partition{
			Project:   t.Project,
			Partition: t.Partition,
		}
	}), nil
}

func (i *imlProjectPartitionService) Save(ctx context.Context, projectID string, partitions []string) error {
	if len(partitions) == 0 {
		return fmt.Errorf("partitions is empty")
	}
	return i.store.Transaction(ctx, func(ctx context.Context) error {
		_, err := i.store.DeleteWhere(ctx, map[string]interface{}{
			"project": projectID,
		})
		if err != nil {
			return err
		}
		now := time.Now()
		return i.store.Insert(ctx, utils.SliceToSlice(partitions, func(t string) *project.Partition {
			return &project.Partition{
				Project:    projectID,
				Partition:  t,
				CreateTime: now,
			}
		})...)
	})

}

func (i *imlProjectPartitionService) GetByProject(ctx context.Context, projectID string) ([]string, error) {
	list, err := i.store.List(ctx, map[string]interface{}{
		"project": projectID,
	})
	if err != nil {
		return nil, err
	}
	return utils.SliceToSlice(list, func(t *project.Partition) string {
		return t.Partition
	}), nil
}
