package project_monitor

import (
	"context"

	"github.com/eolinker/go-common/utils"

	"github.com/eolinker/apipark/service/monitor"

	"github.com/eolinker/apipark/service/project"

	"github.com/eolinker/apipark/service/partition"

	project_monitor_dto "github.com/eolinker/apipark/module/project-monitor/dto"
)

var _ IProjectMonitor = (*imlProjectMonitor)(nil)

type imlProjectMonitor struct {
	partitionService        partition.IPartitionService       `autowired:""`
	projectPartitionService project.IProjectPartitionsService `autowired:""`
	monitorService          monitor.IMonitorService           `autowired:""`
}

func (i *imlProjectMonitor) MonitorPartitions(ctx context.Context, pid string) ([]*project_monitor_dto.MonitorPartition, error) {
	partitionIds, err := i.projectPartitionService.GetByProject(ctx, pid)
	if err != nil {
		return nil, err
	}
	monitorMap, err := i.monitorService.MapByPartition(ctx, partitionIds...)
	if err != nil {
		return nil, err
	}
	partitions, err := i.partitionService.Search(ctx, "", map[string]interface{}{
		"uuid": partitionIds,
	}, "create_at asc")
	if err != nil {
		return nil, err
	}
	return utils.SliceToSlice(partitions, func(s *partition.Partition) *project_monitor_dto.MonitorPartition {
		_, ok := monitorMap[s.UUID]
		return &project_monitor_dto.MonitorPartition{
			Id:            s.UUID,
			Name:          s.Name,
			EnableMonitor: ok,
		}
	}), nil
}
