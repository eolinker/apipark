package permit_middleware

import (
	"errors"
	permit_identity "github.com/eolinker/apipark/middleware/permit/identity"
	permit_type "github.com/eolinker/apipark/service/permit-type"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/go-common/autowire"
	"github.com/eolinker/go-common/permit"
	"github.com/eolinker/go-common/pm3"
	"github.com/eolinker/go-common/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"reflect"
)

var (
	checkSort = []string{permit_identity.ProjectGroup, permit_identity.TeamGroup, permit_identity.SystemGroup}
)

type IPermitMiddleware interface {
	pm3.IMiddleware
}

func init() {
	autowire.Auto[IPermitMiddleware](func() reflect.Value {
		return reflect.ValueOf(new(PermitMiddleware))
	})
}

var (
	_ IPermitMiddleware = (*PermitMiddleware)(nil)
)

type PermitMiddleware struct {
	permitService permit.IPermit `autowired:""`
}

func (p *PermitMiddleware) Sort() int {
	return 99
}

func (p *PermitMiddleware) Check(method string, path string) (bool, []gin.HandlerFunc) {
	// 当前路径是否有配置权限
	accessRules, has := permit.GetPathRule(method, path)

	if !has || len(accessRules) == 0 {
		return false, nil
	}

	return true, []gin.HandlerFunc{
		func(ginCtx *gin.Context) {
			userId := utils.UserId(ginCtx)
			if userId == "" {
				ginCtx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": http.StatusForbidden, "msg": "not login", "success": "fail"})
				ginCtx.Abort()
				return
			}
			if userId == "admin" {
				// 超级管理员不校验
				return
			}

			for _, group := range checkSort {
				accessList, has := accessRules[group]
				if !has {
					// 当前分组没有配置权限
					continue
				}
				domainHandler, has := permit.SelectDomain(group)
				if !has {
					// 当前分组没有配置身份handler
					continue
				}
				domains, myIdentity, ok := domainHandler(ginCtx)
				if !ok {
					// 无效的身份域
					//ginCtx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": http.StatusForbidden, "msg": "domain not found", "success": "fail"})
					//ginCtx.Abort()
					//return
					continue
				}
				myIdentitySet := utils.NewSet(myIdentity...)
				myIdentitySet.Set(permit_type.AnyOne.Key)
				for _, domain := range domains {
					grantsOfAccess, err := p.permitService.GrantForDomain(ginCtx, domain)
					if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
						//ginCtx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": http.StatusForbidden, "msg": "domain not found", "success": "fail"})
						//ginCtx.Abort()
						//return
						continue
					}
					if len(grantsOfAccess) == 0 {
						// 当前域没有配置权限
						continue
					}
					checkCount := 0
					for _, access := range accessList {
						grants, ok := grantsOfAccess[access]
						if !ok {
							// 当前域没有配置目标权限
							continue
						}

						for _, grant := range grants {
							if myIdentitySet.Has(grant) {
								// 当前用户有权限
								return
							}
						}
						checkCount++

					}
					if checkCount > 0 {
						// 当前域有权限配置,且当前用户没有目标权限,则跳过当前group接下来的domain
						break
					}
				}
			}
			//所有group都校验不通过
			log.DebugF("no permission:%s", ginCtx.FullPath())
			ginCtx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": http.StatusForbidden, "msg": "no permission", "success": "fail"})
		},
	}
}

func (p *PermitMiddleware) Name() string {
	return "permit"
}
