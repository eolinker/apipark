package team

import (
	"context"
	permit_dto "github.com/eolinker/apipark/module/permit/dto"
	"github.com/eolinker/go-common/autowire"
	"github.com/gin-gonic/gin"
	"reflect"
)

const (
	accessGroup = "team"
)

type ITeamPermitModule interface {
	Grant(ctx *gin.Context, teamId, access string, key string) error
	Remove(ctx *gin.Context, teamId string, access string, key string) error
	TeamAccess(ctx *gin.Context, teamId string) ([]*permit_dto.Permission, error)
	Options(ctx context.Context, teamId string, keyword string) ([]*permit_dto.Option, error)
	Permissions(ctx context.Context, teamId string) ([]string, error)
}

func init() {
	var m *imlTeamPermitModule

	autowire.Auto[ITeamPermitModule](func() reflect.Value {
		if m == nil {
			m = new(imlTeamPermitModule)
		}
		return reflect.ValueOf(m)
	})

}
