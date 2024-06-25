package api

import (
	"github.com/eolinker/apipark/model/plugin_model"
	"time"

	"github.com/eolinker/apipark/stores/api"
)

type API struct {
	UUID string
	//Name     string
	Project  string
	Team     string
	Creator  string
	Method   string
	Path     string
	CreateAt time.Time
	IsDelete bool
}

type APIInfo struct {
	//Id          int64
	UUID        string
	Name        string
	Description string
	CreateAt    time.Time
	UpdateAt    time.Time
	Project     string
	Team        string
	Creator     string
	Updater     string
	Upstream    string
	Method      string
	Path        string
	Match       string
}

func FromEntity(e *api.Api) *API {
	return &API{
		//Id:          e.Id,
		UUID: e.UUID,
		//Name:        e.Name,
		//Description: e.Description,
		CreateAt: e.CreateAt,
		IsDelete: e.IsDelete != 0,
		//UpdateAt:    e.UpdateAt,
		Project: e.Project,
		Team:    e.Team,
		Creator: e.Creator,
		//Updater:     e.Updater,
		//Upstream:    e.Upstream,
		Method: e.Method,
		Path:   e.Path,
		//Match:       e.Match,
	}
}

type CreateAPI struct {
	UUID        string
	Name        string
	Description string
	Project     string
	Team        string
	Method      string
	Path        string
	Match       string
	//Upstream    string
}

type EditAPI struct {
	//UUID        string
	Name        *string
	Upstream    *string
	Description *string
}

type ExistAPI struct {
	Path   string
	Method string
}

type Document struct {
	Content string `json:"content"`
}
type PluginSetting struct {
	Disable bool                    `json:"disable"`
	Config  plugin_model.ConfigType `json:"config"`
}
type Proxy struct {
	Path    string                   `json:"path"`
	Timeout int                      `json:"timeout"`
	Retry   int                      `json:"retry"`
	Plugins map[string]PluginSetting `json:"plugins"`
	Extends map[string]any           `json:"extends"`
	Headers []*Header                `json:"headers"`
}

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Opt   string `json:"opt"`
}

type Router struct {
	Method     string   `json:"method"`
	Path       string   `json:"path"`
	MatchRules []*Match `json:"match"`
}

type Match struct {
	Position  string `json:"position"`
	MatchType string `json:"match_type"`
	Key       string `json:"key"`
	Pattern   string `json:"pattern"`
}
