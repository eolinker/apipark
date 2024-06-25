package plugin_partition

import (
	"github.com/eolinker/apipark/model/plugin_model"
	"github.com/eolinker/apipark/module/plugin-partition/dto"
	"github.com/eolinker/go-common/autowire"
	"github.com/gin-gonic/gin"
	"reflect"
)

type IPluginPartitionController interface {
	List(ctx *gin.Context, partition string) ([]*dto.Item, error)
	Get(ctx *gin.Context, partition string, name string) (config *dto.PluginOutput, render plugin_model.Render, er error)
	Set(ctx *gin.Context, partition string, name string, config *dto.PluginSetting) error
	Option(ctx *gin.Context, project string) ([]*dto.PluginOption, error)
	Info(ctx *gin.Context, name string) (*dto.Define, error)
}

func init() {
	autowire.Auto[IPluginPartitionController](func() reflect.Value {
		return reflect.ValueOf(new(imlPluginPartitionController))
	})
}
