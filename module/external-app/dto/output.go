package external_app_dto

import (
	"github.com/eolinker/go-common/auto"
)

type ExternalApp struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type ExternalAppItem struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Token      string         `json:"token"`
	Tags       string         `json:"tags"`
	Status     bool           `json:"status"`
	Operator   auto.Label     `json:"operator" aolabel:"user"`
	UpdateTime auto.TimeLabel `json:"update_time"`
}
