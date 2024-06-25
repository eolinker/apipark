package permit_team

import (
	permit_dto "github.com/eolinker/apipark/module/permit/dto"
	"github.com/eolinker/go-common/autowire"
	"github.com/gin-gonic/gin"
	"reflect"
)

type ITeamPermitController interface {
	Grant(ctx *gin.Context, team string, input *permit_dto.SetGrant) error
	Remove(ctx *gin.Context, team string, access string, key string) error
	List(ctx *gin.Context, team string) ([]*permit_dto.Permission, error)
	Options(ctx *gin.Context, team string, keyword string) ([]*permit_dto.Option, error)
	Permissions(ctx *gin.Context, team string) ([]string, error)
}

func init() {
	autowire.Auto[ITeamPermitController](func() reflect.Value {
		return reflect.ValueOf(new(imlTeamPermitController))
	})
}
