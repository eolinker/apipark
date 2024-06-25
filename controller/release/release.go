package release

import (
	"github.com/eolinker/apipark/module/project_diff"
	"github.com/eolinker/apipark/module/release/dto"
	"github.com/eolinker/go-common/autowire"
	"github.com/gin-gonic/gin"
	"reflect"
)

type IReleaseController interface {
	Create(ctx *gin.Context, project string, input *dto.CreateInput) error
	Delete(ctx *gin.Context, project string, id string) error
	Detail(ctx *gin.Context, project string, id string) (*dto.Detail, error)
	List(ctx *gin.Context, project string) ([]*dto.Release, error)
	Preview(ctx *gin.Context, project string) (*dto.Release, *project_diff.DiffOut, bool, error)
}

func init() {
	autowire.Auto[IReleaseController](func() reflect.Value {
		return reflect.ValueOf(new(imlReleaseController))
	})
}
