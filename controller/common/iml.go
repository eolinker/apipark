package common

import (
	"github.com/eolinker/apipark/common/version"
	"github.com/gin-gonic/gin"
)

var _ ICommonController = (*imlCommonController)(nil)

type imlCommonController struct{}

func (i imlCommonController) Version(ctx *gin.Context) (string, string, error) {
	return version.Version, version.BuildTime, nil
}
