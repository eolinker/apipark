package partition

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/apipark/gateway/admin"
	"github.com/eolinker/apipark/service/organization"
	"github.com/eolinker/eosc/log"
	"strings"
	"time"

	"github.com/eolinker/go-common/store"

	"github.com/eolinker/apipark/gateway"

	"github.com/eolinker/apipark/module/monitor/driver"

	"gorm.io/gorm"

	"github.com/eolinker/apipark/service/monitor"

	"github.com/google/uuid"

	"github.com/eolinker/ap-account/service/account"
	paritiondto "github.com/eolinker/apipark/module/partition/dto"
	"github.com/eolinker/apipark/service/cluster"
	"github.com/eolinker/apipark/service/partition"
	"github.com/eolinker/go-common/auto"
	"github.com/eolinker/go-common/utils"
)

var (
	_ IPartitionModule = (*imlPartition)(nil)
)

type imlPartition struct {
	partitionService    partition.IPartitionService       `autowired:""`
	organizationService organization.IOrganizationService `autowired:""`
	clusterService      cluster.IClusterService           `autowired:""`
	userNameService     account.IAccountService           `autowired:""`
	monitorService      monitor.IMonitorService           `autowired:""`
	transaction         store.ITransaction                `autowired:""`
}

func (m *imlPartition) CheckCluster(ctx context.Context, address ...string) ([]*paritiondto.Node, error) {
	info, err := admin.Admin(address...).Info(ctx)
	if err != nil {
		return nil, err
	}
	nodesOut := utils.SliceToSlice(info.Nodes, func(i *admin.Node) *paritiondto.Node {
		return &paritiondto.Node{
			Id:       i.Id,
			Name:     i.Name,
			Admins:   i.Admin,
			Peers:    i.Peer,
			Gateways: i.Server,
		}
	})
	nodeStatus(ctx, nodesOut)

	return nodesOut, nil
}

func (m *imlPartition) ResetCluster(ctx context.Context, partitionId string, address string) ([]*paritiondto.Node, error) {
	info, err := m.partitionService.Get(ctx, partitionId)
	if err != nil {
		return nil, err
	}

	nodes, err := m.clusterService.UpdateAddress(ctx, info.Cluster, address)
	if err != nil {
		return nil, err
	}
	err = m.initGateway(ctx, partitionId, info.Cluster)
	if err != nil {
		return nil, err
	}
	nodesOut := utils.SliceToSlice(nodes, func(i *cluster.Node) *paritiondto.Node {
		return &paritiondto.Node{
			Id:       i.Uuid,
			Name:     i.Name,
			Admins:   i.Admin,
			Peers:    i.Peer,
			Gateways: i.Server,
		}
	})

	nodeStatus(ctx, nodesOut)
	return nodesOut, nil
}
func (m *imlPartition) initGateway(ctx context.Context, partitionId, clusterId string) error {
	client, err := m.clusterService.GatewayClient(ctx, clusterId)
	if err != nil {
		return err
	}
	defer func() {
		err := client.Close(ctx)
		if err != nil {
			log.Warn("close apinto client:", err)
		}
	}()
	return gateway.InitGateway(ctx, partitionId, client)
}
func (m *imlPartition) ClusterNodes(ctx context.Context, partitionId string) ([]*paritiondto.Node, error) {
	info, err := m.partitionService.Get(ctx, partitionId)
	if err != nil {
		return nil, err
	}
	nodes, err := m.clusterService.Nodes(ctx, info.Cluster)
	if err != nil {
		return nil, err
	}
	nodesOut := utils.SliceToSlice(nodes, func(i *cluster.Node) *paritiondto.Node {
		return &paritiondto.Node{
			Id:       i.Uuid,
			Name:     i.Name,
			Admins:   i.Admin,
			Peers:    i.Peer,
			Gateways: i.Server,
		}
	})
	nodeStatus(ctx, nodesOut)

	return nodesOut, nil
}

