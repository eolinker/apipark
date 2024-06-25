package permit_system

import (
	permit_identity "github.com/eolinker/apipark/middleware/permit/identity"
	"github.com/eolinker/apipark/module/permit/dto"
	"github.com/eolinker/apipark/module/permit/system"
	"github.com/eolinker/go-common/autowire"
	"github.com/gin-gonic/gin"
)

var (
	_ ISystemPermitController = (*imlSystemPermitController)(nil)
	_ autowire.Complete       = (*imlSystemPermitController)(nil)
)

type imlSystemPermitController struct {
	systemPermitModule system.ISystemPermitModule `autowired:""`

	identitySystem permit_identity.IdentitySystemService `autowired:""`
}

func (c *imlSystemPermitController) Permissions(ctx *gin.Context) ([]string, error) {
	return c.systemPermitModule.Permissions(ctx)
}

func (c *imlSystemPermitController) RemoveTemplateForTeam(ctx *gin.Context, access string, key string) error {
	return c.systemPermitModule.RemoveTeamTemplateAccess(ctx, access, key)
}

func (c *imlSystemPermitController) RemoveTemplateForSystem(ctx *gin.Context, access string, key string) error {
	return c.systemPermitModule.RemoveProjectTemplateAccess(ctx, access, key)
}

func (c *imlSystemPermitController) GrantTemplateForTeam(ctx *gin.Context, input *permit_dto.SetGrant) error {
	return c.systemPermitModule.GrantTemplateForTeam(ctx, input.Access, input.Key)
}

func (c *imlSystemPermitController) GrantTemplateForProject(ctx *gin.Context, input *permit_dto.SetGrant) error {
	return c.systemPermitModule.GrantTemplateForProject(ctx, input.Access, input.Key)
}

func (c *imlSystemPermitController) OptionsForTeamTemplate(ctx *gin.Context, keyword string) ([]*permit_dto.Option, error) {
	return c.systemPermitModule.OptionsForTeamTemplate(ctx, keyword)
}

func (c *imlSystemPermitController) OptionsForProjectTemplate(ctx *gin.Context, keyword string) ([]*permit_dto.Option, error) {
	return c.systemPermitModule.OptionsForProjectTemplate(ctx, keyword)
}

func (c *imlSystemPermitController) OnComplete() {

}

func (c *imlSystemPermitController) Grant(ctx *gin.Context, input *permit_dto.SetGrant) error {
	return c.systemPermitModule.GrantSystem(ctx, input.Access, input.Key)
}

func (c *imlSystemPermitController) Remove(ctx *gin.Context, access string, key string) error {
	return c.systemPermitModule.RemoveSystemAccess(ctx, access, key)
}

func (c *imlSystemPermitController) List(ctx *gin.Context) ([]*permit_dto.Permission, []*permit_dto.Permission, []*permit_dto.Permission, error) {
	systemAccess, err := c.systemPermitModule.SystemAccess(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	teamAccess, err := c.systemPermitModule.TeamAccess(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	projectAccess, err := c.systemPermitModule.ProjectAccess(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	return systemAccess, teamAccess, projectAccess, err
}

func (c *imlSystemPermitController) Options(ctx *gin.Context, keyword string) ([]*permit_dto.Option, error) {
	return c.systemPermitModule.OptionsForSystem(ctx, keyword)
}
