package project_diff

import (
	"context"
	"github.com/eolinker/apipark/service/project_diff"
	"github.com/eolinker/go-common/autowire"
	"reflect"
)

var (
	_ IProjectDiffModule = (*imlProjectDiff)(nil)
)

type IProjectDiffModule interface {
	Diff(ctx context.Context, projectId string, baseRelease, targetRelease string) (*project_diff.Diff, error)
	DiffForLatest(ctx context.Context, project string, baseRelease string) (*project_diff.Diff, bool, error)
	Out(ctx context.Context, diff *project_diff.Diff) (*DiffOut, error)
}

func init() {
	autowire.Auto[IProjectDiffModule](func() reflect.Value {
		return reflect.ValueOf(new(imlProjectDiff))
	})
}
