package service_dto

import "github.com/eolinker/go-common/auto"

type ServiceItem struct {
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	Partition   []auto.Label   `json:"partition" aolabel:"partition"`
	ServiceType string         `json:"service_type"`
	ApiNum      int64          `json:"api_num"`
	Status      string         `json:"status"`
	CreateTime  auto.TimeLabel `json:"create_time"`
	UpdateTime  auto.TimeLabel `json:"update_time"`
	CanDelete   bool           `json:"can_delete"`
}

type Service struct {
	Id          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Logo        string       `json:"logo"`
	ServiceType string       `json:"service_type"`
	Team        auto.Label   `json:"team" aolabel:"team"`
	Project     auto.Label   `json:"project" aolabel:"project"`
	Catalogue   auto.Label   `json:"group" aolabel:"catalogue"`
	Partition   []auto.Label `json:"partition" aolabel:"partition"`
	Tags        []auto.Label `json:"tags" aolabel:"tag"`
	Status      string       `json:"status"`
}

type ServiceDoc struct {
	Id         string         `json:"id"`
	Name       string         `json:"name"`
	Doc        string         `json:"doc"`
	Creator    auto.Label     `json:"creator" aolabel:"user"`
	CreateTime auto.TimeLabel `json:"create_time"`
	Updater    auto.Label     `json:"updater" aolabel:"user"`
	UpdateTime auto.TimeLabel `json:"update_time"`
}

type SaveServiceDoc struct {
	Doc string `json:"doc"`
}

type ServiceApi struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

type SimpleItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
