package permit_system

import (
	permit_dto "github.com/eolinker/apipark/module/permit/dto"
	"github.com/eolinker/go-common/autowire"
	"github.com/gin-gonic/gin"
	"reflect"
)

type ISystemPermitController interface {
	Grant(ctx *gin.Context, input *permit_dto.SetGrant) error
	GrantTemplateForTeam(ctx *gin.Context, input *permit_dto.SetGrant) error
	GrantTemplateForProject(ctx *gin.Context, input *permit_dto.SetGrant) error
	Remove(ctx *gin.Context, access string, key string) error
	RemoveTemplateForTeam(ctx *gin.Context, access string, key string) error
	RemoveTemplateForSystem(ctx *gin.Context, access string, key string) error
	List(ctx *gin.Context) ([]*permit_dto.Permission, []*permit_dto.Permission, []*permit_dto.Permission, error)
	Options(ctx *gin.Context, keyword string) ([]*permit_dto.Option, error)
	OptionsForTeamTemplate(ctx *gin.Context, keyword string) ([]*permit_dto.Option, error)
	OptionsForProjectTemplate(ctx *gin.Context, keyword string) ([]*permit_dto.Option, error)
	Permissions(ctx *gin.Context) ([]string, error)
}

func init() {
	autowire.Auto[ISystemPermitController](func() reflect.Value {
		return reflect.ValueOf(new(imlSystemPermitController))
	})
}