func (m *imlPartition) DeleteMonitorConfig(ctx context.Context, partitionId string) error {
	_, err := m.monitorService.GetByPartition(ctx, partitionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	return m.transaction.Transaction(ctx, func(ctx context.Context) error {
		clusters, err := m.clusterService.ListByClusters(ctx, partitionId)
		if err != nil {
			return err
		}
		err = m.monitorService.DeleteByPartition(ctx, partitionId)
		if err != nil {
			return err
		}
		id := fmt.Sprintf("%s_influxdb", partitionId)
		for _, c := range clusters {
			err := m.dynamicClient(ctx, c.Uuid, "influxdbv2", func(client gateway.IDynamicClient) error {
				return client.Offline(ctx, &gateway.DynamicRelease{
					BasicItem: &gateway.BasicItem{ID: id},
				})
			})
			if err != nil {
				return err
			}

		}
		return nil
	})
}

func (m *imlPartition) MonitorPartitions(ctx context.Context) ([]*paritiondto.MonitorPartition, error) {
	list, err := m.partitionService.Search(ctx, "", nil, "create_at asc")
	if err != nil {
		return nil, err
	}
	partitionIds := utils.SliceToSlice(list, func(i *partition.Partition) string {
		return i.UUID
	})
	monitorMap, err := m.monitorService.MapByPartition(ctx, partitionIds...)
	if err != nil {
		return nil, err
	}
	return utils.SliceToSlice(list, func(i *partition.Partition) *paritiondto.MonitorPartition {
		_, ok := monitorMap[i.UUID]
		return &paritiondto.MonitorPartition{
			Id:            i.UUID,
			Name:          i.Name,
			EnableMonitor: ok,
		}
	}), nil
}

func (m *imlPartition) dynamicClient(ctx context.Context, clusterId string, resource string, f func(gateway.IDynamicClient) error) error {
	client, err := m.clusterService.GatewayClient(ctx, clusterId)
	if err != nil {
		return err
	}
	defer func() {
		err := client.Close(ctx)
		if err != nil {
			log.Warn("close apinto client:", err)
		}
	}()
	dynamic, err := client.Dynamic(resource)
	if err != nil {
		return err
	}
	return f(dynamic)
}

func (m *imlPartition) SaveMonitorConfig(ctx context.Context, partition string, cfg *paritiondto.SaveMonitorConfig) (*paritiondto.MonitorConfig, error) {
	_, err := m.partitionService.Get(ctx, partition)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(cfg.Config)
	err = driver.Check(cfg.Driver, string(data))
	if err != nil {
		return nil, err
	}

	executor, err := driver.CreateExecutor(cfg.Driver, string(data))
	if err != nil {
		return nil, err
	}
	err = executor.Init(ctx)
	if err != nil {
		return nil, err
	}
	clusters, err := m.clusterService.ListByClusters(ctx, partition)
	if err != nil {
		return nil, err
	}
	version := time.Now().Format("20060102150405")
	id := fmt.Sprintf("%s_influxdb", partition)
	for _, c := range clusters {
		err := m.dynamicClient(ctx, c.Uuid, "influxdbv2", func(client gateway.IDynamicClient) error {
			pubCfg := &gateway.DynamicRelease{
				BasicItem: &gateway.BasicItem{
					ID:          id,
					Description: "",
					Version:     version,
					MatchLabels: map[string]string{
						"module": "monitor",
					},
				},
				Attr: map[string]interface{}{
					"org":    cfg.Config["org"],
					"token":  cfg.Config["token"],
					"url":    cfg.Config["addr"],
					"bucket": "apinto",
					"scopes": []string{"monitor"},
				},
			}
			return client.Online(ctx, pubCfg)
		})
		if err != nil {
			return nil, err
		}

	}

	err = m.monitorService.Save(ctx, &monitor.SaveMonitor{
		Partition: partition,
		Driver:    cfg.Driver,
		Config:    string(data),
	})
	if err != nil {
		return nil, err
	}

	return m.GetMonitorConfig(ctx, partition)
}

func (m *imlPartition) GetMonitorConfig(ctx context.Context, partition string) (*paritiondto.MonitorConfig, error) {
	_, err := m.partitionService.Get(ctx, partition)
	if err != nil {
		return nil, err

	}
	info, err := m.monitorService.GetByPartition(ctx, partition)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return &paritiondto.MonitorConfig{
			Driver: "influxdb-v2",
			Config: map[string]interface{}{},
		}, nil
	}
	cfg := make(map[string]interface{})
	err = json.Unmarshal([]byte(info.Config), &cfg)
	if err != nil {
		return nil, err
	}
	return &paritiondto.MonitorConfig{
		Driver: info.Driver,
		Config: cfg,
	}, nil
}

func (m *imlPartition) CreatePartition(ctx context.Context, create *paritiondto.Create) (*paritiondto.Detail, error) {
	if create.Id == "" {
		create.Id = uuid.New().String()
	}
	if create.Name == "" {
		return nil, errors.New("name is empty")
	}
	clusterId := ""
	err := m.transaction.Transaction(ctx, func(ctx context.Context) error {
		clusterInfo, err := m.clusterService.Create(ctx, create.Id, create.Id, create.Description, create.ManagerAddress)
		if err != nil {
			return err
		}
		if create.Prefix != "" {
			create.Prefix = fmt.Sprintf("/%s", strings.TrimPrefix(create.Prefix, "/"))
		}
		clusterId = clusterInfo.Uuid
		return m.partitionService.Create(ctx, &partition.CreatePartition{
			Uuid:    create.Id,
			Name:    create.Name,
			Resume:  create.Description,
			Prefix:  create.Prefix,
			Url:     create.Url,
			Cluster: clusterInfo.Uuid,
		})
	})
	if err != nil {
		return nil, err
	}
	err = m.initGateway(ctx, create.Id, clusterId)
	if err != nil {
		return nil, err
	}
	return m.Get(ctx, create.Id)
}

