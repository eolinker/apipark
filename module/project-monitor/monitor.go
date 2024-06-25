package project_monitor

import (
	"context"
	"reflect"

	"github.com/eolinker/go-common/autowire"

	project_monitor_dto "github.com/eolinker/apipark/module/project-monitor/dto"
)

type IProjectMonitor interface {
	MonitorPartitions(ctx context.Context, pid string) ([]*project_monitor_dto.MonitorPartition, error)
}

func init() {
	autowire.Auto[IProjectMonitor](func() reflect.Value {
		return reflect.ValueOf(new(imlProjectMonitor))
	})
}
