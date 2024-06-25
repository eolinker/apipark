package project_authorization

import (
	"time"

	"github.com/eolinker/apipark/stores/project"
)

type Authorization struct {
	UUID           string
	Project        string
	Name           string
	Type           string
	Position       string
	TokenName      string
	Config         string
	Creator        string
	Updater        string
	CreateTime     time.Time
	UpdateTime     time.Time
	ExpireTime     int64
	HideCredential bool
}

func FromEntity(e *project.Authorization) *Authorization {
	return &Authorization{
		UUID:           e.UUID,
		Project:        e.Project,
		Name:           e.Name,
		Type:           e.Type,
		Position:       e.Position,
		TokenName:      e.TokenName,
		Config:         e.Config,
		Creator:        e.Creator,
		Updater:        e.Updater,
		CreateTime:     e.CreateAt,
		UpdateTime:     e.UpdateAt,
		ExpireTime:     e.ExpireTime,
		HideCredential: e.HideCredential,
	}
}

type CreateAuthorization struct {
	UUID           string
	Project        string
	Name           string
	Type           string
	Position       string
	TokenName      string
	Config         string
	AuthID         string
	ExpireTime     int64
	HideCredential bool
}

type EditAuthorization struct {
	Name           *string
	Position       *string
	TokenName      *string
	Config         *string
	ExpireTime     *int64
	HideCredential *bool
	AuthID         *string
}
