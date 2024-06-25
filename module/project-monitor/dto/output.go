package project_monitor_dto

type MonitorPartition struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	EnableMonitor bool   `json:"enable_monitor"`
}
