package organization

import (
	organization_dto "github.com/eolinker/apipark/module/organization/dto"
	"github.com/eolinker/go-common/autowire"
	"github.com/gin-gonic/gin"
	"reflect"
)

type IOrganizationController interface {
	Create(ctx *gin.Context, input *organization_dto.CreateOrganization) (*organization_dto.Detail, error)
	Edit(ctx *gin.Context, id string, input *organization_dto.EditOrganization) (*organization_dto.Detail, error)
	Get(ctx *gin.Context, id string) (*organization_dto.Detail, error)
	Search(ctx *gin.Context, keyword string) ([]*organization_dto.Item, error)
	Delete(ctx *gin.Context, id string) (string, error)
	Partitions(ctx *gin.Context, id string) ([]*organization_dto.Partition, error)
	Simple(ctx *gin.Context) ([]*organization_dto.Simple, error)
}

func init() {
	autowire.Auto[IOrganizationController](func() reflect.Value {
		return reflect.ValueOf(new(imlOrganizationController))
	})
}
