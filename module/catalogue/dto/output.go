package catalogue_dto

import "github.com/eolinker/go-common/auto"

type Item struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Children []*Item `json:"children"`
}

type ServiceItem struct {
	Id        string       `json:"id"`
	Name      string       `json:"name"`
	Tags      []auto.Label `json:"tags" aolabel:"tag"`
	Catalogue auto.Label   `json:"catalogue" aolabel:"catalogue"`
	//Partition     []auto.Label `json:"partition" aolabel:"partition"`
	Description   string `json:"description"`
	Logo          string `json:"logo"`
	ApiNum        int64  `json:"api_num"`
	SubscriberNum int64  `json:"subscriber_num"`
}

type ServiceDetail struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Document    string             `json:"document"`
	Basic       *ServiceBasic      `json:"basic"`
	Apis        []*ServiceApi      `json:"apis"`
	DisableApis []*ServiceApiBasic `json:"disable_apis"`
	Partition   []*Partition       `json:"partition"`
}

type ServiceBasic struct {
	//Organization  auto.Label `json:"organization" aolabel:"organization"`
	Project       auto.Label `json:"project" aolabel:"project"`
	Team          auto.Label `json:"team" aolabel:"team"`
	ApiNum        int        `json:"api_num"`
	SubscriberNum int        `json:"subscriber_num"`
}

type ServiceApiBasic struct {
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Method      string         `json:"method"`
	Path        string         `json:"path"`
	Creator     auto.Label     `json:"creator" aolabel:"user"`
	Updater     auto.Label     `json:"updater" aolabel:"user"`
	CreateTime  auto.TimeLabel `json:"create_time"`
	UpdateTime  auto.TimeLabel `json:"update_time"`
}

type ServiceApi struct {
	*ServiceApiBasic
	Doc interface{} `json:"doc"`
}

type Partition struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Prefix string `json:"prefix"`
}
