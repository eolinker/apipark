package project_authorization

import (
	"context"
	"reflect"

	"github.com/eolinker/apipark/gateway"

	"github.com/eolinker/go-common/autowire"

	_ "github.com/eolinker/apipark/module/project-authorization/auth-driver/aksk"
	_ "github.com/eolinker/apipark/module/project-authorization/auth-driver/apikey"
	_ "github.com/eolinker/apipark/module/project-authorization/auth-driver/basic"
	_ "github.com/eolinker/apipark/module/project-authorization/auth-driver/jwt"
	_ "github.com/eolinker/apipark/module/project-authorization/auth-driver/oauth2"
	project_authorization_dto "github.com/eolinker/apipark/module/project-authorization/dto"
)

type IProjectAuthorizationModule interface {
	// AddAuthorization 添加项目鉴权信息
	AddAuthorization(ctx context.Context, pid string, info *project_authorization_dto.CreateAuthorization) (*project_authorization_dto.Authorization, error)
	// EditAuthorization 修改项目鉴权信息
	EditAuthorization(ctx context.Context, pid string, aid string, info *project_authorization_dto.EditAuthorization) (*project_authorization_dto.Authorization, error)
	// DeleteAuthorization 删除项目鉴权
	DeleteAuthorization(ctx context.Context, pid string, aid string) error
	// Authorizations 获取项目鉴权列表
	Authorizations(ctx context.Context, pid string) ([]*project_authorization_dto.AuthorizationItem, error)
	// Detail 获取项目鉴权详情（弹窗用）
	Detail(ctx context.Context, pid string, aid string) ([]project_authorization_dto.DetailItem, error)
	// Info 获取项目鉴权详情
	Info(ctx context.Context, pid string, aid string) (*project_authorization_dto.Authorization, error)
}

func init() {
	autowire.Auto[IProjectAuthorizationModule](func() reflect.Value {
		m := new(imlProjectAuthorizationModule)
		gateway.RegisterInitHandleFunc(m.initGateway)
		return reflect.ValueOf(m)
	})
}
