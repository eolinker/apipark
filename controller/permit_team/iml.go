package permit_team

import (
	permit_dto "github.com/eolinker/apipark/module/permit/dto"
	"github.com/eolinker/apipark/module/permit/team"
	"github.com/gin-gonic/gin"
)

var (
	_ ITeamPermitController = (*imlTeamPermitController)(nil)
)

type imlTeamPermitController struct {
	teamPermitModule team.ITeamPermitModule `autowired:""`
}

func (c *imlTeamPermitController) Permissions(ctx *gin.Context, team string) ([]string, error) {
	return c.teamPermitModule.Permissions(ctx, team)
}

func (c *imlTeamPermitController) Grant(ctx *gin.Context, team string, input *permit_dto.SetGrant) error {
	return c.teamPermitModule.Grant(ctx, team, input.Access, input.Key)
}

func (c *imlTeamPermitController) Remove(ctx *gin.Context, team string, access string, key string) error {
	return c.teamPermitModule.Remove(ctx, team, access, key)
}

func (c *imlTeamPermitController) List(ctx *gin.Context, team string) ([]*permit_dto.Permission, error) {
	return c.teamPermitModule.TeamAccess(ctx, team)
}

func (c *imlTeamPermitController) Options(ctx *gin.Context, team string, keyword string) ([]*permit_dto.Option, error) {
	return c.teamPermitModule.Options(ctx, team, keyword)
}
