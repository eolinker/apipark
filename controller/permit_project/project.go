package permit_project

import (
	permit_dto "github.com/eolinker/apipark/module/permit/dto"
	"github.com/eolinker/go-common/autowire"
	"github.com/gin-gonic/gin"
	"reflect"
)

type IProjectPermitController interface {
	Grant(ctx *gin.Context, project string, input *permit_dto.SetGrant) error
	Remove(ctx *gin.Context, project string, access string, key string) error
	List(ctx *gin.Context, project string) ([]*permit_dto.Permission, error)
	Options(ctx *gin.Context, project string, keyword string) ([]*permit_dto.Option, error)
	Permissions(ctx *gin.Context, project string) ([]string, error)
}

func init() {
	autowire.Auto[IProjectPermitController](func() reflect.Value {
		return reflect.ValueOf(new(imlProjectPermitController))
	})
}
