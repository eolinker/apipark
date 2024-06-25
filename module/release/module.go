package release

import (
	"context"
	"github.com/eolinker/apipark/module/release/dto"
	"github.com/eolinker/apipark/service/project_diff"
	"github.com/eolinker/go-common/autowire"
	"reflect"
)

type IReleaseModule interface {
	Create(ctx context.Context, project string, input *dto.CreateInput) (string, error)
	Detail(ctx context.Context, project string, id string) (*dto.Detail, error)
	List(ctx context.Context, project string) ([]*dto.Release, error)
	Delete(ctx context.Context, project string, id string) error
	Preview(ctx context.Context, project string) (*dto.Release, *project_diff.Diff, bool, error)
}

func init() {
	autowire.Auto[IReleaseModule](func() reflect.Value {
		return reflect.ValueOf(new(imlReleaseModule))
	})
}
