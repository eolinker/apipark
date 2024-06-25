package project

import (
	"context"
	"reflect"

	"github.com/eolinker/apipark/service/universally"
	"github.com/eolinker/go-common/autowire"
)

type IProjectService interface {
	universally.IServiceGet[Project]
	universally.IServiceDelete
	universally.IServiceCreate[CreateProject]
	universally.IServiceEdit[EditProject]
	CountByTeam(ctx context.Context, keyword string) (map[string]int64, error)
	CountTeam(ctx context.Context, teamID string, keyword string) (int64, error)
	CheckProject(ctx context.Context, pid string, rule map[string]bool) (*Project, error)
	AppList(ctx context.Context, appIds ...string) ([]*Project, error)
}

type IProjectPartitionsService interface {
	ListByProject(ctx context.Context, projectIDs ...string) ([]*Partition, error)
	ListByPartition(ctx context.Context, partition ...string) ([]*Partition, error)
	Save(ctx context.Context, projectID string, partitions []string) error
	Delete(ctx context.Context, projectID string) error
	GetByProject(ctx context.Context, projectID string) ([]string, error)
}

func init() {
	autowire.Auto[IProjectService](func() reflect.Value {
		return reflect.ValueOf(new(imlProjectService))
	})
	autowire.Auto[IProjectPartitionsService](func() reflect.Value {
		return reflect.ValueOf(new(imlProjectPartitionService))
	})
}
