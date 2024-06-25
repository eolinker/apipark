package permit_project

import (
	permit_dto "github.com/eolinker/apipark/module/permit/dto"
	"github.com/eolinker/apipark/module/permit/project"
	"github.com/gin-gonic/gin"
)

var (
	_ IProjectPermitController = (*imlProjectPermitController)(nil)
)

type imlProjectPermitController struct {
	projectPermitModule project.IProjectPermitModule `autowired:""`
}

func (c *imlProjectPermitController) Permissions(ctx *gin.Context, project string) ([]string, error) {
	return c.projectPermitModule.Permissions(ctx, project)
}

func (c *imlProjectPermitController) Grant(ctx *gin.Context, project string, input *permit_dto.SetGrant) error {
	return c.projectPermitModule.Grant(ctx, project, input.Access, input.Key)
}

func (c *imlProjectPermitController) Remove(ctx *gin.Context, project string, access string, key string) error {
	return c.projectPermitModule.Remove(ctx, project, access, key)

}

func (c *imlProjectPermitController) List(ctx *gin.Context, project string) ([]*permit_dto.Permission, error) {

	return c.projectPermitModule.ProjectAccess(ctx, project)

}

func (c *imlProjectPermitController) Options(ctx *gin.Context, project string, keyword string) ([]*permit_dto.Option, error) {
	return c.projectPermitModule.Options(ctx, project, keyword)
}
