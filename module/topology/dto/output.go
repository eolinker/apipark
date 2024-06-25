package topology_dto

import "github.com/eolinker/go-common/auto"

type ProjectItem struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	InvokeServices []string `json:"invoke_services"`
	IsServer       bool     `json:"is_server"`
	IsApp          bool     `json:"is_app"`
}

type ServiceItem struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Project string `json:"project"`
}

type TopologyItem struct {
	Project  auto.Label   `json:"project" aolabel:"project"`
	Services []auto.Label `json:"services" aolabel:"service"`
}
