package plugin_partition

import (
	"context"
	"fmt"

	"github.com/eolinker/apipark/gateway"
	"github.com/eolinker/apipark/model/plugin_model"
	"github.com/eolinker/apipark/module/plugin-partition/dto"
	"github.com/eolinker/apipark/service/cluster"
	"github.com/eolinker/apipark/service/partition"
	pluginPartition "github.com/eolinker/apipark/service/plugin-partition"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/go-common/auto"
	"github.com/eolinker/go-common/utils"
)

var (
	_ IPluginPartitionModule = (*imlPluginPartitionModule)(nil)
)

type imlPluginPartitionModule struct {
	service          pluginPartition.IPluginService `autowired:""`
	partitionService partition.IPartitionService    `autowired:""`
	clusterService   cluster.IClusterService        `autowired:""`
}

func (m *imlPluginPartitionModule) UpdateDefine(ctx context.Context, defines []*plugin_model.Define) error {
	err := m.service.SaveDefine(ctx, defines)
	if err != nil {
		return err
	}
	return m.initAllPartition(ctx)
}
func (m *imlPluginPartitionModule) initAllPartition(ctx context.Context) error {
	partitions, err := m.partitionService.List(ctx)
	if err != nil {
		return err
	}

	for _, p := range partitions {
		err := m.initPartition(ctx, p)
		if err != nil {
			log.Warn("init partition:%s %s", p.Name, err.Error())
		}
	}
	return nil
}
func (m *imlPluginPartitionModule) initGateway(ctx context.Context, partitionId string, clientDriver gateway.IClientDriver) error {
	configForPartitions, err := m.service.ListPartition(ctx, partitionId)
	if err != nil {
		return err
	}
	pluginConfigs := utils.SliceToSlice(configForPartitions, func(s *pluginPartition.ConfigPartition) *gateway.PluginConfig {

		return &gateway.PluginConfig{
			Id:     s.Extend,
			Name:   s.Plugin,
			Config: s.Config.Config,
			Status: s.Status.String(),
		}
	})

	return clientDriver.PluginSetting().Set(ctx, pluginConfigs)
}
func (m *imlPluginPartitionModule) GetDefine(ctx context.Context, name string) (*dto.Define, error) {
	define, err := m.service.GetDefine(ctx, name)
	if err != nil {
		return nil, err
	}
	return &dto.Define{
		Name:    define.Name,
		Cname:   define.Cname,
		Desc:    define.Desc,
		Default: define.Config,
		Render:  define.Render,
		Extend:  define.Extend,
	}, nil
}

func (m *imlPluginPartitionModule) Options(ctx context.Context) ([]*dto.PluginOption, error) {
	defines, err := m.service.Defines(ctx, plugin_model.OpenKind)
	if err != nil {
		return nil, err
	}

	return utils.SliceToSlice(defines, func(s *pluginPartition.PluginDefine) *dto.PluginOption {
		return &dto.PluginOption{
			Name:    s.Name,
			Cname:   s.Cname,
			Desc:    s.Desc,
			Default: s.Config,
			Render:  s.Render,
		}
	}), nil
}

func (m *imlPluginPartitionModule) List(ctx context.Context, partition string) ([]*dto.Item, error) {

	configPartitions, err := m.service.ListPartition(ctx, partition, plugin_model.OpenKind)
	if err != nil {
		return nil, err
	}
	return utils.SliceToSlice(configPartitions, func(s *pluginPartition.ConfigPartition) *dto.Item {
		return &dto.Item{

			Name:     s.Plugin,
			Cname:    s.Cname,
			Desc:     s.Desc,
			Extend:   s.Extend,
			Operator: auto.UUIDP(s.Operator),
			Update:   (*auto.TimeLabel)(s.Update),
			Create:   (*auto.TimeLabel)(s.Create),
		}
	}), nil
}

func (m *imlPluginPartitionModule) Get(ctx context.Context, partition string, name string) (config *dto.PluginOutput, render plugin_model.Render, er error) {
	if partition == "" {
		return nil, nil, fmt.Errorf("partition is require")
	}
	cf, define, err := m.service.GetConfig(ctx, partition, name)
	if err != nil {
		return nil, nil, err
	}
	if define.Kind != plugin_model.OpenKind {
		return nil, nil, fmt.Errorf("plugin %s [extend:%s] not support for setting ", name, define.Extend)
	}
	out := &dto.PluginOutput{
		//Partition: auto.UUID(cf.Partition),
		Name:   cf.Plugin,
		Cname:  define.Cname,
		Extend: define.Extend,
		Desc:   define.Desc,
		Status: cf.Status,
		Config: cf.Config,
	}
	if cf.Operator != "" {
		out.Operator = auto.UUIDP(cf.Operator)
	}
	if cf.Create != nil {
		out.Create = (*auto.TimeLabel)(cf.Create)
	}
	if cf.Update != nil {
		out.Update = (*auto.TimeLabel)(cf.Update)
	}
	return out, define.Render, nil

}

func (m *imlPluginPartitionModule) Set(ctx context.Context, partition string, name string, config *dto.PluginSetting) error {

	err := m.service.SetPartition(ctx, partition, name, config.Status, config.Config)
	if err != nil {
		return err
	}
	partitionInfo, err := m.partitionService.Get(ctx, partition)
	if err != nil {
		return err
	}
	return m.initPartition(ctx, partitionInfo)
}

func (m *imlPluginPartitionModule) initPartition(ctx context.Context, partitionInfo *partition.Partition) error {
	client, err := m.clusterService.GatewayClient(ctx, partitionInfo.Cluster)
	if err != nil {
		return err
	}
	defer func() {
		err := client.Close(ctx)
		if err != nil {
			log.Warn("close apinto client:", err)
		}
	}()
	return m.initGateway(ctx, partitionInfo.UUID, client)
}
