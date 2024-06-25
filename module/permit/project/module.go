package project

import (
	"context"
	permit_dto "github.com/eolinker/apipark/module/permit/dto"
	"github.com/eolinker/go-common/autowire"
	"github.com/gin-gonic/gin"
	"reflect"
)

const (
	accessGroup = "project"
)

type IProjectPermitModule interface {
	Grant(ctx *gin.Context, projectId, access string, key string) error
	Remove(ctx *gin.Context, projectId string, access string, key string) error
	ProjectAccess(ctx *gin.Context, projectId string) ([]*permit_dto.Permission, error)
	Options(ctx context.Context, projectId string, keyword string) ([]*permit_dto.Option, error)
	Permissions(ctx *gin.Context, projectId string) ([]string, error)
}

func init() {

	autowire.Auto[IProjectPermitModule](func() reflect.Value {
		m := new(imlProjectPermitModule)

		return reflect.ValueOf(m)
	})

}
