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

func init() {
	autowire.Auto[IProjectService](func() reflect.Value {
		return reflect.ValueOf(new(imlProjectService))
	})
}
