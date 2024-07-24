package upstream

import (
	"reflect"

	"github.com/eolinker/go-common/autowire"

	upstream_dto "github.com/eolinker/apipark/module/upstream/dto"
	"github.com/gin-gonic/gin"
)

type IUpstreamController interface {
	Get(ctx *gin.Context, pid string) (upstream_dto.UpstreamConfig, error)
	Save(ctx *gin.Context, pid string, upstream *upstream_dto.UpstreamConfig) (upstream_dto.UpstreamConfig, error)
}

func init() {
	autowire.Auto[IUpstreamController](func() reflect.Value {
		return reflect.ValueOf(new(imlUpstreamController))
	})
}