func (m *imlPartition) Search(ctx context.Context, keyword string) ([]*paritiondto.Item, error) {
	partitions, err := m.partitionService.Search(ctx, keyword, nil)
	if err != nil {
		return nil, err
	}
	countMap, err := m.clusterService.CountByPartition(ctx)
	if err != nil {
		return nil, err
	}
	items := utils.SliceToSlice(partitions, func(i *partition.Partition) *paritiondto.Item {

		return &paritiondto.Item{
			Creator:     auto.UUID(i.Creator),
			Updater:     auto.UUID(i.Updater),
			Id:          i.UUID,
			Name:        i.Name,
			Description: i.Resume,
			ClusterNum:  countMap[i.UUID],
			CreateTime:  auto.TimeLabel(i.CreateTime),
			UpdateTime:  auto.TimeLabel(i.UpdateTime),
		}
	})
	if len(items) > 0 {
		counts, err := m.clusterService.CountByPartition(ctx)
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			item.ClusterNum = counts[item.Id]
		}
	}

	return items, nil
}

func (m *imlPartition) Get(ctx context.Context, id string) (*paritiondto.Detail, error) {
	pm, err := m.partitionService.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	oDetails, err := m.organizationService.Search(ctx, "")
	if err != nil {
		return nil, err
	}
	canDelete := true
	for _, o := range oDetails {
		for _, p := range o.Partitions {
			if p == id {
				canDelete = false
				break
			}
		}
		if !canDelete {
			break
		}
	}

	pd := &paritiondto.Detail{
		Creator:     auto.UUID(pm.Creator),
		Updater:     auto.UUID(pm.Updater),
		Id:          pm.UUID,
		Name:        pm.Name,
		Description: pm.Resume,
		Prefix:      pm.Prefix,
		CreateTime:  auto.TimeLabel(pm.CreateTime),
		UpdateTime:  auto.TimeLabel(pm.UpdateTime),
		CanDelete:   canDelete,
	}
	return pd, nil
}

func (m *imlPartition) Update(ctx context.Context, id string, edit *paritiondto.Edit) (*paritiondto.Detail, error) {
	err := m.partitionService.Save(ctx, id, &partition.EditPartition{
		Name:   edit.Name,
		Resume: edit.Description,
		Prefix: edit.Prefix,
		Url:    edit.Url,
	})
	if err != nil {
		return nil, err
	}
	return m.Get(ctx, id)
}

func (m *imlPartition) Delete(ctx context.Context, id string) error {
	return m.transaction.Transaction(ctx, func(ctx context.Context) error {
		info, err := m.partitionService.Get(ctx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}
		err = m.clusterService.Delete(ctx, info.Cluster)
		if err != nil {
			return err
		}
		return m.partitionService.Delete(ctx, id)
	})

}

func (m *imlPartition) Simple(ctx context.Context) ([]*paritiondto.Simple, error) {
	pm, err := m.partitionService.Search(ctx, "", nil)
	if err != nil {
		return nil, err
	}
	pd := utils.SliceToSlice(pm, func(i *partition.Partition) *paritiondto.Simple {
		return &paritiondto.Simple{
			Id:   i.UUID,
			Name: i.Name,
		}
	})
	return pd, nil
}

func (m *imlPartition) SimpleByIds(ctx context.Context, ids []string) ([]*paritiondto.Simple, error) {
	pm, err := m.partitionService.Search(ctx, "", map[string]interface{}{
		"uuid": ids,
	})
	if err != nil {
		return nil, err
	}
	pd := utils.SliceToSlice(pm, func(i *partition.Partition) *paritiondto.Simple {
		return &paritiondto.Simple{
			Id:   i.UUID,
			Name: i.Name,
		}
	})
	return pd, nil

}
func (m *imlPartition) SimpleWithCluster(ctx context.Context) ([]*paritiondto.SimpleWithCluster, error) {
	pm, err := m.partitionService.Search(ctx, "", nil)
	if err != nil {
		return nil, err
	}

	clusterList, err := m.clusterService.List(ctx)
	if err != nil {
		return nil, err
	}

	clusterMap := utils.SliceToMapArrayO(clusterList, func(i *cluster.Cluster) (string, *paritiondto.Cluster) {
		return i.Partition, &paritiondto.Cluster{
			Id:          i.Uuid,
			Name:        i.Name,
			Description: i.Resume,
		}
	})
	pd := utils.SliceToSlice(pm, func(i *partition.Partition) *paritiondto.SimpleWithCluster {
		return &paritiondto.SimpleWithCluster{
			Id:       i.UUID,
			Name:     i.Name,
			Clusters: clusterMap[i.UUID],
		}
	})
	return pd, nil
}
