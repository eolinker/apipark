package organization

import (
	"github.com/eolinker/apipark/module/organization"
	organization_dto "github.com/eolinker/apipark/module/organization/dto"
	"github.com/gin-gonic/gin"
)

var (
	_ IOrganizationController = (*imlOrganizationController)(nil)
)

type imlOrganizationController struct {
	module organization.IOrganizationModule `autowired:""`
}

func (c *imlOrganizationController) Simple(ctx *gin.Context) ([]*organization_dto.Simple, error) {
	return c.module.Simple(ctx)

}

func (c *imlOrganizationController) Create(ctx *gin.Context, input *organization_dto.CreateOrganization) (*organization_dto.Detail, error) {
	return c.module.Create(ctx, input)
}

func (c *imlOrganizationController) Edit(ctx *gin.Context, id string, input *organization_dto.EditOrganization) (*organization_dto.Detail, error) {
	return c.module.Edit(ctx, id, input)
}

func (c *imlOrganizationController) Get(ctx *gin.Context, id string) (*organization_dto.Detail, error) {
	return c.module.Get(ctx, id)
}

func (c *imlOrganizationController) Search(ctx *gin.Context, keyword string) ([]*organization_dto.Item, error) {
	return c.module.Search(ctx, keyword)
}

func (c *imlOrganizationController) Delete(ctx *gin.Context, id string) (string, error) {
	return c.module.Delete(ctx, id)
}

func (c *imlOrganizationController) Partitions(ctx *gin.Context, id string) ([]*organization_dto.Partition, error) {
	return c.module.Partitions(ctx, id)
}
