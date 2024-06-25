package plugin_partition

import (
	"github.com/eolinker/apipark/model/plugin_model"
	plugin_partition "github.com/eolinker/apipark/module/plugin-partition"
	"github.com/eolinker/apipark/module/plugin-partition/dto"
	"github.com/gin-gonic/gin"
)

var (
	_ IPluginPartitionController = (*imlPluginPartitionController)(nil)
)

type imlPluginPartitionController struct {
	module plugin_partition.IPluginPartitionModule `autowired:""`
}

func (i *imlPluginPartitionController) Info(ctx *gin.Context, name string) (*dto.Define, error) {
	return i.module.GetDefine(ctx, name)
}

func (i *imlPluginPartitionController) Option(ctx *gin.Context, project string) ([]*dto.PluginOption, error) {
	return i.module.Options(ctx)
}

func (i *imlPluginPartitionController) List(ctx *gin.Context, partition string) ([]*dto.Item, error) {
	return i.module.List(ctx, partition)
}

func (i *imlPluginPartitionController) Get(ctx *gin.Context, partition string, name string) (config *dto.PluginOutput, render plugin_model.Render, er error) {
	return i.module.Get(ctx, partition, name)
}

func (i *imlPluginPartitionController) Set(ctx *gin.Context, partition string, name string, config *dto.PluginSetting) error {
	return i.module.Set(ctx, partition, name, config)
}
