package project_authorization

import (
	project_authorization "github.com/eolinker/apipark/module/project-authorization"
	"github.com/gin-gonic/gin"

	project_authorization_dto "github.com/eolinker/apipark/module/project-authorization/dto"
)

var _ IProjectAuthorizationController = (*imlProjectAuthorizationController)(nil)

type imlProjectAuthorizationController struct {
	module project_authorization.IProjectAuthorizationModule `autowired:""`
}

func (i *imlProjectAuthorizationController) AddAuthorization(ctx *gin.Context, pid string, info *project_authorization_dto.CreateAuthorization) (*project_authorization_dto.Authorization, error) {
	return i.module.AddAuthorization(ctx, pid, info)
}

func (i *imlProjectAuthorizationController) EditAuthorization(ctx *gin.Context, pid string, aid string, info *project_authorization_dto.EditAuthorization) (*project_authorization_dto.Authorization, error) {
	return i.module.EditAuthorization(ctx, pid, aid, info)
}

func (i *imlProjectAuthorizationController) DeleteAuthorization(ctx *gin.Context, pid string, aid string) error {
	return i.module.DeleteAuthorization(ctx, pid, aid)
}

func (i *imlProjectAuthorizationController) Authorizations(ctx *gin.Context, pid string) ([]*project_authorization_dto.AuthorizationItem, error) {
	return i.module.Authorizations(ctx, pid)
}

func (i *imlProjectAuthorizationController) Detail(ctx *gin.Context, pid string, aid string) ([]project_authorization_dto.DetailItem, error) {
	return i.module.Detail(ctx, pid, aid)
}

func (i *imlProjectAuthorizationController) Info(ctx *gin.Context, pid string, aid string) (*project_authorization_dto.Authorization, error) {
	return i.module.Info(ctx, pid, aid)
}
